package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"strings"

	"github.com/oz/osdb"
)

var extensions = []string{".mp4", ".avi", ".mkv"}
var languages = []string{"eng"}

func isExtOk(file string) bool {
	for _, extension := range extensions {
		if extension == filepath.Ext(file) {
			return true
		}
	}
	return false
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s", name, elapsed)
}

func main() {
	defer timeTrack(time.Now(), "autosub")

	//Define the base directory for the search
	serieDir := "."
	if len(os.Args) > 1 {
		serieDir = os.Args[1]
	}

	fmt.Println("Autosub starting ...")

	var subNotFound int

	c, err := osdb.NewClient()
	if err != nil {
		log.Fatalln("Cannot initiate opensubtitle client")
	}
	if err = c.LogIn("", "", ""); err != nil {
		log.Fatalln("Canot connect to the opensubtitle api")
	}
	fmt.Println("Connection to opensubtitle OK ...")
	videoFiles := getAllVideoFiles(serieDir)
	fmt.Printf("Searching subtitles for %d video files\n\n", len(videoFiles))
	for _, video := range videoFiles {
		srtPath := video[:strings.LastIndex(video, ".")] + ".srt"
		subs, err := c.FileSearch(video, languages)
		if err != nil {
			log.Println(err)
		}
		if len(subs) < 1 {
			log.Println("Cannot find subtitle for : ", path.Base(video))
			subNotFound++
		} else {

			if err := c.DownloadTo(&subs[0], srtPath); err != nil {
				log.Println(err)
			}
		}

	}
	fmt.Println()
	fmt.Println(len(videoFiles)-subNotFound, "subtitles found")
	fmt.Println(subNotFound, "subtitles not found")

}

func getAllVideoFiles(serieDir string) []string {
	var videoFiles []string
	src, err := os.Open(serieDir)
	defer src.Close()
	if err != nil {
		log.Fatalf("Could not open serieDir folder : %s", err)
	}
	if filemode, _ := src.Stat(); !filemode.IsDir() {
		log.Fatalf("%s is not a directory", serieDir)
	}
	files, err := src.Readdir(-1)
	if err != nil {
		log.Fatalf("Readding %s failed : %s", serieDir, err)
	}
	for _, file := range files {
		filename := file.Name()
		if isExtOk(filename) {
			videoFiles = append(videoFiles, filepath.Join(serieDir, filename))
		} else if file.IsDir() {
			videoFiles = append(videoFiles, getAllVideoFiles(filepath.Join(serieDir, filename))...)
		}
	}
	return videoFiles
}
