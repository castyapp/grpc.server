package theater

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type Service struct {}

func (s *Service) GetUserTheaters(ctx context.Context, req *proto.GetAllUserTheatersRequest) (*proto.UserTheatersResponse, error) {

	var (
		theaters   = make([]*proto.Theater, 0)
		collection = db.Connection.Collection("theaters")
	)

	user, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return nil, err
	}

	qOpts := options.Find()
	qOpts.SetSort(bson.D{
		{"created_at", -1},
	})

	cursor, err := collection.Find(ctx, bson.M{"user_id": user.ID}, qOpts)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Could not find any theaters!")
	}

	for cursor.Next(ctx) {
		theater := new(models.Theater)
		if err := cursor.Decode(theater); err != nil {
			continue
		}
		th, err := helpers.NewTheaterProto(ctx, theater)
		if err != nil {
			continue
		}
		theaters = append(theaters, th)
	}

	return &proto.UserTheatersResponse{
		Result:  theaters,
		Code:    http.StatusOK,
		Message: "success",
	}, nil
}