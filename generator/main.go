package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var (
	values    = []string{"red", "green", "blue"}
	votingURL string
	workers   int
	help      bool
)

func main() {

	flag.StringVar(&votingURL, "url", "http://localhost:8080/voting", "URL to post votes")
	flag.IntVar(&workers, "workers", 10, "Number of concurrent workers")
	flag.BoolVar(&help, "help", false, "Show this help")

	flag.Parse()

	if help {
		fmt.Printf("Usage: %s\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	ctx, cancel := context.WithCancel(context.Background())

	valuesStream := make(chan string)

	for i := 0; i < workers; i++ {
		go generateVotes(ctx, valuesStream, time.Second*time.Duration((i+1)))
	}

	go voter(ctx, valuesStream)

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)

	<-signals

	cancel()

}

func voter(ctx context.Context, values <-chan string) {

	for {
		select {
		case <-ctx.Done():
			return
		case v := <-values:
			fmt.Println("Voting", v)
			data := make(url.Values)
			data["vote"] = []string{v}
			resp, err := http.PostForm(votingURL, data)
			if nil != err {
				log.Println(err)
				continue
			}

			fmt.Printf("Voting result: %d [%s]\n", resp.StatusCode, resp.Status)
		}
	}
}

func generateVotes(ctx context.Context, stream chan<- string, interval time.Duration) {
	timer := time.NewTicker(interval)

	defer timer.Stop()

	for {
		idx := rand.Intn(len(values))
		stream <- values[idx]
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}
	}
}
