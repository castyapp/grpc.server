package theater

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"movie.night.gRPC.server/db"
	"movie.night.gRPC.server/db/models"
	"movie.night.gRPC.server/proto"
	"movie.night.gRPC.server/proto/messages"
	"movie.night.gRPC.server/services/auth"
	"movie.night.gRPC.server/services/user"
	"net/http"
	"time"
)

func (s *Service) GetMembers(ctx context.Context, req *proto.GetTheaterMembersRequest) (*proto.TheaterMembersResponse, error) {

	var (
		theater     = new(models.Theater)
		collection  = db.Connection.Collection("theater_members")
		theaterCol  = db.Connection.Collection("theaters")
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

	if err := theaterCol.FindOne(mCtx, bson.M{ "_id": theaterObjectId }).Decode(theater); err != nil {
		return emptyResponse, err
	}

	cursor, err := collection.Find(mCtx, bson.M{ "theater_id": theater.ID })
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
		member          = new(models.User)
		theaterCol      = db.Connection.Collection("theaters")
		collection      = db.Connection.Collection("theater_members")
		usersCollection = db.Connection.Collection("users")
		failedResponse  = &proto.Response{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Could not add member to theater!",
		}
	)

	mCtx, _ := context.WithTimeout(ctx, 10 * time.Second)

	if _, err := auth.Authenticate(req.AuthRequest); err != nil {
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

	memberObjectId, err := primitive.ObjectIDFromHex(req.MemberId)
	if err != nil {
		return failedResponse, err
	}; {
		if err := usersCollection.FindOne(mCtx, bson.M{ "_id": memberObjectId }).Decode(member); err != nil {
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

func (s *Service) RemoveMember(ctx context.Context, request *proto.AddOrRemoveMemberRequest) (*proto.Response, error) {



	return nil, nil
}
