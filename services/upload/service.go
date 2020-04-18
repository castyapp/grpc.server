package upload

import (
	"fmt"
	"github.com/CastyLab/grpc.server/services/upload/handlers"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func NewUploadService() {

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(handlers.CORSMiddleware)

	theaterService := router.Group("/json.TheaterService"); {
		theaterService.Use(handlers.Authentication)
		theaterService.POST(":theater_id/poster", handlers.TheaterPosterUpload)
		theaterService.POST(":theater_id/cc", handlers.TheaterSubtitlesUpload)
	}

	userService := router.Group("/json.UserService"); {
		userService.Use(handlers.Authentication)
		userService.POST("avatar", handlers.UserAvatarUpload)
	}

	unixFile := fmt.Sprintf("%s/upload.service.sock", os.Getenv("SOCKETS_PATH"))

	log.Println("UploadService is up and running on:", unixFile)
	log.Printf("UploadService Err: [%v]", router.RunUnix(unixFile))
}
