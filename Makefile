run: 
	go run src/main.go

run-single:
	go run src/main.go -link="https://mangaclash.com/manga/shadowless-night/chapter-4/" -single=true

run-range:
	go run src/main.go -link="https://mangaclash.com/manga/shadowless-night/" -start=0 -end=3

build:
	go build -o bin/main src/main.go