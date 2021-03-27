package helpers

import (
	"context"
	"fmt"

	"github.com/CastyLab/grpc.proto/proto"
	"github.com/castyapp/grpc.server/core"
	"github.com/castyapp/grpc.server/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewMemberProto(ctx *core.Context, member *models.TheaterMember) (*proto.User, error) {

	dbConn, err := ctx.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db       = dbConn.(*mongo.Database)
		dbmember = new(models.User)
		decoder  = db.Collection("users").
				FindOne(context.Background(), bson.M{"_id": member.UserId})
	)
	if err := decoder.Decode(dbmember); err != nil {
		return nil, fmt.Errorf("could not decode theater member: %v", err)
	}
	return NewProtoUser(dbmember), nil
}

func GetTheaterMembers(ctx *core.Context, theater *models.Theater) ([]*proto.User, error) {

	dbConn, err := ctx.Get("db.mongo")
	if err != nil {
		return nil, err
	}

	var (
		db      = dbConn.(*mongo.Database)
		members = make([]*proto.User, 0)
	)

	cursor, err := db.Collection("theater_members").Find(ctx, bson.M{"theater_id": theater.ID})
	if err != nil {
		return nil, status.Error(codes.Internal, "Could not get theater members")
	}

	for cursor.Next(ctx) {
		member := new(models.TheaterMember)
		if err := cursor.Decode(member); err != nil {
			continue
		}
		protoMember, err := NewMemberProto(ctx, member)
		if err != nil {
			continue
		}
		members = append(members, protoMember)
	}

	return members, nil
}
