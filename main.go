package main

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"sync"
	"time"
)

func setExpire(item string, client redis.Client, ttl int, wg *sync.WaitGroup, processed *int) {
	client.Expire(item, (time.Duration(ttl))*time.Second)
	*processed = *processed + 1
	wg.Done()
}

func log(msg string, isDetail bool) {
	if isDetail {
		fmt.Println(msg)
	}
}

func main() {

	host := flag.String("host", "127.0.0.1", "host")
	port := flag.String("port", "6379", "port")
	pattern := flag.String("pattern", "", "pattern")
	db := flag.Int("db", 5, "db number")
	ttl := flag.Int("ttl", 3600, "ttl")
	isDetail := flag.Bool("detail", false, "detail")

	flag.Parse()

	if *isDetail {
		fmt.Println("asdsada")
	}

	if *isDetail {
		fmt.Printf("host:%s port:%s db:%d pattern:%s ttl:%d\n", *host, *port, *db, *pattern, *ttl)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     *host + ":" + *port,
		Password: "",  // no password set
		DB:       *db, // use default DB
	})

	list, err := client.Keys(*pattern).Result()

	if err != nil {
		flag.PrintDefaults()
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	listLength := len(list)
	fmt.Printf("total %d key found\n", listLength)
	now := time.Now()
	var wg sync.WaitGroup
	remain := len(list)
	processed := 0
	for _, item := range list {
		remain = remain - 1
		keyTtl, _ := client.TTL(item).Result()
		if keyTtl == -1 {
			wg.Add(1)
			go setExpire(item, *client, *ttl, &wg, &processed)
			if *isDetail {
				fmt.Printf("not expired\n")
			}
		} else {
			if *isDetail {
				fmt.Printf("%s %f\n", item, time.Duration(keyTtl).Seconds())
			}
		}
	}
	wg.Wait()

	diff := time.Now().Sub(now)

	fmt.Printf("duration:%f total:%d processed:%d\n", diff.Seconds(), listLength, processed)

}
