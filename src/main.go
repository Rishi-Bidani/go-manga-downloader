package main

import (
	mc "downloader/src/mangaclash"
	rd "downloader/src/readmorg"
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
	_link := flag.String("link", "https://readm.org/manga/owari-no-seraph/", "link to manga")
	_rootFolder := flag.String("folder", "test", "root folder to download manga")
	_downloadSingleChapter := flag.Bool("single", false, "download single chapter boolean. If true, link must be a chapter link")
	_downloadStart := flag.Int("start", -1, "start chapter. Do not provide if attempting to download single or all chapters. Link must be a manga link")
	_downloadEnd := flag.Int("end", -1, "end chapter. Do not provide if attempting to download single or all chapters. Link must be a manga link")
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

	var start, end int
	downloadChapterRange := false
	if *_downloadStart != -1 && *_downloadEnd != -1 {
		start = *_downloadStart
		end = *_downloadEnd
		downloadChapterRange = true
	}

	// =================================================================================
	// switch case for different baseURL
	// =================================================================================
	switch baseURL {
		case "https://mangaclash.com":
			if downloadSingleChapter {
				// download single chapter
				mc.DownloadChapter(pathRoot, link)
			} else if downloadChapterRange {
				// download chapter range
				mc.DownloadChapterRange(link, pathRoot, start, end)
			} else {
				mc.Download(link, pathRoot)
			}

		case "https://readm.org":
			if downloadSingleChapter {
				rd.DownloadChapter(pathRoot, link)
			} else {
				rd.Download(pathRoot, link)
			}
		
		default:
			// exit with error
			fmt.Println("Invalid URL")
			os.Exit(1)
	}
	// =================================================================================
}