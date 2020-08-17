package theater

import "go.mongodb.org/mongo-driver/mongo"

type Service struct {
	db *mongo.Database
}