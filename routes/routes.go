package routes

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"vidz/video"
)

func Index(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func MediaStreamHandler(videoService video.VideoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Path
		v, err := os.Open(filename)
		if err != nil {
			// Video was deleted, remove it from database
			vid, err := videoService.FindByPath(filename)
			if err != nil {
				videoService.Delete(vid.ID)
			}

			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		defer v.Close()
		fileinfo, err := v.Stat()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filesize := fileinfo.Size()
		w.Header().Set("Content-Type", video.GuessMimeType(filename))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", filesize))
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", fileinfo.Name()))
		http.ServeContent(w, r, fileinfo.Name(), fileinfo.ModTime(), v)
	}
}

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

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}

		var ip string

		if r.Header.Get("X-Forwarded-For") == "" {
			ip = r.RemoteAddr
		} else {
			ip = r.Header.Get("X-Forwarded-For")
		}

		next.ServeHTTP(w, r)
		log.Printf("%s %d %s %s %s", r.Method, recorder.Status, r.URL, r.Proto, ip)
	})
}

func NewServeMux(videoService *video.VideoService, tmpl *template.Template, static embed.FS) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", LoggerMiddleware(Index(tmpl)))
	mux.Handle("/static/", http.FileServer(http.FS(static)))
	mux.Handle("/media/", http.StripPrefix("/media", http.HandlerFunc(MediaStreamHandler(*videoService))))
	mux.HandleFunc("/api/videos", VideoListApi(*videoService))
	mux.HandleFunc("/api/videos/detail", VideoDetailApi(*videoService))

	return mux
}
