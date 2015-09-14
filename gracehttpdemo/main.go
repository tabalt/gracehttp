package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tabalt/gracehttp"
)

func main() {

	mux := http.NewServeMux()
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

	log.Println(fmt.Sprintf("Serving localhost:8080 with pid %d.", os.Getpid()))

	err := gracehttp.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Println(err)
	}
}
