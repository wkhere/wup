go:
	go vet

install: go
	go install

build-dev:
	go build -o wup-dev

restart-dev:
	pkill wup-dev || true
	./wup-dev &

dev: build-dev restart-dev


.PHONY: go install build-dev install-dev
