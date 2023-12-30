package main

import (
	mc "downloader/src/mangaclash"
	"flag"
	"fmt"
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

func main() {
	// link := "https://mangaclash.com/manga/shadowless-night/"
	_link := flag.String("link", "https://mangaclash.com/manga/shadowless-night/", "link to manga")
	_rootFolder := flag.String("folder", "test", "root folder to download manga")
	_downloadSingleChapter := flag.Bool("single", false, "download single chapter")
	flag.Parse()

	link := *_link
	rootFolder := *_rootFolder
	downloadSingleChapter := *_downloadSingleChapter

	pathRoot, err := filepath.Abs(filepath.Join(rootFolder))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting absolute path: %v\n", err)
		os.Exit(1)
	}

	baseURL, err := getBaseURL(link)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing url: %v\n", err)
		os.Exit(1)
	}

	// =================================================================================
	// switch case for different baseURL
	// =================================================================================
	switch baseURL {
		case "https://mangaclash.com":
			if downloadSingleChapter {
				// download single chapter
				mc.DownloadChapter(pathRoot, link)
			} else {
				mc.Download(link, pathRoot)
			}
		default:
			// exit with error
			fmt.Println("Invalid URL")
			os.Exit(1)
	}
	// =================================================================================
}