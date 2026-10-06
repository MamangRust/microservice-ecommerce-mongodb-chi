package merchant_policy_test

import (
	"context"
	policy_handler "github.com/MamangRust/microservice-ecommerce-grpc-merchant_policy/handler"
	policy_repo "github.com/MamangRust/microservice-ecommerce-grpc-merchant_policy/repository"
	policy_service "github.com/MamangRust/microservice-ecommerce-grpc-merchant_policy/service"
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	pbmerchant_policy "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant_policy"
	"testing"

	policy_cache "github.com/MamangRust/microservice-ecommerce-grpc-merchant_policy/cache"
	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	tests "github.com/MamangRust/microservice-ecommerce-test"
	"google.golang.org/protobuf/types/known/emptypb"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type MerchantPolicyGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient   pbmerchant_policy.MerchantPolicyQueryServiceClient
	commandClient pbmerchant_policy.MerchantPolicyCommandServiceClient
}

func (s *MerchantPolicyGapiTestSuite) SetupSuite() {
	s.PrimaryContext = database.MerchantDB
	s.BaseTestSuite.SetupSuite()

	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	// Setup dependencies
	s.SetupUserService()
	s.SetupMerchantService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)

	// Policy dependencies
	mencache := policy_cache.NewMencache(cacheStore)
	repos := policy_repo.NewRepositories(
		mongoDB,
		pbmerchant.NewMerchantQueryServiceClient(s.Conns["merchant"]),
	)
	svc := policy_service.NewService(&policy_service.Deps{
		Cache:         mencache,
		Repository:    repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := policy_handler.NewHandler(&policy_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pbmerchant_policy.RegisterMerchantPolicyQueryServiceServer(server, handler.MerchantPolicyQuery)
	pbmerchant_policy.RegisterMerchantPolicyCommandServiceServer(server, handler.MerchantPolicyCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pbmerchant_policy.NewMerchantPolicyQueryServiceClient(conn)
	s.commandClient = pbmerchant_policy.NewMerchantPolicyCommandServiceClient(conn)
}

func (s *MerchantPolicyGapiTestSuite) TestMerchantPolicyGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := int32(s.SeedMerchant(ctx, userID))

	// 2. Create
	createRes, err := s.commandClient.Create(ctx, &pbmerchant_policy.CreateMerchantPoliciesRequest{
		MerchantId:  merchID,
		PolicyType:  "Warranty",
		Title:       "Warranty Policy",
		Description: "1 year limited warranty",
	})
	s.NoError(err)
	s.NotNil(createRes)
	policyID := createRes.Data.Id

	// 3. FindById
	getRes, err := s.queryClient.FindById(ctx, &pbmerchant_policy.FindByIdMerchantPoliciesRequest{Id: policyID})
	s.NoError(err)
	s.Equal("Warranty", getRes.Data.PolicyType)

	// 4. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pbmerchant.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pbmerchant.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	updateRes, err := s.commandClient.Update(ctx, &pbmerchant_policy.UpdateMerchantPoliciesRequest{
		MerchantPolicyId: policyID,
		PolicyType:       "Warranty Updated",
		Title:            "Warranty Policy Updated",
		Description:      "2 years limited warranty",
	})
	s.NoError(err)
	s.Equal("Warranty Updated", updateRes.Data.PolicyType)

	// 7. Trash
	_, err = s.commandClient.TrashedMerchantPolicies(ctx, &pbmerchant_policy.FindByIdMerchantPoliciesRequest{Id: policyID})
	s.NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pbmerchant.FindAllMerchantRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.commandClient.RestoreMerchantPolicies(ctx, &pbmerchant_policy.FindByIdMerchantPoliciesRequest{Id: policyID})
	s.NoError(err)

	// 10. DeletePermanent
	_, _ = s.commandClient.TrashedMerchantPolicies(ctx, &pbmerchant_policy.FindByIdMerchantPoliciesRequest{Id: policyID})
	_, err = s.commandClient.DeleteMerchantPoliciesPermanent(ctx, &pbmerchant_policy.FindByIdMerchantPoliciesRequest{Id: policyID})
	s.NoError(err)

	// 11. RestoreAll
	_, err = s.commandClient.RestoreAllMerchantPolicies(ctx, &emptypb.Empty{})
	s.NoError(err)

	// 12. DeleteAll
	_, err = s.commandClient.DeleteAllMerchantPoliciesPermanent(ctx, &emptypb.Empty{})
	s.NoError(err)
}

func TestMerchantPolicyGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantPolicyGapiTestSuite))
}
