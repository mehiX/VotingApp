package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var values = []string{"red", "green", "blue"}
var votingURL = "http://localhost:8080/voting"

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	valuesStream := make(chan string)

	for i := 0; i < 10; i++ {
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
