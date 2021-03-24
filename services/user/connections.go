package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/oauth/spotify"
	"github.com/CastyLab/grpc.server/services/auth"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) UpdateConnection(ctx context.Context, req *proto.ConnectionRequest) (*proto.ConnectionsResponse, error) {

	var (
		connection = new(models.Connection)
		collection = db.Connection.Collection("connections")
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"type":    req.Connection.Type,
		"user_id": user.ID,
	}

	if err := collection.FindOne(ctx, filter).Decode(connection); err != nil {
		return nil, status.Error(codes.NotFound, "Could not find connection!")
	}

	token, err := spotify.RefreshToken(connection.RefreshedToken)
	if err != nil {
		sentry.CaptureException(err)
		return nil, status.Errorf(codes.Unauthenticated, "could not refresh the token")
	}

	var (
		updateFilter  = bson.M{"_id": connection.ID}
		updatePayload = bson.M{
			"$set": bson.M{
				"access_token": token.AccessToken,
				"updated_at":   time.Now(),
			},
		}
	)

	result, err := collection.UpdateOne(ctx, updateFilter, updatePayload)
	if err != nil {
		sentry.CaptureException(err)
		return nil, status.Errorf(codes.Unauthenticated, "could not refresh the token")
	}

	if result.ModifiedCount == 1 {

		connection.AccessToken = token.AccessToken
		connection.UpdatedAt = time.Now()

		return &proto.ConnectionsResponse{
			Status: "success",
			Code:   http.StatusOK,
			Result: []*proto.Connection{helpers.NewProtoConnection(connection)},
		}, nil
	}

	return nil, status.Errorf(codes.Aborted, "could not update the token in db")
}

func (s *Service) GetConnection(ctx context.Context, req *proto.ConnectionRequest) (*proto.ConnectionsResponse, error) {

	var (
		connection = new(models.Connection)
		collection = db.Connection.Collection("connections")
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"type":    req.Connection.Type,
		"user_id": user.ID,
	}

	if err := collection.FindOne(ctx, filter).Decode(connection); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "Could not find connection!")
		}
		return nil, fmt.Errorf("could not get connection :%v", err)
	}

	return &proto.ConnectionsResponse{
		Status: "success",
		Code:   http.StatusOK,
		Result: []*proto.Connection{helpers.NewProtoConnection(connection)},
	}, nil
}

func (s *Service) GetConnections(ctx context.Context, req *proto.AuthenticateRequest) (*proto.ConnectionsResponse, error) {

	var (
		connections = make([]*proto.Connection, 0)
		collection  = db.Connection.Collection("connections")
	)

	user, err := auth.Authenticate(req)
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, bson.M{"user_id": user.ID})
	if err != nil {
		return nil, status.Error(codes.NotFound, "Could not find connections!")
	}

	for cursor.Next(ctx) {
		connection := new(models.Connection)
		if err := cursor.Decode(connection); err != nil {
			continue
		}
		connections = append(connections, helpers.NewProtoConnection(connection))
	}

	return &proto.ConnectionsResponse{
		Status: "success",
		Code:   http.StatusOK,
		Result: connections,
	}, nil

}
