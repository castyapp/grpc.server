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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"os"
)

func TheaterPosterUpload(ctx *gin.Context)  {

	var (
		user = ctx.MustGet("user").(*models.User)
		collection = db.Connection.Collection("theaters")
		theater = new(models.Theater)
		rules = govalidator.MapData{
			"file:poster": []string{"ext:jpg,jpeg,png", "size:2000000"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: true,
		}
	)

	if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {

		validations := helpers.GetValidationErrorsFromGoValidator(validate)
		ctx.JSON(respond.Default.ValidationErrors(validations))
		return
	}

	theaterObjectId, err := primitive.ObjectIDFromHex(ctx.Param("theater_id"))
	if err != nil {
		log.Println(err)
		return
	}

	filter := bson.M{
		"_id": theaterObjectId,
		"user_id": user.ID,
	}

	if err := collection.FindOne(ctx, filter).Decode(theater); err != nil {
		ctx.JSON(respond.Default.NotFound())
		return
	}

	storagePath := os.Getenv("STORAGE_PATH")
	posterFile, err := ctx.FormFile("poster")
	if err != nil {

		ctx.JSON(respond.Default.SetStatusCode(http.StatusBadRequest).
			SetStatusText("Failed!").
			RespondWithMessage("Bad request!"))
		return
	}

	fileName := services.RandomNumber(20)
	filePath := fmt.Sprintf("%s/uploads/posters/%s.png", storagePath, fileName)
	if err := ctx.SaveUploadedFile(posterFile, filePath); err != nil {
		sentry.CaptureException(err)
		ctx.JSON(respond.Default.SetStatusText("Failed!").
			SetStatusCode(http.StatusInternalServerError).
			RespondWithMessage("Internal server error, Please try again later!"))
		return
	}

	updateFilter := bson.M{"_id": theater.ID}
	update := bson.M{
		"$set": bson.M{
			"movie.poster": fileName,
		},
	}

	result, err := collection.UpdateOne(ctx, updateFilter, update)
	if err != nil {
		ctx.JSON(respond.Default.SetStatusText("Failed!").
			SetStatusCode(http.StatusBadRequest).
			RespondWithMessage("Could not update theater's poster!"))
		return
	}

	if result.ModifiedCount == 1 {
		ctx.JSON(respond.Default.UpdateSucceeded())
		return
	}

	ctx.JSON(respond.Default.SetStatusText("Failed!").
		SetStatusCode(http.StatusBadRequest).
		RespondWithMessage("Bad Request!"))
	return
}
