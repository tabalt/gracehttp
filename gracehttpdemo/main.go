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
	http.HandleFunc("/sleep/", func(w http.ResponseWriter, r *http.Request) {
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

	pid := os.Getpid()
	address := ":8080"

	log.Printf("process with pid %d serving %s.\n", pid, address)
	err := gracehttp.ListenAndServe(address, nil)
	log.Printf("process with pid %d stoped, error: %s.\n", pid, err)
}
