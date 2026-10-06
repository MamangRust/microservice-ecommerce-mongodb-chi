package repository

import "go.mongodb.org/mongo-driver/v2/mongo"

type Repositories struct {
	BannerQuery   BannerQueryRepository
	BannerCommand BannerCommandRepository
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		BannerQuery:   NewBannerQueryRepository(db),
		BannerCommand: NewBannerCommandRepository(db),
	}
}
