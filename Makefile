export GO111MODULE=on

.PHONY: build

build:
	go build -o bin/prometheus-rough-proxy .

crossbuild:
	gox -osarch='linux/amd64' -output='bin/{{.Dir}}_{{.OS}}_{{.Arch}}'

download:
	go mod download
