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
	// get binary path
	path, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	abspath := filepath.Join(path, "videos.db")
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

	// make sure db exists
	if _, err := os.Stat(*dbflag); os.IsNotExist(err) {
		fmt.Println("Creating database file: ", *dbflag)
		_, err = os.Create(*dbflag)
		if err != nil {
			panic("Could not create database file: " + err.Error())
		}
	}

	if *dirname == "" {
		panic("No folder to watch specified in arguments and environment variable for HOME not set")
	}

	// convert to absolute path
	folderToWatch, err := filepath.Abs(*dirname)
	if err != nil {
		log.Fatalf("Could not covert %s to absolute path: %s ", folderToWatch, err)
	}

	// create a new video service
	videoService := video.NewVideoService(db.New(*dbflag))
	handler := routes.NewServeMux(videoService, tmpl, static)

	// start watching the folder
	count, err := video.SyncDatabase(videoService, folderToWatch, *refresh_db)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d videos in %s\n", count, folderToWatch)

	// start the server
	log.Printf("Listening on port %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}
