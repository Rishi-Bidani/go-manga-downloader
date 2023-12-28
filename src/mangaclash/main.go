package mangaclash

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type MangaMetaData struct {
	Name string
	Description string
	Genres []string
	NumberOfChapters int
	Chapters [] struct {
		Name string
		ReleaseDate string
		Link string
		NumberOfImages int
		Images [] struct {
			ImageLink string
			ImageName string
		}
	}
}

func mangaMetaData (filePath string, mangaName string, mangaDescription string, mangaGenres []string, chapterArr []ChapterData, imageArr []ChapterImage){
	// this function will write a mangaName.yaml file in the root folder
	// the yaml file will contain the following information:
	// - mangaName
	// - mangaDescription
	// - mangaGenres
	// - number of chapters
	// - chapterName
	// 		- chapter release date
	// 		- chapter link
	// 		- number of images
	// 		- image links array

	mangaMeta := MangaMetaData{}
	mangaMeta.Name = mangaName
	mangaMeta.Description = mangaDescription
	mangaMeta.Genres = mangaGenres
	mangaMeta.NumberOfChapters = len(chapterArr)

	for _, chapter := range chapterArr {
		chapterData := struct {
			Name string
			ReleaseDate string
			Link string
			NumberOfImages int
			Images [] struct {
				ImageLink string
				ImageName string
			}
		}{}
		chapterData.Name = chapter.Name
		chapterData.ReleaseDate = chapter.ReleaseDate
		chapterData.Link = chapter.Link
		chapterData.NumberOfImages = 0

		for _, image := range imageArr {
			if image.ChapterName == chapter.Name {
				chapterData.NumberOfImages++
				chapterData.Images = append(chapterData.Images, struct {
					ImageLink string
					ImageName string
				}{image.ImageLink, image.ImageName})
			}
		}
		mangaMeta.Chapters = append(mangaMeta.Chapters, chapterData)
	
	}

	// write to yaml file
	f, err := os.OpenFile(filepath.Join(filePath, mangaName + ".yaml"), os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening yaml file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(mangaMeta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error encoding yaml file: %v\n", err)
		os.Exit(1)
	}
}

func Download(link string, downloadPath string){
	/*
		Download function will download the manga from the link provided
		It will download all the chapters and images
	*/
	chapterLinks, mangaDetails := getChapterLinks(link)
	
	// create a folder for the manga
	mangaName := mangaDetails.Name
	rootPath, _ := filepath.Abs(filepath.Join(downloadPath, mangaName))
	err := os.MkdirAll(rootPath, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating manga folder: %v\n", err)
		os.Exit(1)
	}

	{
		// =================================================================================
		// setup logger
		// =================================================================================
		f, err := os.OpenFile(filepath.Join(rootPath, "log.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening log file: %v\n", err)
			os.Exit(1)
		}
		log.SetOutput(f)
		defer f.Close()
		// =================================================================================
	}
	
	// =====================================================================================
	// get image links for each chapter
	// =====================================================================================
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
	// =====================================================================================

	// write manga metadata to yaml file
	mangaMetaData(rootPath, mangaName, mangaDetails.Description, mangaDetails.Genres, chapterLinks, masterImageArr)


	// =======================================================
	// create a folder for each chapter
	// =======================================================
	for _, cd := range chapterLinks {
		chapterName := cd.Name
		chapterPath, _ := filepath.Abs(filepath.Join(rootPath, chapterName))
		err := os.MkdirAll(chapterPath, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating chapter folder: %v\n", err)
			os.Exit(1)
		}
	}

	// =======================================================
	// download images
	// =======================================================
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
	// =======================================================
}