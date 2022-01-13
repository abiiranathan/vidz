package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"vidz/db"
	"vidz/routes"
	"vidz/video"
)

//go:embed static/*
var static embed.FS

//go:embed static/templates
var templates embed.FS

var tmpl *template.Template
var dbPath string
var timeout = 10 * time.Second

func init() {
	tmpl = template.Must(template.ParseFS(templates, "static/templates/*.html"))
	path, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dir := filepath.Dir(path)
	abspath := filepath.Join(dir, "videos.db")
	if err != nil {
		panic(err)
	}

	dbPath = abspath
}

func startServer(handler http.Handler, port int) {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      routes.LoggerMiddleware(handler),
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	done := make(chan error, 1)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		ctx := context.Background()
		var cancel context.CancelFunc
		if timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		done <- srv.Shutdown(ctx)
	}()

	// Recovers from panic

	log.Printf("Listening on port %d\n", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	log.Println("Server stopped successfully")
}

func main() {
	defaultFolder, _ := os.UserHomeDir()
	dirname := flag.String("dirname", defaultFolder, "The directory to scan for videos")
	port := flag.Int("port", 8080, "The port to listen on")
	dbflag := flag.String("db", dbPath, "The database file to use")
	refresh_db := flag.Bool("refresh_db", false, "Refresh the database")

	flag.Parse()

	folderToWatch, err := filepath.Abs(*dirname)

	if err != nil {
		log.Fatalf("Could not covert %s to absolute path: %s ", folderToWatch, err)
	}

	videoService := video.NewVideoService(db.New(*dbflag))
	handler := routes.NewServeMux(&videoService, tmpl, static)

	count, err := video.SyncDatabase(&videoService, folderToWatch, *refresh_db)
	fmt.Printf("Found %d videos in %s\n", count, folderToWatch)

	if err != nil {
		log.Fatalf("Error syncing db: %v", err.Error())
	}

	startServer(handler, *port)
}
