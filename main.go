package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	http.HandleFunc("/", fish)

	port := ":8080"

	log.Infof("Listening in port %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("error : %s", err.Error())
	}
}

func fish(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())

	imageList := getImageList()
	imageListLen := len(imageList)

	randInt := rand.Intn(imageListLen - 1)
	filepath := fmt.Sprintf("fish/%s", imageList[randInt])

	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("open file error : %s", err.Error())
	}
	defer f.Close()

	w.Header().Set("Content-Type", "image/jpeg")
	_, err = io.Copy(w, f)
	if err != nil {
		log.Errorf("error copying image : %s", err.Error())
	} else {
		log.Infof("fish pic '%s'", imageList[randInt])
	}
}

func getImageList() []string {
	var imageList []string

	files, err := ioutil.ReadDir("fish/")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		imageList = append(imageList, f.Name())
	}
	return imageList
}
