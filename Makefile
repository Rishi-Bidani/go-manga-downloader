clean:
	rm -rf test/*

run: clean
	go run src/main.go

run-single: clean
	go run src/main.go -link="https://mangaclash.com/manga/shadowless-night/chapter-4/" -single=true

run-range: clean
	go run src/main.go -link="https://mangaclash.com/manga/shadowless-night/" -start=0 -end=3

run-rdm: clean
	go run src/main.go -link="https://readm.org/manga/owari-no-seraph/"

run-rdm-ch: clean
	go run src/main.go -link="https://readm.org/manga/owari-no-seraph/125/all-pages" -single=true

run-rdm-range: clean
	go run src/main.go -link="https://readm.org/manga/owari-no-seraph/" -start=0 -end=2

build:
	go build -o bin/main src/main.go
	