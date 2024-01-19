package main

import (
	plugin "downloader/src/plugins"
	mc "downloader/src/plugins/mangaclash"
	rd "downloader/src/plugins/readmorg"
	"flag"
	"fmt"
	"math"
	"net/url"
	"os"
	"path/filepath"
)

// getBaseURL returns the base url of the link
// e.g. https://example.com/abc/def/ghi returns https://example.com
// Error is returned if link is invalid i.e. it can't be parsed by url.Parse()
func getBaseURL(link string) (string, error) {
	url, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	return url.Scheme + "://" + url.Host, nil
}

// getPlugin returns the plugin for the website
// currently supported websites:
// - https://mangaclash.com
// - https://readm.org
// If the website is not supported, the program will exit with error
func getPlugin(baseUrl string) plugin.IMangaDownloader {
	switch baseUrl {
		case "https://mangaclash.com":
			return &mc.MangaClash{}
		
		case "https://readm.org":
			return &rd.ReadmOrg{}
		
		default:
			// exit with error
			fmt.Println("Invalid URL")
			os.Exit(1)
	}
	return nil
}


func main() {
	// link := "https://mangaclash.com/manga/shadowless-night/"
	// ==========================================================================================================================================================
	// flags
	// ==========================================================================================================================================================
	_link := flag.String("link", "", "link to manga")
	_rootFolder := flag.String("folder", "test", "root folder to download manga")
	_downloadSingleChapter := flag.Bool("single", false, "download single chapter boolean. If true, link must be a chapter link")
	_downloadStart := flag.Int("start", 0, "start chapter. Do not provide if attempting to download single or all chapters. Link must be a manga link")
	_downloadEnd := flag.Int("end", int(math.Inf(-1)), "end chapter. Do not provide if attempting to download single or all chapters. Link must be a manga link")
	flag.Parse()

	link := *_link
	rootFolder := *_rootFolder
	downloadSingleChapter := *_downloadSingleChapter
	// ==========================================================================================================================================================

	pathRoot, err := filepath.Abs(filepath.Join(rootFolder))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting absolute path: %v\n", err)
		os.Exit(1)
	}

	// check if link is valid and get base url to determine which plugin to use
	baseURL, err := getBaseURL(link)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing url: %v\n", err)
		os.Exit(1)
	}

	// checking if user is downloading a specific range of chapters
	var start, end int
	downloadChapterRange := false
	if *_downloadStart != 0 || *_downloadEnd != int(math.Inf(-1)) {
		start = *_downloadStart
		end = *_downloadEnd
		downloadChapterRange = true
	}

	// plugin for different website
	plugin := getPlugin(baseURL)
	if plugin == nil {
		fmt.Fprintf(os.Stderr, "error getting plugin\n")
		os.Exit(1)
	}

	// ==============================================================
	// download manga
	// ==============================================================
	if downloadSingleChapter {
		// download single chapter
		plugin.DownloadChapter(pathRoot, link)
	} else if downloadChapterRange {
		// download chapter range
		plugin.DownloadChapterRange(pathRoot, link, start, end)
	} else {
		plugin.Download(pathRoot, link)
	}
	// ==============================================================
}