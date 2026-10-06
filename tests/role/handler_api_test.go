package role_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	pbrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/role"
	pbuserrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/user_role"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apicache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache"
	rolehandler "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/handler/role"
	role_cache "github.com/MamangRust/microservice-ecommerce-grpc-role/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-role/handler"
	"github.com/MamangRust/microservice-ecommerce-grpc-role/repository"
	"github.com/MamangRust/microservice-ecommerce-grpc-role/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	tests "github.com/MamangRust/microservice-ecommerce-test"

	"github.com/MamangRust/microservice-ecommerce-grpc-apigateway/apierror"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RoleApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	redisClient *redis.Client
	router      chi.Router
	grpcServer  *grpc.Server
	conn        *grpc.ClientConn
	roleID      int
}

func (s *RoleApiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite(database.IdentityDB)
	s.Require().NoError(err)
	s.ts = ts

	mongoDB, err := s.ts.MongoDB(s.ts.Ctx)
	s.Require().NoError(err)

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)

	repos := repository.NewRepositories(mongoDB)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	obs, _ := observability.NewObservability("test", log)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	mencache := role_cache.NewMencache(cacheStore)

	roleService := service.NewService(&service.Deps{
		Repository:    repos,
		Logger:        log,
		Cache:         mencache,
		Observability: obs,
	})

	// Start internal gRPC Server for Role module
	roleHandlerGrpc := handler.NewHandler(&handler.Deps{
		Service: roleService,
		Logger:  log,
	})
	server := grpc.NewServer()
	pbrole.RegisterRoleCommandServiceServer(server, roleHandlerGrpc.RoleCommand)
	pbrole.RegisterRoleQueryServiceServer(server, roleHandlerGrpc.RoleQuery)
	pbuserrole.RegisterUserRoleServiceServer(server, roleHandlerGrpc.UserRole)
	s.grpcServer = server

	lis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn

	// Setup Echo and API Handler
	e := chi.NewRouter()
	s.router = e

	// Bypass auth middleware by setting user_id and seeding roles in Redis
	e.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "user_id", 1))
			next.ServeHTTP(w, r)
		})
	})

	roles := []string{"Admin_Role_10", "Admin_Admin_14"}
	// F6-2: the gateway namespaces its cache keys with apigw: — the role
	// middleware reads cacheRoleKey "apigw:user_roles:%s", so the seeded cache
	// entry must use the same prefixed key.
	cache.SetToCache(s.ts.Ctx, cacheStore, "apigw:user_roles:1", &roles, 5*time.Minute)

	apiErrorHandler := apierror.NewApiHandler(obs, log)
	rolehandler.RegisterRoleHandler(&rolehandler.DepsRole{
		Client:     conn,
		Kafka:      nil,
		Router:     e,
		Logger:     log,
		CacheStore: cacheStore,
		Cache:      apicache.NewRoleCache(cacheStore),
		ApiHandler: apiErrorHandler,
	})
}

func (s *RoleApiTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.ts.Teardown()
}

func (s *RoleApiTestSuite) TestRoleApiLifecycle() {
	// 1. Create
	reqBody := requests.CreateRoleRequest{
		Name: "API Role",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/role-command/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	s.Equal(reqBody.Name, data["name"])
	s.roleID = int(data["id"].(float64))

	// 2. FindAll
	req = httptest.NewRequest(http.MethodGet, "/api/role-query", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. FindById
	s.Require().NotZero(s.roleID)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/role-query/%d", s.roleID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &res)
	data = res["data"].(map[string]interface{})
	s.Equal(float64(s.roleID), data["id"])

	// 4. FindByActive
	req = httptest.NewRequest(http.MethodGet, "/api/role-query/active", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. FindByTrashed
	req = httptest.NewRequest(http.MethodGet, "/api/role-query/trashed", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. FindByUserId
	req = httptest.NewRequest(http.MethodGet, "/api/role-query/user/1", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. Update
	updateBody := requests.UpdateRoleRequest{
		Name: "Updated API Role",
	}
	body, _ = json.Marshal(updateBody)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/role-command/update/%d", s.roleID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 8. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/role-command/trashed/%d", s.roleID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 9. FindByTrashed (verify trashed)
	req = httptest.NewRequest(http.MethodGet, "/api/role-query/trashed", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 10. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/role-command/restore/%d", s.roleID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 11. Trash again for DeletePermanent
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/role-command/trashed/%d", s.roleID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 12. DeletePermanent
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/role-command/permanent/%d", s.roleID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 13. RestoreAll
	req = httptest.NewRequest(http.MethodPost, "/api/role-command/restore/all", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 14. DeleteAll
	req = httptest.NewRequest(http.MethodPost, "/api/role-command/permanent/all", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestRoleApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(RoleApiTestSuite))
}
