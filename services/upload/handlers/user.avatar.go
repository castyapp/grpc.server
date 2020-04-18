package handlers

import (
	"fmt"
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/services"
	"github.com/MrJoshLab/go-respond"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/thedevsaddam/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"os"
)

func UserAvatarUpload(ctx *gin.Context)  {

	var (
		user = ctx.MustGet("user").(*models.User)
		rules = govalidator.MapData{
			"file:avatar": []string{"ext:jpg,jpeg,png", "size:2000000"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: true,
		}
		collection = db.Connection.Collection("users")
	)

	if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {

		validations := helpers.GetValidationErrorsFromGoValidator(validate)
		ctx.JSON(respond.Default.ValidationErrors(validations))
		return
	}

	avatarFile, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusBadRequest).
			SetStatusText("Failed!").
			RespondWithMessage("Bad request!"))
		return
	}

	var (
		storagePath = os.Getenv("STORAGE_PATH")
		avatar = services.RandomNumber(20)
		avatarPath = fmt.Sprintf("%s/uploads/avatars/%s.png", storagePath, avatar)
	)

	if err := ctx.SaveUploadedFile(avatarFile, avatarPath); err != nil {
		sentry.CaptureException(err)
		ctx.JSON(respond.Default.SetStatusText("Failed!").
			SetStatusCode(http.StatusInternalServerError).
			RespondWithMessage("Internal server error, Please try again later!"))
		return
	}

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{"avatar": avatar}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		ctx.JSON(respond.Default.SetStatusText("Failed!").
			SetStatusCode(http.StatusBadRequest).
			RespondWithMessage("Could not update user's avatar!"))
		return
	}

	if result.ModifiedCount == 1 {
		ctx.JSON(respond.Default.Succeed(map[string] interface{} {
			"avatar": avatar,
		}))
		return
	}

	ctx.JSON(respond.Default.SetStatusText("Failed!").
		SetStatusCode(http.StatusBadRequest).
		RespondWithMessage("Bad Request!"))
	return
}