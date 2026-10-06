package cache

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-user/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type UserQueryCache interface {
	GetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, bool)
	SetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers, data []*repository.UserResult, total *int)

	GetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, bool)
	SetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers, data []*repository.UserResult, total *int)

	GetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, bool)
	SetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers, data []*repository.UserResult, total *int)

	GetCachedUserCache(ctx context.Context, id int) (*models.User, bool)
	SetCachedUserCache(ctx context.Context, data *models.User)
}

type UserCommandCache interface {
	DeleteUserCache(ctx context.Context, id int)
}
