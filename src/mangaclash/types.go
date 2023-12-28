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