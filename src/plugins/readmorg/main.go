package readmorg

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"downloader/src/helpers"
)

// implement IMangaDownloader
type ReadmOrg struct {}

func (m ReadmOrg) Download(pathRoot string, link string) {
	download(pathRoot, link)
}
func (m ReadmOrg) DownloadChapter(pathRoot string, link string) {
	downloadChapter(pathRoot, link)
}
func (m ReadmOrg) DownloadChapterRange(pathRoot string, link string, start int, end int) {
	downloadChapterRange(pathRoot, link, start, end)
}

// ===========================================================================================================
// DOWNLOADERS ===============================================================================================
// ===========================================================================================================

func download(pathRoot string, link string) {
	const FULL_MANGA = -1
	downloadChapterRange(pathRoot, link, 0, FULL_MANGA)
}

func downloadChapter(pathRoot string, link string) ChapterData {
	chapterDetails, err := getChapterDetails(link)
	if err != nil {
		if err.Error() == "404" {
			fmt.Println("Chapter not found")
		} else {
			fmt.Fprintf(os.Stderr, "error getting chapter details: %v\n", err)
			os.Exit(1)
		}
	}

	// ===========================================================================================================
	// IMAGE NAME CORRECTION =====================================================================================
	// if image names are not numbers, go in order of imageLinks and name them 1, 2, 3, ...s
	// else, use the image names
	imageNames := helpers.Map(chapterDetails.ImageLinks, func (image ImageData) string { return image.Name })
	areImageNamesNumbers := helpers.Map(imageNames, func (imageName string) bool { 
		// check if image name is a number
		return isNumber(strings.Split(imageName, ".")[0])
	 })
	anyNonNumbers := helpers.Any(areImageNamesNumbers, func (b bool) bool { return !b })
	
	if anyNonNumbers {
		for index, image := range chapterDetails.ImageLinks {
			name := strconv.Itoa(index + 1) + filepath.Ext(image.Name)
			chapterDetails.ImageLinks[index].Name = name
		}
	}
	// ===========================================================================================================

	// create folder for chapter in root/mangaName/chapterName
	pathRootManga := filepath.Join(pathRoot, chapterDetails.MangaName)
	pathRootMangaChapter := filepath.Join(pathRootManga, chapterDetails.Name)

	pathRootMangaChapter = strings.ReplaceAll(pathRootMangaChapter, ".", "_")
	err = os.MkdirAll(pathRootMangaChapter, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating manga folder: %v\n", err)
		os.Exit(1)
	}

	// get image links
	imageLinks := chapterDetails.ImageLinks

	// download images
	var wg sync.WaitGroup
	for _, image := range imageLinks {
		wg.Add(1)
		go func(_image ImageData) {
			downloadImage(_image.Link, _image.Name, pathRootMangaChapter)
			defer wg.Done()
		}(image)
	}
	wg.Wait()

	// write chapter details to yaml file
	writeChapterDetailsToFile(pathRootMangaChapter, chapterDetails)

	fmt.Println("Done downloading")
	return chapterDetails
}

func downloadChapterRange(pathRoot string, link string, start int, end int) {
	mangaDetails := getMangaData(link)
	chapterLinks := mangaDetails.ChapterLinks
	// reverse chapterLinks
	for i, j := 0, len(chapterLinks)-1; i < j; i, j = i+1, j-1 {
		chapterLinks[i], chapterLinks[j] = chapterLinks[j], chapterLinks[i]
	}

	if end == -1 {
		end = len(chapterLinks) - 1
	}

	var wg sync.WaitGroup
	for index, chapterLink := range chapterLinks {
		if index >= start && index <= end {
			wg.Add(1)
			go func(_chapterLink string) {
				defer wg.Done()
				downloadChapter(pathRoot, _chapterLink)
			}(chapterLink)
		}
	}

	wg.Wait()
}