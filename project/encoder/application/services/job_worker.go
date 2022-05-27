package services

import (
	"encoder/domain"
	"encoder/framework/utils"
	"encoding/json"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

func JobWorker(messageChannel chan amqp.Delivery, returnChan chan JobWorkerResult, jobService JobService, job domain.Job, workerId int) {
	for message := range messageChannel {
		if err := utils.IsJson(string(message.Body)); err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		if err := json.Unmarshal(message.Body, &jobService.VideoService.Video); err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}
		jobService.VideoService.Video.ID = uuid.NewV4().String()

		if err := jobService.VideoService.Video.Validate(); err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}
		
		if err := jobService.VideoService.InsertVideo(); err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		job.Video = jobService.Job.Video
		job.ID = uuid.NewV4().String()
		job.Status = domain.StatusStart
		job.OutputBucketPath = os.Getenv("outputBucketName")
		job.CreatedAt = time.Now()
		if _, err := jobService.JobRepository.Insert(&job); err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		jobService.Job = &job

		if err := jobService.Start(); err != nil {
			returnChan <- returnJobResult(job, message, err)
			continue
		}

		returnChan <-returnJobResult(job, message, nil)
	}
}

func returnJobResult(job domain.Job, message amqp.Delivery, err error) JobWorkerResult {
	return JobWorkerResult{
		Job: job,
		Message: &message,
		Error: err,
	}
}