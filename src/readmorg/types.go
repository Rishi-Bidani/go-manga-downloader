package readmorg

type ChapterData struct {
	MangaName string
	Name string
	Link string
	ReleaseDate string
	ImageLinks []ImageData
	NumberOfImages int
}

type ImageData struct {
	Link string
	Name string
}

type MangaData struct {
	Name string
	Genres []string
	Description string
	ChapterLinks []string
}