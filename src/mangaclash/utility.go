package mangaclash

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly"
)


func getChapterLinksAndMangaDetails(link string) ([]ChapterData, MangaData){
	c := colly.NewCollector()

	// chapterlink array
	chapterLinkArr := []ChapterData{}

	c.OnHTML("li.wp-manga-chapter", func(e *colly.HTMLElement) {
		// within the link, find direct anchor tag and span.chapter-release-date
		link := e.ChildAttr("a", "href")
		name := e.ChildText("a")
		releaseDate := e.ChildText("span.chapter-release-date")
		chapterLinkArr = append(chapterLinkArr, ChapterData{name, link, releaseDate})
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
		os.Exit(1)
	})

	mangaDetails := MangaData{}

	c.OnHTML("h1", func(e *colly.HTMLElement) {
		mangaDetails.Name = strings.TrimSpace(e.Text)
	})

	c.OnHTML("summary__content p", func(e *colly.HTMLElement) {
		mangaDetails.Description = e.Text
	})

	c.OnHTML("div.genres-content a", func(e *colly.HTMLElement) {
		mangaDetails.Genres = append(mangaDetails.Genres, e.Text)
	})


	c.Visit(link)

	return chapterLinkArr, mangaDetails
}

func getImageLinks(chapterName string, chapterLink string) []ChapterImage{
	c := colly.NewCollector()
	ChapterImageArr := []ChapterImage{}

	c.OnHTML("div.reading-content .page-break img", func(e *colly.HTMLElement) {
		// find all images in the div
		imageLink := strings.TrimSpace(e.Attr("data-src"))
		// find the image number
		imageName := strings.TrimSpace(e.Attr("id"))
		// add file extension to image name using the link
		imageName = imageName + filepath.Ext(imageLink)

		// append to array
		ChapterImageArr = append(ChapterImageArr, ChapterImage{chapterName, imageLink, imageName})
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
		os.Exit(1)
	})
	
	c.Visit(chapterLink)

	return ChapterImageArr
}

func downloadImage(rootPath string, imageData ChapterImage) {
	chapterName := imageData.ChapterName
	imageLink := imageData.ImageLink
	imageName := imageData.ImageName

	imagePath := filepath.Join(rootPath, chapterName, imageName)

	// download file
	response, e := http.Get(imageLink)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	// open a file for writing
	file, err := os.Create(imagePath)
	if err != nil {
		log.Fatal("Error creating file: ", err)
	}
	defer file.Close()

	// write the body to file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Download Success!", imagePath)
}