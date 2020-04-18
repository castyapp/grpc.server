package handlers

import (
	"github.com/CastyLab/grpc.server/jwt"
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Authentication(ctx *gin.Context)  {

	token := ctx.GetHeader("Authorization")
	if token == "" {

		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(422).
			SetStatusText("Failed!").
			RespondWithMessage("Token is required!"))
		return
	}

	stringToken := strings.ReplaceAll(token, "Bearer ", "")
	user, err := jwt.DecodeAuthToken([]byte(stringToken))
	if err != nil {

		ctx.AbortWithStatusJSON(respond.Default.Error(http.StatusUnauthorized, 3012))
		return
	}

	ctx.Set("user", user)
	ctx.Next()
}