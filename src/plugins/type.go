package modules

type IMangaDownloader interface {
	Download(pathRoot string, link string)
	DownloadChapter(pathRoot string, link string)
	DownloadChapterRange(pathRoot string, link string, start int, end int)
}