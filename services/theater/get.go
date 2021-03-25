package theater

import (
	"context"
	"github.com/CastyLab/grpc.proto/proto"
	"github.com/castyapp/grpc.server/db"
	"github.com/castyapp/grpc.server/db/models"
	"github.com/castyapp/grpc.server/helpers"
	"github.com/castyapp/grpc.server/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

func (s *Service) GetTheater(ctx context.Context, req *proto.GetTheaterRequest) (*proto.UserTheaterResponse, error) {

	var (
		err error
		database        = db.Connection
		authenticated   = false
		authUser        = new(models.User)
		dbTheater       = new(models.Theater)
		usersCollection = database.Collection("users")
		collection      = database.Collection("theaters")
		failedResponse  = status.Error(codes.Internal, "Could not get theater, Please try again later!")
	)

	if req.AuthRequest != nil {
		authUser, err = auth.Authenticate(req.AuthRequest)
		if err != nil {
			return nil, err
		}
		authenticated = true
	}

	if req.TheaterId != "" {
		theaterObjectId, err := primitive.ObjectIDFromHex(req.TheaterId)
		if err != nil {
			return nil, status.Error(codes.NotFound, "Theater object id is invalid!")
		}
		if err := collection.FindOne(ctx, bson.M{"_id": theaterObjectId}).Decode(dbTheater); err != nil {
			return nil, status.Error(codes.NotFound, "Could not find theater!")
		}
	} else if req.User != "" {
		user := new(models.User)
		if err := usersCollection.FindOne(ctx, bson.M{ "username": req.User }).Decode(user); err != nil {
			return nil, status.Error(codes.NotFound, "Could not find the user!")
		}
		if err := collection.FindOne(ctx, bson.M{"user_id": user.ID}).Decode(dbTheater); err != nil {
			return nil, status.Error(codes.NotFound, "Could not find theater!")
		}
	} else {
		if err := collection.FindOne(ctx, bson.M{"user_id": authUser.ID}).Decode(dbTheater); err != nil {
			return nil, status.Error(codes.NotFound, "Could not find theater!")
		}
	}

	if !authenticated {
		switch dbTheater.Privacy {
		case proto.PRIVACY_PRIVATE:
			return nil, status.Error(codes.PermissionDenied, "Permission Denied!")
		}
	} else {
		if dbTheater.UserId.Hex() != authUser.ID.Hex() {
			switch dbTheater.Privacy {
			case proto.PRIVACY_PRIVATE:
				return nil, status.Error(codes.PermissionDenied, "Permission Denied!")
			}
		}
	}

	theater, err := helpers.NewTheaterProto(ctx, dbTheater)
	if err != nil {
		log.Println(err)
		return nil, failedResponse
	}

	theater.Followed = false

	if req.AuthRequest != nil {
		findFilter := bson.M{ "theater_id": dbTheater.ID, "user_id": authUser.ID }
		countResult, err := database.Collection("follows").CountDocuments(ctx, findFilter)
		if err == nil && countResult != 0 {
			theater.Followed = true
		}
	}

	return &proto.UserTheaterResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Result:  theater,
	}, nil
}
