package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var (
	errMissingEnv      = "Missing env variable %s"
	addr               string
	redisAddr          = "${REDIS_HOST}:${REDIS_PORT}"
	redisPass          string
	clientIDcookieName = "clientID"
)

func main() {
	fmt.Println("Voting app started")

	tmp := "VOTING_ADDR"
	addr = os.Getenv(tmp)
	if "" == addr {
		panic(fmt.Errorf(errMissingEnv, tmp))
	}

	redisAddr = os.ExpandEnv(redisAddr)

	tmp = "REDIS_PASS"
	redisPass = os.Getenv(tmp)
	if "" == redisPass {
		panic(fmt.Errorf(errMissingEnv, tmp))
	}

	server := http.Server{
		Addr:              addr,
		Handler:           getHandler(),
		IdleTimeout:       time.Second * 30,
		ReadHeaderTimeout: time.Second * 5,
		ReadTimeout:       time.Second * 5,
		WriteTimeout:      time.Second * 15,
	}

	log.Println(server.ListenAndServe())
}

func getHandler() http.Handler {

	router := mux.NewRouter()

	router.HandleFunc("/voting", showVoting).Methods(http.MethodGet)
	router.HandleFunc("/voting", processVote).Methods(http.MethodPost)

	return router
}

//RespData data structure for html template
type RespData struct {
	Message string
}

func showVoting(w http.ResponseWriter, r *http.Request) {

	resp := &RespData{}

	tmpl := template.Must(template.ParseFiles("tmpl/voting.html"))
	if err := tmpl.Execute(w, resp); nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type RedisRecord struct {
	ClientID string `json:"clientID"`
	Vote     string `json:"vote"`
}

func processVote(w http.ResponseWriter, r *http.Request) {
	//TODO process POST

	vote := r.FormValue("vote")
	cookieClientID, err := r.Cookie(clientIDcookieName)
	if nil != err {
		log.Println(err)
	}

	clientID := generateSecureToken(32)

	if nil != cookieClientID {
		clientID = cookieClientID.Value
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:    clientIDcookieName,
			Value:   clientID,
			Expires: time.Now().Add(time.Hour * 24),
		})
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
	})

	defer client.Close()

	ctx := context.Background()

	record, err := json.Marshal(&RedisRecord{
		ClientID: clientID,
		Vote:     vote,
	})
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cmd := client.RPush(ctx, "votes", string(record), 0)

	msg := fmt.Sprintf("Thanks for voting for %s", vote)
	if nil != cmd.Err() {
		msg = fmt.Sprintf("ERROR: %s", err.Error())
	}

	resp := &RespData{
		Message: msg,
	}

	tmpl := template.Must(template.ParseFiles("tmpl/voting.html"))
	if err := tmpl.Execute(w, resp); nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
