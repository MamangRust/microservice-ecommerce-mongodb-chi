package repository

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repositories struct {
	ReviewDetailQuery   ReviewDetailQueryRepository
	ReviewDetailCommand ReviewDetailCommandRepository
}

func NewRepositories(DB *mongo.Database) *Repositories {
	return &Repositories{
		ReviewDetailQuery:   NewReviewDetailQueryRepository(DB),
		ReviewDetailCommand: NewReviewDetailCommandRepository(DB),
	}
}
