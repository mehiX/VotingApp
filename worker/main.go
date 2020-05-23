package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

var (
	redisAddr             string
	redisPass             string
	mysqlConnectionString = "${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DATABASE}"
	errMissingEnv         = "Missing env variable %s"
)

func main() {
	fmt.Println("Worker started")

	redisAddr = os.Getenv("REDIS_ADDR")
	if "" == redisAddr {
		panic(fmt.Errorf(errMissingEnv, "REDIS_ADDR"))
	}

	redisPass = os.Getenv("REDIS_PASS")
	if "" == redisPass {
		panic(fmt.Errorf(errMissingEnv, "REDIS_PASS"))
	}

	mysqlConnectionString = os.ExpandEnv(mysqlConnectionString)

	ctx := context.Background()

	redis := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
	})

	defer redis.Close()

	mySQL, err := sql.Open("mysql", mysqlConnectionString)
	if nil != err {
		panic(err)
	}

	defer mySQL.Close()

	timeout := time.Minute * 30
	for {
		cmd := redis.BLPop(ctx, timeout, "votes")

		result, err := cmd.Result()
		if nil != err {
			log.Println(err)
			continue
		}

		record := result[1]
		if "0" != record {
			decoded := make(map[string]string)
			if err := json.Unmarshal([]byte(record), &decoded); nil != err {
				log.Println(err)
				continue
			}

			clientID := decoded["clientID"]
			vote := decoded["vote"]

			if err := upsertRecord(mySQL, clientID, vote); nil != err {
				log.Println(err)
			}
		}

	}
}

func upsertRecord(conn *sql.DB, clientID string, vote string) error {

	if err := conn.Ping(); nil != err {
		return err
	}

	insertSQL := "insert into votes (voterID, vote) values (?, ?)"
	updateSQL := "update votes set vote=? where voterID=?"
	selectSQL := "select vote from votes where voterID=?"

	row := conn.QueryRow(selectSQL, clientID)
	if err := row.Scan(); nil != err {
		if err == sql.ErrNoRows {
			_, err := conn.Exec(insertSQL, clientID, vote)
			return err
		}
	}

	_, err := conn.Exec(updateSQL, vote, clientID)
	return err

}
