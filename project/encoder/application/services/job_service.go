package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"errors"
	"os"
	"strconv"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

func (js *JobService) Start() error {
	//downloand 
	if err := js.changeJobStatus(domain.StatusDownloading); err != nil {
		return js.failJob(err)
	}
	if err := js.VideoService.Download(os.Getenv("inputBucketName")); err != nil {
		return js.failJob(err)
	}
	//fragmentação
	if err := js.changeJobStatus(domain.StatusFragmenting); err != nil {
		return js.failJob(err)
	}
	if err := js.VideoService.Fragment(); err != nil {
		return js.failJob(err)
	}
	//encoder
	if err := js.changeJobStatus(domain.StatusEncoding); err != nil {
		return js.failJob(err)
	}
	if err := js.VideoService.Encode(); err != nil {
		return js.failJob(err)
	}
	//Upload
	if err := js.performUpload(); err != nil {
		return js.failJob(err)
	}
	//finaliza
	if err := js.changeJobStatus(domain.StatusFinishing); err != nil {
		return js.failJob(err)
	}
	if err := js.VideoService.Finish(); err != nil {
		return js.failJob(err)
	}
	//completo
	if err := js.changeJobStatus(domain.StatusCompleted); err != nil {
		return js.failJob(err)
	}

	return nil
}

func (js *JobService) performUpload() error {
	if err := js.changeJobStatus(domain.StatusUploading); err != nil {
		return js.failJob(err)
	}

	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv("outputBucketName")
	videoUpload.VideoPath = os.Getenv("localStoragePath" + "/" + js.VideoService.Video.ID)

	concurrency, _ := strconv.Atoi(os.Getenv("concurrencyUpload"))
	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(concurrency, doneUpload)

	var uploadResult string
	uploadResult = <-doneUpload

	if uploadResult != MsgUploadCompleted {
		return js.failJob(errors.New(uploadResult))
	}

	return nil
}

func (js *JobService) changeJobStatus(status string) error {
	var err error

	js.Job.Status = status
	js.Job, err = js.JobRepository.Update(js.Job)

	if err != nil {
		return js.failJob(err)
	}

	return nil
}

func (js *JobService) failJob(errJob error) error {
	js.Job.Status = domain.StatusFaild
	js.Job.Error = errJob.Error()

	_, err := js.JobRepository.Update(js.Job)
	if err != nil {
		return err
	}

	return errJob
}
