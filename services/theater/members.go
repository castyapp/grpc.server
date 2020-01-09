package theater

import (
	"context"
	"gitlab.com/movienight1/grpc.proto"
	"gitlab.com/movienight1/grpc.proto/messages"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"movie.night.gRPC.server/db"
	"movie.night.gRPC.server/db/models"
	"movie.night.gRPC.server/services/auth"
	"movie.night.gRPC.server/services/user"
	"net/http"
	"time"
)

func (s *Service) GetMembers(ctx context.Context, req *proto.GetTheaterMembersRequest) (*proto.TheaterMembersResponse, error) {

	var (
		collection  = db.Connection.Collection("theater_members")
		usersCollection = db.Connection.Collection("users")
		emptyResponse   = &proto.TheaterMembersResponse{
			Status:  "success",
			Code:    http.StatusOK,
			Result:  make([]*messages.User, 0),
		}
	)

	mCtx, _ := context.WithTimeout(ctx, 10 * time.Second)

	if _, err := auth.Authenticate(req.AuthRequest); err != nil {
		return &proto.TheaterMembersResponse{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	theaterObjectId, err := primitive.ObjectIDFromHex(req.TheaterId)
	if err != nil {
		return emptyResponse, err
	}

	cursor, err := collection.Find(mCtx, bson.M{ "theater_id": theaterObjectId })
	if err != nil {
		return emptyResponse, err
	}

	var users = make([]*messages.User, 0)
	for cursor.Next(mCtx) {
		dbTheaterMember := new(models.TheaterMember)
		if err := cursor.Decode(dbTheaterMember); err != nil {
			break
		}
		var theaterMemberUser = new(models.User)
		findUser := usersCollection.FindOne(mCtx, bson.M{ "_id": dbTheaterMember.UserId })
		if err := findUser.Decode(theaterMemberUser); err != nil {
			break
		}
		protoUser, err := user.SetDBUserToProtoUser(theaterMemberUser)
		if err != nil {
			break
		}
		users = append(users, protoUser)
	}

	return &proto.TheaterMembersResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Result:  users,
	}, nil
}

func (s *Service) AddMember(ctx context.Context, req *proto.AddOrRemoveMemberRequest) (*proto.Response, error) {

	var (
		theater         = new(models.Theater)
		theaterCol      = db.Connection.Collection("theaters")
		collection      = db.Connection.Collection("theater_members")
		failedResponse  = &proto.Response{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Could not add member to theater!",
		}
	)

	mCtx, _ := context.WithTimeout(ctx, 10 * time.Second)

	member, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	theaterObjectId, err := primitive.ObjectIDFromHex(req.TheaterId)
	if err != nil {
		return failedResponse, err
	}; {
		if err := theaterCol.FindOne(mCtx, bson.M{ "_id": theaterObjectId }).Decode(theater); err != nil {
			return failedResponse, err
		}
	}

	count, err := collection.CountDocuments(mCtx, bson.M{"theater_id": theater.ID, "user_id": member.ID})
	if err != nil {
		return failedResponse, err
	}

	if count == 0 {

		// adding member to theater_members collection
		theaterMember := bson.M{
			"theater_id": theater.ID,
			"user_id":    member.ID,
			"created_at": time.Now(),
		}

		if _, err := collection.InsertOne(mCtx, theaterMember); err != nil {
			return failedResponse, nil
		}

		return &proto.Response{
			Status:  "success",
			Code:    http.StatusOK,
			Message: "User added to theater successfully!",
		}, nil
	}

	return failedResponse, nil
}

func (s *Service) RemoveMember(ctx context.Context, req *proto.AddOrRemoveMemberRequest) (*proto.Response, error) {

	var (
		collection      = db.Connection.Collection("theater_members")
		failedResponse  = &proto.Response{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Could not add member to theater!",
		}
	)

	mCtx, _ := context.WithTimeout(ctx, 10 * time.Second)

	member, err := auth.Authenticate(req.AuthRequest)
	if err != nil {
		return &proto.Response{
			Status:  "failed",
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized!",
		}, nil
	}

	theaterObjectId, err := primitive.ObjectIDFromHex(req.TheaterId)
	if err != nil {
		return failedResponse, err
	}

	// removing member from theater_members collection
	filter := bson.M{
		"theater_id": theaterObjectId,
		"user_id": member.ID,
	}

	if _, err := collection.DeleteOne(mCtx, filter); err != nil {
		return failedResponse, nil
	}

	return &proto.Response{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "User removed from theater successfully!",
	}, nil
}
