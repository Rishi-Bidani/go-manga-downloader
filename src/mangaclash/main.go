package mangaclash

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

var LOG_SET bool = false

func mangaMetaData (filePath string, manga MangaData, chapterArr []ChapterData, imageArr []ChapterImage){
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
	mangaMeta.Name = manga.Name
	mangaMeta.Description = manga.Description
	mangaMeta.Genres = manga.Genres
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
	f, err := os.OpenFile(filepath.Join(filePath, manga.Name + ".yaml"), os.O_CREATE|os.O_WRONLY, 0664)
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

func Download(link string, pathRoot string){
	/*
		Download function will download the manga from the link provided
		It will download all the chapters and images
	*/
	chapterLinks, manga := getChapterLinksAndMangaDetails(link)
	
	// create a folder for the manga
	pathRootManga, _ := filepath.Abs(filepath.Join(pathRoot, manga.Name))
	err := os.MkdirAll(pathRootManga, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating manga folder: %v\n", err)
		os.Exit(1)
	}

	{
		// =================================================================================
		// setup logger
		// =================================================================================
		var f *os.File
		var err error
		if LOG_SET == false {
			f, err = os.OpenFile(filepath.Join(pathRootManga, "log.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error opening log file: %v\n", err)
				os.Exit(1)
			}
			log.SetOutput(f)
			defer f.Close()
			LOG_SET = true
		}
		// = end setup logger =============================================================
	}
	masterImageArr := []ChapterImage{}

	// download all chapters
	var wg sync.WaitGroup
	for _, chapter := range chapterLinks {
		wg.Add(1)
		go func(_chapter ChapterData) {
			downloadChapter(pathRoot, _chapter.Link, _chapter, manga, &masterImageArr)
			defer wg.Done()
		}(chapter)
	}
	wg.Wait()
	
	// write manga metadata to yaml file
	mangaMetaData(pathRootManga, manga, chapterLinks, masterImageArr)
}

func DownloadChapterRange(link string, pathRoot string, start int, end int){
	/*
		DownloadChapterRange function will download the manga from the link provided
		It will download all the chapters and images
	*/
	chapterLinks, manga := getChapterLinksAndMangaDetails(link)
	// reverse the chapterLinks array
	for i := len(chapterLinks)/2-1; i >= 0; i-- {
		opp := len(chapterLinks)-1-i
		chapterLinks[i], chapterLinks[opp] = chapterLinks[opp], chapterLinks[i]
	}
	
	
	// create a folder for the manga
	pathRootManga, _ := filepath.Abs(filepath.Join(pathRoot, manga.Name))
	err := os.MkdirAll(pathRootManga, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating manga folder: %v\n", err)
		os.Exit(1)
	}

	{
		// =================================================================================
		// setup logger
		// =================================================================================
		var f *os.File
		var err error
		if LOG_SET == false {
			f, err = os.OpenFile(filepath.Join(pathRootManga, "log.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error opening log file: %v\n", err)
				os.Exit(1)
			}
			log.SetOutput(f)
			defer f.Close()
			LOG_SET = true
		}
		// = end setup logger =============================================================
	}
	masterImageArr := []ChapterImage{}
	fmt.Println("Downloading chapters", start, "to", end, "of", len(chapterLinks))
	// download all chapters
	var wg sync.WaitGroup
	for i, chapter := range chapterLinks {
		if i >= start && i <= end {
			wg.Add(1)
			go func(_chapter ChapterData) {
				downloadChapter(pathRoot, _chapter.Link, _chapter, manga, &masterImageArr)
				defer wg.Done()
			}(chapter)
		}
	}
	wg.Wait()
}

func DownloadChapter(rootPath string, chapterLink string){
	chapter, mangaDetails := getSingleChapterLinkAndMangaDetails(chapterLink)
	downloadChapter(rootPath, chapterLink, chapter, mangaDetails, nil)
}

func downloadChapter(rootPath string, chapterLink string, chapter ChapterData, mangaDetails MangaData, masterImageArr *[]ChapterImage) {
	// if chapter data is not provided, get it
	if (chapter == (ChapterData{})) && (mangaDetails.Name == "") {
		chapter, mangaDetails = getSingleChapterLinkAndMangaDetails(chapterLink)
	}
	
	pathRootManga := filepath.Join(rootPath, mangaDetails.Name)
	
	err := os.MkdirAll(pathRootManga, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating manga folder: %v\n", err)
		os.Exit(1)
	}

	{
		// =================================================================================
		// setup logger
		// =================================================================================
		var f *os.File
		var err error
		if LOG_SET == false {
			f, err = os.OpenFile(filepath.Join(pathRootManga, "log.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error opening log file: %v\n", err)
				os.Exit(1)
			}
			log.SetOutput(f)
			defer f.Close()
			LOG_SET = true
		}
		// = end setup logger =============================================================
	}
	
	chapterImgArr := getImageLinks(chapter.Name, chapter.Link)

	// check if master image array is provided

	if masterImageArr != nil {
		*masterImageArr = append(*masterImageArr, chapterImgArr...)
	}

	pathRootMangaChapter := filepath.Join(pathRootManga, chapter.Name)
	err = os.MkdirAll(pathRootMangaChapter, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating chapter folder: %v\n", err)
		os.Exit(1)
	}

	// download all images
	var wg sync.WaitGroup
	for _, ci := range chapterImgArr {
		wg.Add(1)
		go func(_ci ChapterImage) {
			downloadImage(pathRootManga, _ci)
			defer wg.Done()
		}(ci)
	}
	wg.Wait()
}
