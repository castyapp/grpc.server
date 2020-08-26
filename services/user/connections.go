package user

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func (s *Service) GetConnection(ctx context.Context, req *proto.GetConnectionRequest) (*proto.ConnectionsResponse, error) {

	var (
		connections = make([]*proto.Connection, 0)
		collection  = db.Connection.Collection("connections")
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"type":    req.Connection.Type,
		"user_id": user.ID,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Could not find connections!")
	}

	for cursor.Next(ctx) {
		connection := new(models.Connection)
		if err := cursor.Decode(connection); err != nil {
			continue
		}
		protoConnection, err := helpers.NewProtoConnection(connection)
		if err != nil {
			continue
		}
		connections = append(connections, protoConnection)
	}

	return &proto.ConnectionsResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Result:  connections,
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

	cursor, err := collection.Find(ctx, bson.M{ "user_id": user.ID })
	if err != nil {
		return nil, status.Error(codes.NotFound, "Could not find connections!")
	}

	for cursor.Next(ctx) {
		connection := new(models.Connection)
		if err := cursor.Decode(connection); err != nil {
			continue
		}
		protoConnection, err := helpers.NewProtoConnection(connection)
		if err != nil {
			continue
		}
		connections = append(connections, protoConnection)
	}

	return &proto.ConnectionsResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Result:  connections,
	}, nil

}
