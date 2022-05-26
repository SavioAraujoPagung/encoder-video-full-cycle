package main

import (
	"encoder/application/repositories"
	"encoder/application/services"
	"encoder/domain"
	"encoder/framework/database"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Something noteworthy happened")
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "neymito.mp4"
	video.CreatedAt = time.Now()

	repo := &repositories.VideoRepositoryDb{
		Db: db,
	}

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo
	err := videoService.Download("full-cycle-encoder")
	if err != nil {

	}

}
