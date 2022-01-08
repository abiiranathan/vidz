package routes

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"vidz/video"
)

// filterVideos returns a slice of videos that match the query
func filterVideos(videos *[]video.Video, filter string) []video.Video {
	filter = strings.ToLower(filter)
	filteredVideos := []video.Video{}
	for _, video := range *videos {
		title := strings.ToLower(video.Title)
		if strings.Contains(title, filter) {
			filteredVideos = append(filteredVideos, video)
		}
	}

	return filteredVideos
}

// not found handler
func NotFound(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "404.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// VideoListRoute returns a handler for the video list
func Index(videoService video.VideoService, tmpl *template.Template) http.HandlerFunc {
	page_size := 10

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			id := r.URL.Query().Get("id")
			videoID, err := strconv.Atoi(id)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// find video by id from the database
			err = videoService.Delete(videoID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
			return

		} else {
			// Get all videos from the database
			videos, err := videoService.GetAll()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			query := r.URL.Query()
			q := query.Get("query")
			page := query.Get("page")

			if page == "" {
				page = "1"
			}

			pageNum, err := strconv.Atoi(page)

			if err != nil {
				pageNum = 1
			}

			if q == "" {
				err := tmpl.ExecuteTemplate(w, "index.html", PaginateVideos(videos, pageNum, page_size))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

			} else {
				filteredVideos := filterVideos(&videos, q)
				err := tmpl.ExecuteTemplate(w, "index.html", PaginateVideos(filteredVideos, pageNum, page_size))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		}
	}
}

// VideoRoute returns a handler for the video page
func VideoDetailRoute(videoService video.VideoService, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("id")
		videoID, err := strconv.Atoi(id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// find video by id
		v, err := videoService.Get(videoID)
		if err != nil {
			// if not found, redirect to NotFound handler
			http.Redirect(w, r, "/404", http.StatusFound)
			return
		}

		nextId, err := videoService.GetNextID(videoID)

		if err != nil {
			nextId = 0
		}

		prevId, err := videoService.GetPrevID(videoID)

		if err != nil {
			prevId = 0
		}

		videoDetail := video.VideoDetail{
			Video: v,
			Next:  nextId,
			Prev:  prevId,
		}

		err = tmpl.ExecuteTemplate(w, "play.html", videoDetail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Custom handler that streams a video to the client
func MediaStreamHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the filename from the request path
		filename := r.URL.Path

		// open the file
		v, err := os.Open(filename)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// close the file when the function returns
		defer v.Close()

		// get the file info
		fileinfo, err := v.Stat()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get the file size
		filesize := fileinfo.Size()

		// set the headers
		w.Header().Set("Content-Type", video.GuessMimeType(filename))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", filesize))
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", fileinfo.Name()))

		// stream the file to the client
		http.ServeContent(w, r, fileinfo.Name(), fileinfo.ModTime(), v)
	}
}

func PaginateVideos(videos []video.Video, page, pageSize int) video.PaginatedVideos {
	if pageSize == 0 || pageSize > len(videos) {
		pageSize = len(videos)
	}

	if page > 0 {
		page--
	}

	if page*pageSize > len(videos) {
		page = 0
	}

	if (page*pageSize)+pageSize > len(videos) {
		pageSize = len(videos) - page*pageSize
	}

	return video.PaginatedVideos{
		Videos: videos[page*pageSize : page*pageSize+pageSize],
		Page:   page + 1,
		Next:   page + 2,
		Prev:   page,
	}
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ip string

		if r.Header.Get("X-Forwarded-For") == "" {
			ip = r.RemoteAddr
		} else {
			ip = r.Header.Get("X-Forwarded-For")
		}

		// sort date and time
		now := time.Now().Format("2006-01-02 15:04 PM")

		log.Printf("%s %s %s %s %s", now, r.Method, r.URL, r.Proto, ip)
		next.ServeHTTP(w, r)
	})
}

// Create a new serveMux with the routes
func NewServeMux(videoService video.VideoService, tmpl *template.Template, static embed.FS) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", LoggerMiddleware(Index(videoService, tmpl)))

	// serve static files
	mux.Handle("/static/", http.FileServer(http.FS(static)))

	// server media files (videos)
	VideoStreamHandler := http.HandlerFunc(MediaStreamHandler())
	mux.Handle("/media/", http.StripPrefix("/media", VideoStreamHandler))

	// video detail handler
	VideoDetailHandler := http.HandlerFunc(VideoDetailRoute(videoService, tmpl))
	mux.Handle("/videoplayer", VideoDetailHandler)

	// handle 404
	mux.HandleFunc("/404", NotFound(tmpl))

	return mux
}
