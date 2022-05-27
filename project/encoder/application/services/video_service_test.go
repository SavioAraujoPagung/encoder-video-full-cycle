package services_test

import (
	"encoder/application/repositories"
	"encoder/application/services"
	"encoder/domain"
	"encoder/framework/database"
	"testing"
	"time"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		logrus.Fatalf("Error loading .env file")
	}
}

func prepare() (*domain.Video, *repositories.VideoRepositoryDb) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "neymito.mp4"
	video.CreatedAt = time.Now()

	repo := &repositories.VideoRepositoryDb{
		Db: db,
	}

	return video, repo
}

func TestVideoDownloand(t *testing.T) {
	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo
	err := videoService.Download("full-cycle-encoder")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)
}
