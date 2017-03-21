package main

import (
        "fmt"
	"log"
	"net/http"
)

func main() {
    router := http.NewServeMux()

    router.HandleFunc("/", waitAndHash)

    err := http.ListenAndServe(":8080", router)
    log.Fatal(err)
}

func waitAndHash(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "hello friendly friends")
}
