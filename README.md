# Manga Downloader

## How to use the program

1. Download the latest release from the [releases page](https://github.com/Rishi-Bidani/go-manga-downloader/releases/tag/v1.0)

2. Get usage information by running the following command

```bash
./main.exe -h
```

3. Quick start

Download full manga

```bash
./main.exe -link="https://mangaclash.com/manga/[mangaName]"
```

Download single chapter

```bash
./main.exe -link="https://mangaclash.com/manga/[mangaName]/[chapterName]" -single=true
```

Downloading a range of chapters

```bash
./main.exe -link="https://mangaclash.com/manga/[mangaName]" -start=0 -end=9
```

## Build it yourself

```bash
make build
```

or with go cli

```bash
go build -o main src/main.go
```
