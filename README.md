# Redis Key Expire

Applies TTL for each matching key

## Build

> GOOS=linux GOARCH=amd64 go build main.go

## Usage

Show help menu
> ./main --help

Detail ( without keys )
> ./main --host=127.0.0.1 --port=6379 --db=1 --pattern=ABC* --ttl=3600 --limit=2

Detail ( with keys )
> ./main --host=127.0.0.1 --port=6379 --db=1 --pattern=ABC* --ttl=3600 --detail  --limit=2