package services_test

import (
	"encoder/application/services"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		logrus.Fatalf("Error loading .env file")
	}
}

func TestVideoServiceUpload(t *testing.T) {
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

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "full-cycle-encoder"
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + video.ID
	
	dodeUpload := make(chan string)
	go videoUpload.ProcessUpload(50, dodeUpload)

	result := <-dodeUpload
	require.Equal(t, result, services.MsgUploadCompleted)
}
