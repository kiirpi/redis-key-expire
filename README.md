# Redis Key Expire

Search with pattern if ttl not exist then apply ttl

## Build

> GOOS=linux GOARCH=amd64 go build main.go

## Usage

Detail ( without keys )
> ./main --host=127.0.0.1 --port=6379 --db=1 --pattern=ABC* --ttl=3600

Detail ( with keys )
> ./main --host=127.0.0.1 --port=6379 --db=1 --pattern=ABC* --ttl=3600 --detail