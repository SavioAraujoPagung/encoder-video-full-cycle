package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/queue"
	"encoding/json"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type JobManager struct {
	Db               *gorm.DB
	Domain           domain.Job
	MessageChannel   chan amqp.Delivery
	JobReturnChannel chan JobWorkerResult
	RabbitMQ         *queue.RabbitMQ
}

type JobNotificationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewJobManager(db *gorm.DB,
	rabbitMQ *queue.RabbitMQ,
	jobReturnChannel chan JobWorkerResult,
	messageChannel chan amqp.Delivery) *JobManager {

	return &JobManager{
		Db:               db,
		Domain:           domain.Job{},
		MessageChannel:   messageChannel,
		JobReturnChannel: jobReturnChannel,
		RabbitMQ:         rabbitMQ,
	}
}

func (jm *JobManager) Start(ch *amqp.Channel) {
	videoService := NewVideoService()
	videoService.VideoRepository = &repositories.VideoRepositoryDb{Db: jm.Db}

	jobService := JobService{
		JobRepository: &repositories.JobRepositoryDb{Db: jm.Db},
		VideoService:  videoService,
	}

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKDER"))
	if err != nil {

	}

	for qtdProcesses := 0; qtdProcesses < concurrency; qtdProcesses++ {
		go JobWorker(jm.MessageChannel, jm.JobReturnChannel, jobService, jm.Domain, qtdProcesses)
	}

	for jobResult := range jm.JobReturnChannel {
		if jobResult.Error != nil {
			err = jm.checkParceErrors(jobResult)
		} else {
			err = jm.notifySuccess(jobResult, ch)
		}

		if err != nil {
			jobResult.Message.Reject(false)
		}
	}
}

func (jm *JobManager) checkParceErrors(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		logrus.Printf("messageID:", jobResult.Message.DeliveryTag, ". error parsing job:", jobResult.Job.ID)
	} else {
		logrus.Printf("messageID:", jobResult.Message.DeliveryTag, ". error parsing message:", jobResult.Error)
	}
	
	//falta implementar a notificação
	errMsg := JobNotificationError{
		Message: string(jobResult.Message.Body),
		Error: jobResult.Error.Error(),
	}

	jobJson, err := json.Marshal(errMsg)
	if err != nil {
		return err
	}

	err = jm.notify(jobJson)
	if err != nil {
		return err
	}

	err = jobResult.Message.Reject(false)
	if err != nil {
		return err
	}

	return nil
}

func (jm *JobManager) notifySuccess(jobResult JobWorkerResult, ch *amqp.Channel) error {
	jobJson, err := json.Marshal(jobResult)
	if err != nil {
		return err
	}

	err = jm.notify(jobJson)
	if err != nil {
		return err
	}

	err = jobResult.Message.Ack(false)
	if err != nil {
		return err
	}

	return nil
}

func (jm *JobManager) notify(jobJson []byte) error{
	err := jm.RabbitMQ.Notify(
		string(jobJson),
		"application/json",
		os.Getenv("RABBITMQ_NOTFICATION_EX"),
		os.Getenv("RABBITMQ_NOTFICATION_ROUTING_KEY"),
	)

	if err != nil {
		return err
	}

	return nil
}
