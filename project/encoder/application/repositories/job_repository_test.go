package repositories_test

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestJobRepositoryDbNewInsert(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()
	
	db := database.NewDbTest()
	defer db.Close()

	{	
		repo := repositories.VideoRepositoryDb{
			Db: db,
		}
		repo.Insert(video)
	}

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)

	repo := repositories.JobRepositoryDb{
		Db: db,
	}

	j, err := repo.Insert(job)

	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
}

func TestJobRepositoryDbNewFind(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()
	
	db := database.NewDbTest()
	defer db.Close()

	{	
		repo := repositories.VideoRepositoryDb{
			Db: db,
		}
		repo.Insert(video)
	}

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)

	repo := repositories.JobRepositoryDb{
		Db: db,
	}
	repo.Insert(job)
	j, err := repo.Find(job.ID)

	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.VideoID, video.ID)
}

func TestJobRepositoryDbNewUpdate(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()
	
	db := database.NewDbTest()
	defer db.Close()

	{	
		repo := repositories.VideoRepositoryDb{
			Db: db,
		}
		repo.Insert(video)
	}

	job, err := domain.NewJob("output_path", "Pending", video)
	require.Nil(t, err)

	repo := repositories.JobRepositoryDb{
		Db: db,
	}

	repo.Insert(job)

	job.Status = "Complete"
	repo.Update(job)
	
	j, err := repo.Find(job.ID)
	require.Nil(t, err)
	require.Equal(t, j.Status, job.Status)
}