run: 
	go run src/main.go

run-single:
	go run src/main.go -link="https://mangaclash.com/manga/shadowless-night/chapter-4/" -single=true

run-range:
	go run src/main.go -link="https://mangaclash.com/manga/shadowless-night/" -start=0 -end=3

run-rdm:
	go run src/main.go -link="https://readm.org/manga/owari-no-seraph/"

run-rdm-ch:
	go run src/main.go -link="https://readm.org/manga/owari-no-seraph/125/all-pages" -single=true

build:
	go build -o bin/main src/main.go