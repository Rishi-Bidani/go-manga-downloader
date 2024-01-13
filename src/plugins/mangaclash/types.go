package mangaclash

type ChapterData struct {
	Name string
	Link string
	ReleaseDate string
}

type ChapterImage struct {
	ChapterName string
	ImageLink string
	ImageName string
}

type MangaData struct {
	Name string
	Genres []string
	Description string
}

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