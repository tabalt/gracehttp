
# gracehttp

This is a simple and graceful http server for golang.


Usage
-----

```
package main

import (
    "fmt"
    "net/http"

    "github.com/tabalt/gracehttp"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "hello world")
    })

    err := gracehttp.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println(err)
    }
}
```

Demo
-----

1. Install the demo application

        go get github.com/tabalt/gracehttp/gracehttpdemo

1. Start it in the first terminal

        gracehttpdemo

    This will output something like:

        2015/09/14 20:01:08 Serving localhost:8080 with pid 4388.

1. In a second terminal start a slow HTTP request

        curl 'http://localhost:8080/sleep/?duration=20s'

1. In a third terminal trigger a graceful server restart (using the pid from your output):

        kill -SIGUSR2 $pid

1. Trigger another shorter request that finishes before the earlier request:

        curl 'http://localhost:8080/sleep/?duration=0s'



