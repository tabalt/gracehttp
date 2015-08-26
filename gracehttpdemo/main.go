package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tabalt/gracehttp"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello app on gracehttp!\n"))
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/sleep/", func(w http.ResponseWriter, r *http.Request) {
		duration, err := time.ParseDuration(r.FormValue("duration"))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		time.Sleep(duration)

		fmt.Fprintf(
			w,
			"started at %s slept for %d nanoseconds from pid %d.\n",
			time.Now(),
			duration.Nanoseconds(),
			os.Getpid(),
		)
	})

	err := gracehttp.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Println(err)
	}
}
