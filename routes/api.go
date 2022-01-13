package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"vidz/video"
)

func VideoListApi(videoService video.VideoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		videos, err := videoService.GetAll()
		query := r.URL.Query()
		title := query.Get("title")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if title != "" {
			filteredVideos := filterVideos(&videos, title)
			err = json.NewEncoder(w).Encode(filteredVideos)
		} else {
			err = json.NewEncoder(w).Encode(videos)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func VideoDetailApi(videoService video.VideoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the id from the query string
		query := r.URL.Query()
		id := query.Get("id")
		videoID, err := strconv.Atoi(id)

		if err != nil {
			http.Error(w, "Error parsing video id", http.StatusBadRequest)
			return
		}

		// find video by id
		v, err := videoService.Get(videoID)
		if err != nil {
			http.Error(w, fmt.Sprintf("No video mathes id: %d", videoID), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(v)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
