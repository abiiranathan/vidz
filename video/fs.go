package video

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"vidz/constants"
)

// Return the number of videos in the database or an error
func SyncDatabase(videoService VideoService, path string, refresh_db bool) (int, error) {
	if !refresh_db {
		return videoService.Count()
	}

	fmt.Printf("Crawling %s for videos in path: \n", path)
	// print all supported formats in a single line
	fmt.Printf("Supported formats: %s\n", constants.GetSupportedFormats())

	var videos []Video

	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		// ignore hidden directories
		if d.IsDir() && d.Name()[0] == '.' {
			return filepath.SkipDir

		}

		// ignore directories in ignoreDirs
		for _, ignoreDir := range constants.GetIgnoreDirs() {
			if d.IsDir() && d.Name() == ignoreDir {
				return filepath.SkipDir
			}
		}

		// ignore directories
		if d.IsDir() {
			return nil
		}

		// ignore non-video files
		if IsVideo(path, constants.GetSupportedFormats()) {
			fileinfo, err := os.Stat(path)

			if err != nil {
				return err
			}

			video := Video{
				Title:        d.Name(),
				Size:         int(fileinfo.Size()),
				Path:         path,
				Type:         GuessMimeType(path),
				LastModified: fileinfo.ModTime(),
			}

			// Ignore files smaller than 2MB
			if video.Size < 2*1024*1024 {
				return nil
			}

			videos = append(videos, video)
		}

		return nil
	})

	// save videos to database
	err := videoService.SaveMany(videos)
	if err != nil {
		return 0, err
	}

	return len(videos), nil

}

func IsVideo(path string, formats []string) bool {
	for _, format := range formats {
		if filepath.Ext(path) == "."+format {
			return true
		}
	}
	return false
}
