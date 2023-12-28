package main

import (
	mc "downloader/src/mangaclash"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

func getBaseURL(link string) (string, error) {
	url, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	return url.Scheme + "://" + url.Host, nil
}

// func logger(filePath string, text string) {
// 	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)

// }

func main() {
	{
		// setup logger
		f, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		log.SetOutput(f)
		defer f.Close()
	}

	// link := "https://mangaclash.com/manga/shadowless-night/"
	link := *flag.String("link", "https://mangaclash.com/manga/shadowless-night/", "link to manga")
	flag.Parse()

	downloadPath, err := filepath.Abs(filepath.Join("test"))
	if err != nil {
		os.Exit(1)
	}

	baseURL, err := getBaseURL(link)
	if err != nil {
		log.Fatal(err)
	}

	switch baseURL {
	case "https://mangaclash.com":
		mc.Download(link, downloadPath)
	default:
		// exit with error
		fmt.Println("Invalid URL")
		os.Exit(1)
	}
}