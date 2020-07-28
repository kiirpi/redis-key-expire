package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

func setExpire(item string, client redis.Client, ttl int, limitter chan struct{}) {
	client.Expire(item, (time.Duration(ttl))*time.Second)
	<-limitter
}

func main() {

	host := flag.String("host", "127.0.0.1", "host")
	port := flag.String("port", "6379", "port")
	pattern := flag.String("pattern", "", "pattern")
	db := flag.Int("db", 0, "db number")
	limit := flag.Int("limit", 2, "limit")
	ttl := flag.Int("ttl", 3600, "ttl")
	isDetail := flag.Bool("detail", false, "show some details")
	isShowHelp := flag.Bool("help", false, "show help menu")
	isForceTTL := flag.Bool("force", false, "force ttl ( if key has not ttl )")

	flag.Parse()

	if *isShowHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *isDetail {
		fmt.Printf("host:%s port:%s db:%d pattern:%s ttl:%d\n", *host, *port, *db, *pattern, *ttl)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     *host + ":" + *port,
		Password: "",
		DB:       *db,
	})

	list, err := client.Keys(*pattern).Result()

	if err != nil {
		flag.PrintDefaults()
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	defer client.Close()

	listLength := len(list)
	fmt.Printf("Total %d key found\n", listLength)
	now := time.Now()
	limiter := make(chan struct{}, *limit)
	remain := len(list)
	processed := 0
	for _, item := range list {
		remain = remain - 1
		keyTTL, _ := client.TTL(item).Result()
		if keyTTL == -1 || *isForceTTL {
			limiter <- struct{}{}
			go setExpire(item, *client, *ttl, limiter)
			processed++
		} else {
			if *isDetail {
				fmt.Printf("%s %f\n", item, time.Duration(keyTTL).Seconds())
			}
		}
	}
	diff := time.Now().Sub(now)
	fmt.Printf("Duration:%f Total:%d Processed:%d\n", diff.Seconds(), listLength, processed)
}
