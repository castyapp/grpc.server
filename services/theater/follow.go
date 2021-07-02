package theater

import (
	"context"
	"net/http"
	"time"

	"github.com/castyapp/grpc.server/helpers"
	"github.com/castyapp/grpc.server/models"
	"github.com/castyapp/grpc.server/services/auth"
	"github.com/castyapp/libcasty-protocol-go/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetFollowedTheaters(ctx context.Context, req *proto.AuthenticateRequest) (*proto.FollowedTheatersResponse, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db                = dbConn.(*mongo.Database)
		theaters          = make([]*proto.Theater, 0)
		followsCollection = db.Collection("follows")
	)

	user, err := auth.Authenticate(s.Context, req)
	if err != nil {
		return nil, err
	}

	cursor, err := followsCollection.Find(ctx, bson.M{"user_id": user.ID})
	if err != nil {
		return nil, status.Error(codes.NotFound, "Could not find any theaters!")
	}

	for cursor.Next(ctx) {
		follow := new(models.Follow)
		if err := cursor.Decode(follow); err != nil {
			continue
		}
		theater := new(models.Theater)
		err := db.Collection("theaters").FindOne(ctx, bson.M{"_id": follow.TheaterID}).Decode(theater)
		if err != nil {
			continue
		}
		protoTheater, err := helpers.NewTheaterProto(ctx, db, theater)
		if err != nil {
			continue
		}
		theaters = append(theaters, protoTheater)
	}

	return &proto.FollowedTheatersResponse{
		Status: "success",
		Code:   http.StatusOK,
		Result: theaters,
	}, nil
}

func (s *Service) Follow(ctx context.Context, req *proto.TheaterAuthRequest) (*proto.Response, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db                = dbConn.(*mongo.Database)
		theater           = new(models.Theater)
		followsCollection = db.Collection("follows")
		theaterCollection = db.Collection("theaters")
	)

	user, err := auth.Authenticate(s.Context, req.AuthRequest)
	if err != nil {
		return nil, err
	}

	theaterObjectID, err := primitive.ObjectIDFromHex(req.Theater.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Could not parse theater id!")
	}

	if err := theaterCollection.FindOne(ctx, bson.M{"_id": theaterObjectID}).Decode(theater); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find theater!")
	}

	countFollow := bson.M{
		"user_id":    user.ID,
		"theater_id": theater.ID,
	}

	countResult, err := followsCollection.CountDocuments(ctx, countFollow)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Could not find theater!")
	}

	if countResult == 0 {

		follow := bson.M{
			"user_id":            user.ID,
			"theater_id":         theater.ID,
			"email_notification": true,
			"push_notification":  true,
			"created_at":         time.Now(),
			"updated_at":         time.Now(),
		}

		if _, err := followsCollection.InsertOne(ctx, follow); err != nil {
			return nil, status.Error(codes.Internal, "Could not follow this theater!")
		}

		return &proto.Response{
			Status:  "success",
			Code:    http.StatusOK,
			Message: "Followed successfully!",
		}, nil
	}

	return nil, status.Error(codes.AlreadyExists, "Theater followed already!")
}

func (s *Service) Unfollow(ctx context.Context, req *proto.TheaterAuthRequest) (*proto.Response, error) {

	dbConn, err := s.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db                = dbConn.(*mongo.Database)
		follow            = new(models.Follow)
		followsCollection = db.Collection("follows")
	)

	user, err := auth.Authenticate(s.Context, req.AuthRequest)
	if err != nil {
		return nil, err
	}

	theaterObjectID, err := primitive.ObjectIDFromHex(req.Theater.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Could not parse theater id!")
	}

	findFilter := bson.M{
		"theater_id": theaterObjectID,
		"user_id":    user.ID,
	}

	if err := followsCollection.FindOne(ctx, findFilter).Decode(follow); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find theater!")
	}

	deletedResult, err := followsCollection.DeleteOne(ctx, bson.M{"_id": follow.ID})
	if err != nil {
		return nil, status.Error(codes.Internal, "Could not unfollow this theater!")
	}

	if deletedResult.DeletedCount == 1 {
		return &proto.Response{
			Status:  "success",
			Code:    http.StatusOK,
			Message: "Unfollowed successfully!",
		}, nil
	}

	return nil, status.Error(codes.Internal, "Could not unfollow this theater!")
}
