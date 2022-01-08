package video

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	ID           int       `json:"id" gorm:"primary_key"`
	Title        string    `json:"title" gorm:"not null" validate:"required"`
	Size         int       `json:"size" gorm:"not null" validate:"required"`
	Path         string    `json:"path" gorm:"unique, not null" sql:"type:varchar(255)"`
	Type         string    `json:"type"`
	LastModified time.Time `json:"last_modified"`
}

type VideoDetail struct {
	Video Video `json:"video"`
	Next  int   `json:"next"`
	Prev  int   `json:"prev"`
}

type PaginatedVideos struct {
	Videos []Video `json:"videos"`
	Page   int     `json:"page"`
	Next   int     `json:"next"`
	Prev   int     `json:"prev"`
}

type VideoService interface {
	GetAll() ([]Video, error)
	Get(id int) (Video, error)
	Save(video Video) error
	SaveMany(videos []Video) error
	Delete(id int) error
	FindByPath(path string) (Video, error)
	Count() (int, error)
	GetNextID(id int) (int, error)
	GetPrevID(id int) (int, error)
}

// implement VideoServiceImpl
type VideoServiceImpl struct {
	db *gorm.DB
}

// implement VideoService Count
func (v VideoServiceImpl) Count() (int, error) {
	var count int64
	err := v.db.Model(&Video{}).Count(&count).Error
	return int(count), err
}

// implement VideoService GetAll
func (v VideoServiceImpl) GetAll() ([]Video, error) {
	var videos []Video
	err := v.db.Find(&videos).Error
	if err != nil {
		return nil, err

	}
	return videos, nil
}

// implement VideoService Get
func (v VideoServiceImpl) Get(id int) (Video, error) {
	var video Video
	err := v.db.First(&video, id).Error
	if err != nil {
		return Video{}, err
	}

	return video, nil
}

// implement VideoService Save
func (v VideoServiceImpl) Save(video Video) error {
	err := v.db.Save(&video).Error
	if err != nil {
		return err
	}

	return nil
}

// Delete all videos from the database that are not in the given list
// of videos (videosToKeep).
func (v VideoServiceImpl) SaveMany(videos []Video) error {
	// fetch all videos from the database
	var existingVideos []Video
	err := v.db.Find(&existingVideos).Error
	if err != nil {
		return err
	}

	// create a map of existing videos
	existingVideosMap := make(map[string]Video)
	for _, video := range existingVideos {
		existingVideosMap[video.Path] = video
	}

	// create a map of videos to be saved
	videosToBeSaved := make(map[string]Video)
	for _, video := range videos {
		videosToBeSaved[video.Path] = video
	}

	// delete videos that are not in the videosToBeSaved map
	for _, video := range existingVideos {
		if _, ok := videosToBeSaved[video.Path]; !ok {
			// permanently delete the video
			err := v.db.Unscoped().Delete(&video).Error
			if err != nil {
				return err
			}
		}
	}

	// save videos that are in the videosToBeSaved map
	for _, video := range videos {
		if _, ok := existingVideosMap[video.Path]; !ok {
			err := v.db.Create(&video).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// implement VideoService Delete
func (v VideoServiceImpl) Delete(id int) error {
	var video Video
	err := v.db.First(&video, id).Error
	if err != nil {
		return err
	}

	err = v.db.Unscoped().Delete(&video).Error
	if err != nil {
		return err
	}
	// delete the file
	os.Remove(video.Path)
	return nil
}

// implement VideoService FindByPath
func (v VideoServiceImpl) FindByPath(path string) (Video, error) {
	var video Video
	err := v.db.Where("path = ?", path).First(&video).Error
	if err != nil {
		return Video{}, err
	}

	return video, nil
}

// Returns the next video id after the given id
func (v VideoServiceImpl) GetNextID(id int) (int, error) {
	var nextVideo Video
	err := v.db.Where("id > ?", id).First(&nextVideo).Error

	if err != nil {
		return 0, err
	}

	return nextVideo.ID, nil
}

// Returns the prev video id after the given id
func (v VideoServiceImpl) GetPrevID(id int) (int, error) {
	var prevVideo Video
	err := v.db.Where("id < ?", id).Last(&prevVideo).Error

	if err != nil {
		return 0, err
	}

	return prevVideo.ID, nil

}

// Returns the human readable size of a file
func (v Video) HSize() string {
	return HumanReadableSize(v.Size)
}

// concrete implementation of a function that returns a human readable size of a file
func HumanReadableSize(size int) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/1024/1024)
	} else {
		return fmt.Sprintf("%.2f GB", float64(size)/1024/1024/1024)
	}
}

// returns the mime type of a file based on its extension
func GuessMimeType(path string) string {
	ext := filepath.Ext(path)

	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".avi":
		return "video/x-msvideo"
	case ".mkv":
		return "video/x-matroska"
	case ".mov":
		return "video/quicktime"
	case ".flv":
		return "video/x-flv"
	case ".wmv":
		return "video/x-ms-wmv"
	case ".mpeg":
		return "video/mpeg"
	case ".m4v":
		return "video/x-m4v"
	case ".3gp":
		return "video/3gpp"
	case ".ts":
		return "video/mp2t"
	default:
		return "application/octet-stream"
	}
}

// Create and returns a new VideoService implementation
func NewVideoService(db *gorm.DB) VideoService {
	return &VideoServiceImpl{db}
}
