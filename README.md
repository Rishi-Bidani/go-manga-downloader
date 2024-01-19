# Manga Downloader

## How to use the program

### 1. Download the latest release

Available from the [releases page](https://github.com/Rishi-Bidani/go-manga-downloader/releases/tag/v1.0)

### 2. Get usage information by running the following command

```bash
./main.exe -h
```

### 3. Quick start

#### Download full manga

It is not recommended to use this command for very large mangas, as it is likely to break. Use the range command instead and download in batches of 100.

```bash
./main.exe -link="https://mangaclash.com/manga/[mangaName]"
```

#### Download single chapter

```bash
./main.exe -link="https://mangaclash.com/manga/[mangaName]/[chapterName]" -single=true
```

#### Downloading a range of chapters

```bash
./main.exe -link="https://mangaclash.com/manga/[mangaName]" -start=0 -end=9
```

#### Supported sites

- [x] [Manga Clash](https://mangaclash.com/)
- [x] [Readm.org](https://readm.org/)

## Build it yourself

```bash
make build
```

or with go cli

```bash
go build -o main src/main.go
```

## How to read the manga

Use my other manga viewer project [here](https://github.com/Rishi-Bidani/mangaviewer)
