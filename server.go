package main

import (
        "fmt"
	"log"
	"net/http"
        "time"
)

func main() {
    router := http.NewServeMux()

    router.HandleFunc("/", hasher)

    waitingRouter := waiter(router)

    err := http.ListenAndServe(":8080", waitingRouter)
    log.Fatal(err)
}

func hasher(writer http.ResponseWriter, request *http.Request) {
    fmt.Fprintln(writer, "!!!!!!!!!!!!!!!!!1!!!!1!! HASHES COMING SOON !!!!!!!!!!!!!!!!1!!!!1!!!!!")
}

func waiter(handler http.Handler) http.Handler {
    return http.HandlerFunc(
        func(writer http.ResponseWriter, request *http.Request) {
            time.Sleep(time.Duration(5) * time.Second)
            handler.ServeHTTP(writer, request)
        })

}

