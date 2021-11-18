package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	upstreamURL   string
	refreshPeriod int
	listen        string
	cacheContents atomic.Value
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func loadData() (string, error) {
	resp, err := http.Get(upstreamURL)
	if err != nil {
		// handle error
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return string(body), nil
}

func main() {
	var err error
	upstreamURL = getEnv("UPSTREAM_URL", "http://worldclockapi.com/api/json/utc/now")
	refreshPeriod, err = strconv.Atoi(getEnv("REFRESH_PERIOD", "3"))
	if err != nil {
		log.Fatal("REFRESH_PERIOD needs to be an integer.")
	}
	listen := ":8080"

	contents, _ := loadData()
	cacheContents.Store(contents)

	go func() {
		// Update `cacheContents` every `refreshPeriod` seconds
		for {
			time.Sleep(time.Duration(refreshPeriod) * time.Second) // Sleep first, we started initialized.
			contents, err := loadData()
			if err != nil {
				log.Print("Something went wrong updating the data, re-using prior cache.")
				log.Print(err)
			} else {
				log.Print("Refreshing data.")
				cacheContents.Store(contents)
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s\n", cacheContents.Load())
	})

	log.Fatal(http.ListenAndServe(listen, nil))

}
