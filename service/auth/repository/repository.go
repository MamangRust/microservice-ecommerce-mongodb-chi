package repository

import (
	pbrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/role"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc-pb/user"
	pbuserrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/user_role"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	roleadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/role"
	useradapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user"
	userroleadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user_role"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repositories struct {
	User         UserRepository
	RefreshToken RefreshTokenRepository
	UserRole     userroleadapter.CommandRepository
	Role         roleadapter.QueryRepository
	ResetToken   ResetTokenRepository
}

type GuardOptions struct {
	User     []adapter.GuardOption
	UserRole []adapter.GuardOption
	Role     []adapter.GuardOption
}

var (
	_ UserRepository                    = (*useradapter.Repository)(nil)
	_ roleadapter.QueryRepository       = (*roleadapter.Repository)(nil)
	_ userroleadapter.CommandRepository = (*userroleadapter.Repository)(nil)
)

type Deps struct {
	Db          *mongo.Database
	User        pbuser.UserQueryServiceClient
	UserCommand pbuser.UserCommandServiceClient
	Role        pbrole.RoleQueryServiceClient
	UserRole    pbuserrole.UserRoleServiceClient
	Guards      GuardOptions
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		User:         useradapter.New(deps.User, deps.UserCommand, deps.Guards.User...),
		RefreshToken: NewRefreshTokenRepository(deps.Db),
		UserRole:     userroleadapter.New(deps.UserRole, deps.Guards.UserRole...),
		Role:         roleadapter.New(deps.Role, deps.Guards.Role...),
		ResetToken:   NewResetTokenRepository(deps.Db),
	}
}
