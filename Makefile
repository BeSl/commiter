.PHONY: run
run:
	go run cmd/uploader/main.go

.PHONY: build
build:
	go build -o committer.exe cmd/uploader/main.go