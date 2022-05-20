package repositories

import (
	"encoder/domain"
	"fmt"

	"github.com/jinzhu/gorm"
)

type JobRepository interface {
	Insert(job *domain.Job) (*domain.Job, error)
	Find(id string) (*domain.Job, error)
	Update(job *domain.Job) (*domain.Job, error)
}

type JobRepositoryDb struct {
	Db *gorm.DB
}

func (repo *JobRepositoryDb) Insert(job *domain.Job) (*domain.Job, error) {
	if err := repo.Db.Create(job).Error; err != nil {
		return nil, err
	}

	return job, nil
}

func (repo *JobRepositoryDb) Find(id string) (*domain.Job, error) {
	job := &domain.Job{}

	if err := repo.Db.Preload("Video").First(job, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if job.ID == "" {
		return nil, fmt.Errorf("video does not exist")
	}

	return job, nil
}

func (repo *JobRepositoryDb) Update(job *domain.Job) (*domain.Job, error) {
	if err := repo.Db.Save(job).Error; err != nil { return nil, err }
	return job, nil
}
