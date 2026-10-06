package repository

import "go.mongodb.org/mongo-driver/v2/mongo"

type Repositories struct {
	SliderQuery   SliderQueryRepository
	SliderCommand SliderCommandRepository
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		SliderQuery:   NewSliderQueryRepository(db),
		SliderCommand: NewSliderCommandRepository(db),
	}
}
