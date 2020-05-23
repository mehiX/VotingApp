package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var (
	resultsAddr           string
	mysqlConnectionString = "${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DATABASE}"
	queryVotes            = "select vote, count(*) as total from votes group by vote order by total desc"
	db                    *sql.DB
)

func main() {
	fmt.Println("Results app started")

	resultsAddr = os.Getenv("RESULTS_ADDR")
	mysqlConnectionString = os.ExpandEnv(mysqlConnectionString)

	if "" == resultsAddr {
		panic(errors.New("Missing env variable RESULTS_ADDR"))
	}

	var err error

	db, err = sql.Open("mysql", mysqlConnectionString)
	if nil != err {
		panic(err)
	}
	defer db.Close()

	server := http.Server{
		Addr:              resultsAddr,
		Handler:           getHandler(),
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 10,
		IdleTimeout:       time.Second * 30,
	}

	go func() {
		if err := server.ListenAndServe(); nil != err {
			log.Println(err)
		}
	}()

	signStream := make(chan os.Signal)
	signal.Notify(signStream, os.Interrupt)

	<-signStream
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	go func() {
		defer cancel()
		log.Println(server.Shutdown(ctx))
	}()

	<-ctx.Done()
	fmt.Println("Results app done")
}

func getHandler() http.Handler {

	router := mux.NewRouter()

	router.HandleFunc("/", showResults).Methods(http.MethodGet)
	router.HandleFunc("/json", showJSON).Methods(http.MethodGet)

	return router
}

func showResults(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("tmpl/results.html"))
	if err := tmpl.Execute(w, nil); nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func showJSON(w http.ResponseWriter, r *http.Request) {

	resp, err := getJSON()
	if nil != err {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// VotingResult record to marshal in the response
type VotingResult struct {
	Vote  string `json:"vote"`
	Total int    `json:"total"`
}

func getJSON() ([]*VotingResult, error) {
	votes := make([]*VotingResult, 0)

	rows, err := db.Query(queryVotes)
	if nil != err {
		return nil, err
	}

	for rows.Next() {
		var vote VotingResult
		if err := rows.Scan(&vote.Vote, &vote.Total); nil != err {
			return nil, err
		}
		votes = append(votes, &vote)
	}

	return votes, nil
}
