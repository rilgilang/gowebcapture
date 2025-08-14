package repositories

import (
	"context"
	"github.com/google/uuid"
	consts "go/src/github.com/rilgilang/gowebcapture/const"
	"go/src/github.com/rilgilang/gowebcapture/entities"
	"gorm.io/gorm"
	"time"
)

type VideoRepository interface {
	FetchAll(ctx context.Context) ([]entities.Video, error)
	SaveProcessedVideoURL(ctx context.Context, uniqueId, videoUrl, status string) (*entities.Video, error)
}

type videoRepository struct {
	db *gorm.DB
}

func NewVideoRepo(db *gorm.DB) VideoRepository {
	return &videoRepository{
		db: db,
	}
}

func (r *videoRepository) FetchAll(ctx context.Context) ([]entities.Video, error) {
	videos := []entities.Video{}

	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT 
		    id,
		    user_id,
		    url,
		    created_at,
		    updated_at
    	From videos`).Rows()

	if err != nil {
		if err.Error() == consts.SqlNoRow {
			return nil, nil
		}
		return nil, err
	}

	for rows.Next() {
		video := entities.Video{}
		if err = rows.Scan(
			&video.ID,
			&video.UserID,
			&video.URL,
		); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (r *videoRepository) SaveProcessedVideoURL(ctx context.Context, uniqueId, videoUrl, status string) (*entities.Video, error) {
	video := entities.Video{}

	id := uuid.New().String()
	// Use Model(&room) and Where before calling Updates
	err := r.db.WithContext(ctx).
		Model(&video).
		Create(map[string]interface{}{
			"id":         id,
			"unique_id":  uniqueId,
			"url":        videoUrl,
			"status":     status,
			"created_at": time.Now(),
			"updated_at": time.Now(),
		}).Error

	if err != nil {
		if err.Error() == consts.SqlNoRow {
			return nil, nil
		}
		return nil, err
	}

	err = r.db.WithContext(ctx).First(&video, "id = ?", id).Error

	if err != nil {
		if err.Error() == consts.SqlNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &video, nil
}
