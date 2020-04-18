package handlers

import (
	"github.com/CastyLab/grpc.server/db"
	"github.com/CastyLab/grpc.server/db/models"
	"github.com/CastyLab/grpc.server/helpers"
	"github.com/CastyLab/grpc.server/helpers/subtitle"
	"github.com/MrJoshLab/go-respond"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func TheaterSubtitlesUpload(ctx *gin.Context)  {

	var (
		user = ctx.MustGet("user").(*models.User)
		theater = new(models.Theater)
		collection = db.Connection.Collection("theaters")
		ccCollection = db.Connection.Collection("subtitles")
		rules = govalidator.MapData{
			"lang":           []string{"required", "min:4", "max:30"},
			"file:subtitle":  []string{"required", "ext:srt", "size:20000000"},
		}
		opts = govalidator.Options{
			Request:         ctx.Request,
			Rules:           rules,
			RequiredDefault: true,
		}
	)

	theaterObjectId, err := primitive.ObjectIDFromHex(ctx.Param("theater_id"))
	if err != nil {
		ctx.JSON(respond.Default.NotFound())
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

	if validate := govalidator.New(opts).Validate(); validate.Encode() != "" {

		validations := helpers.GetValidationErrorsFromGoValidator(validate)
		ctx.JSON(respond.Default.ValidationErrors(validations))
		return
	}

	if subtitleFile, err := ctx.FormFile("subtitle"); err == nil {

		filename, err := subtitle.Save(subtitleFile)
		if err != nil {
			sentry.CaptureException(err)
			ctx.JSON(respond.Default.
				SetStatusText("Failed!").
				SetStatusCode(400).
				RespondWithMessage("Upload failed. Please try again later!"))
			return
		}

		document := bson.M{
			"theater_id": theater.ID,
			"lang": ctx.PostForm("lang"),
			"file": filename,
			"created_at": time.Now(),
			"updated_at": time.Now(),
		}

		if _, err := ccCollection.InsertOne(ctx, document); err != nil {
			ctx.JSON(respond.Default.SetStatusText("failed").
				SetStatusCode(http.StatusBadRequest).
				RespondWithMessage("Could not add subtitle, please try again later!"))
			return
		}

		ctx.JSON(respond.Default.InsertSucceeded())
		return
	}

	ctx.JSON(respond.Default.SetStatusText("failed").
		SetStatusCode(http.StatusBadRequest).
		RespondWithMessage("Could not add subtitle, please try again later!"))
	return
}
