package readmorg

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"gopkg.in/yaml.v3"
)

var URLRoot string = "https://readm.org"

func getChapterDetails(link string) (ChapterData, error) {
	c := colly.NewCollector()

	chapterDetails := ChapterData{}
	chapterDetails.Link = link
	errorValue := ""

	c.OnHTML("div.ch-images.ch-image-container img", func(e *colly.HTMLElement) {
		imageLink := e.Attr("src")
		imageDetails := ImageData{}

		// /uploads/chapter_files/5618/66/1.jpg?v=12
		// parse the url, remove the query params and get the last element

		imageLinkParsed, _ := url.Parse(imageLink)
		imageLinkPath := strings.Split(imageLinkParsed.Path, "/")
		imageDetails.Name = imageLinkPath[len(imageLinkPath)-1]
		imageDetails.Link = URLRoot + imageLink
		

		chapterDetails.ImageLinks = append(chapterDetails.ImageLinks, imageDetails)
		chapterDetails.NumberOfImages++
	})

	c.OnHTML("h1.page-title a", func(e *colly.HTMLElement) {
		chapterDetails.MangaName = e.Text
	})

	c.OnHTML("h1.page-title span", func(e *colly.HTMLElement) {
		chapterDetails.Name = e.Text
	})

	c.OnHTML("div.media-date span", func(e *colly.HTMLElement) {
		chapterDetails.ReleaseDate = e.Text
	})

	c.OnError(func(e *colly.Response, err error) {
		if (e.StatusCode == 404) {
			errorValue = "404"
			return
		}
		fmt.Println("Something went wrong while getting chapter details:", err)
	})

	c.Visit(link)

	if errorValue == "404" {
		return chapterDetails, errors.New("404")
	}

	return chapterDetails, nil
}

func downloadImage(imageLink string, imageName string, pathRootMangaChapter string) {
	response, e := http.Get(imageLink)
	if e != nil {
		fmt.Fprintf(os.Stderr, "[DOWNLOAD FAILED] Error getting image %s: %v\n", imageLink, e)
		os.Exit(1)
	}

	defer response.Body.Close()
	imagePath := filepath.Join(pathRootMangaChapter, imageName)

	// create file
	file, err := os.Create(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[DOWNLOAD FAILED] Error creating file %s: %v\n", imagePath, err)
		os.Exit(1)
	}
	defer file.Close()

	// write to file
	_, err = io.Copy(file, response.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "[DOWNLOAD FAILED] Error writing to file %s: %v\n", imagePath, err)
		os.Exit(1)
	}

	fmt.Println("[DOWNLOAD SUCCESS] downloaded: ", imagePath)
}

func getMangaData(link string) MangaData {
	c := colly.NewCollector()

	mangaDetails := MangaData{}
	description := ""

	c.OnHTML("h1.page-title", func(e *colly.HTMLElement) {
		mangaDetails.Name = e.Text
	})

	c.OnHTML("div.series-summary-wrapper p", func(e *colly.HTMLElement) {
		description += strings.TrimSpace(e.Text)
	})

	c.OnHTML("div.series-summary-wrapper div.item a", func(e *colly.HTMLElement) {
		mangaDetails.Genres = append(mangaDetails.Genres, e.Text)
	})

	c.OnHTML(".table-episodes-title h6 a", func(e *colly.HTMLElement) {
		chapterLink := URLRoot + e.Attr("href")
		mangaDetails.ChapterLinks = append(mangaDetails.ChapterLinks, chapterLink)
	})

	c.OnError(func(e *colly.Response, err error) {
		fmt.Println("Something went wrong while getting manga data:", err)
	})

	c.Visit(link)

	mangaDetails.Description = description
	return mangaDetails
}

func writeChapterDetailsToFile(pathRootMangaChapter string, chapterDetails ChapterData) {

	f, err := os.Create(filepath.Join(pathRootMangaChapter, "chapter_details.yaml"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating chapter details file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(chapterDetails)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error encoding chapter details: %v\n", err)
		os.Exit(1)
	}
}

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}