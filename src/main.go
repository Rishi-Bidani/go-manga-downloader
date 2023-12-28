package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

type chapterData struct {
	name string
	link string
	releaseDate string
}

type chapterImage struct {
	chapterName string
	imageLink string
	imageName string
}

type mangaData struct {
	name string
	genres []string
	description string
}

func getChapterLinks(link string) ([]chapterData, mangaData){
	c := colly.NewCollector()

	// chapterlink array
	chapterLinkArr := []chapterData{}

	c.OnHTML("li.wp-manga-chapter", func(e *colly.HTMLElement) {
		// within the link, find direct anchor tag and span.chapter-release-date
		link := e.ChildAttr("a", "href")
		name := e.ChildText("a")
		releaseDate := e.ChildText("span.chapter-release-date")
		chapterLinkArr = append(chapterLinkArr, chapterData{name, link, releaseDate})
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
		os.Exit(1)
	})

	mangaDetails := mangaData{}

	c.OnHTML("h1", func(e *colly.HTMLElement) {
		mangaDetails.name = strings.TrimSpace(e.Text)
	})

	c.OnHTML("summary__content p", func(e *colly.HTMLElement) {
		mangaDetails.description = e.Text
	})

	c.OnHTML("div.genres-content a", func(e *colly.HTMLElement) {
		mangaDetails.genres = append(mangaDetails.genres, e.Text)
	})


	c.Visit(link)

	return chapterLinkArr, mangaDetails
}

func getImageLinks(chapterName string, chapterLink string) []chapterImage{
	c := colly.NewCollector()
	chapterImageArr := []chapterImage{}

	c.OnHTML("div.reading-content .page-break img", func(e *colly.HTMLElement) {
		// find all images in the div
		imageLink := strings.TrimSpace(e.Attr("data-src"))
		// find the image number
		imageName := strings.TrimSpace(e.Attr("id"))
		// add file extension to image name using the link
		imageName = imageName + filepath.Ext(imageLink)

		// append to array
		chapterImageArr = append(chapterImageArr, chapterImage{chapterName, imageLink, imageName})
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
		os.Exit(1)
	})
	
	c.Visit(chapterLink)

	return chapterImageArr
}

func downloadImage(rootPath string, imageData chapterImage) {
	chapterName := imageData.chapterName
	imageLink := imageData.imageLink
	imageName := imageData.imageName

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

	link := "https://mangaclash.com/manga/shadowless-night/"

	chapterLinks, mangaDetails := getChapterLinks(link)

	masterImageArr := []chapterImage{}

	// // synchronous
	// for _, cd := range chapterLinks {
	// 	chapterImageArr := getImageLinks(cd.name, cd.link)
	// 	// append to masterImageArr
	// 	masterImageArr = append(masterImageArr, chapterImageArr...)
	// }

	// asynchronous
	var wg sync.WaitGroup
	for _, cd := range chapterLinks {
		wg.Add(1)
		// create a go routine for each chapter
		go func(_cd chapterData) {
			// get image links for each chapter
			chapterImageArr := getImageLinks(_cd.name, _cd.link)
			// append to masterImageArr
			masterImageArr = append(masterImageArr, chapterImageArr...)
			// listout chapter name, number of images and image links to log file
			log.Println("Chapter Name: ", _cd.name, "Number of Images: ", len(chapterImageArr), "Image Links: ", chapterImageArr)

			defer wg.Done()
		}(cd)

	}
	wg.Wait()
	// print length of masterImageArr
	fmt.Println(len(masterImageArr))

	// download images =================
	// create a folder for the manga
	mangaName := mangaDetails.name
	rootPath, _ := filepath.Abs(filepath.Join("test", mangaName))
	err := os.MkdirAll(rootPath, os.ModePerm)
	if err != nil {
		fmt.Errorf("error creating manga folder: %w", err)
		os.Exit(1)
	}

	// create a folder for each chapter
	for _, cd := range chapterLinks {
		chapterName := cd.name
		chapterPath, _ := filepath.Abs(filepath.Join(rootPath, chapterName))
		err := os.MkdirAll(chapterPath, os.ModePerm)
		if err != nil {
			fmt.Errorf("error creating chapter folder: %w", err)
			os.Exit(1)
		}
	}

	var wg2 sync.WaitGroup
	// download images to the chapter folder
	for _, ci := range masterImageArr {
		wg2.Add(1)
		go func(_ci chapterImage) {
			downloadImage(rootPath, _ci)
			defer wg2.Done()
		}(ci)
	}
	wg2.Wait()
}