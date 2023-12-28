package mangaclash

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func Download(link string, downloadPath string){
	chapterLinks, mangaDetails := getChapterLinks(link)

	masterImageArr := []ChapterImage{}

	// asynchronous
	var wg sync.WaitGroup
	for _, cd := range chapterLinks {
		wg.Add(1)
		// create a go routine for each chapter
		go func(_cd ChapterData) {
			// get image links for each chapter
			chapterImageArr := getImageLinks(_cd.Name, _cd.Link)
			// append to masterImageArr
			masterImageArr = append(masterImageArr, chapterImageArr...)
			// listout chapter name, number of images and image links to log file
			log.Println("Chapter Name: ", _cd.Name, "Number of Images: ", len(chapterImageArr), "Image Links: ", chapterImageArr)

			defer wg.Done()
		}(cd)

	}
	wg.Wait()
	// print length of masterImageArr
	fmt.Println(len(masterImageArr))

	// download images =================
	// create a folder for the manga
	mangaName := mangaDetails.Name
	rootPath, _ := filepath.Abs(filepath.Join(downloadPath, mangaName))
	err := os.MkdirAll(rootPath, os.ModePerm)
	if err != nil {
		fmt.Errorf("error creating manga folder: %w", err)
		os.Exit(1)
	}

	// create a folder for each chapter
	for _, cd := range chapterLinks {
		chapterName := cd.Name
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
		go func(_ci ChapterImage) {
			downloadImage(rootPath, _ci)
			defer wg2.Done()
		}(ci)
	}
	wg2.Wait()
}