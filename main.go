package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d videos in %s\n", count, folderToWatch)
	log.Printf("Listening on port %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}
