package repository

import (
	pbrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/role"
	pbuserrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/user_role"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	roleadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/role"
	userroleadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user_role"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GuardOptions struct {
	Role     []adapter.GuardOption
	UserRole []adapter.GuardOption
}

type Repositories struct {
	UserCommand UserCommandRepository
	UserQuery   UserQueryRepository
	Role        roleadapter.QueryRepository
	UserRole    userroleadapter.CommandRepository
}

type Deps struct {
	Db              *mongo.Database
	RoleQueryClient pbrole.RoleQueryServiceClient
	UserRoleClient  pbuserrole.UserRoleServiceClient
	Guard           GuardOptions
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		UserCommand: NewUserCommandRepository(deps.Db),
		UserQuery:   NewUserQueryRepository(deps.Db),
		Role:        roleadapter.New(deps.RoleQueryClient, deps.Guard.Role...),
		UserRole:    userroleadapter.New(deps.UserRoleClient, deps.Guard.UserRole...),
	}
}
