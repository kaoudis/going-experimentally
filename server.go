package main

import (
    "crypto/sha512"
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
    sha512Hasherator := sha512.New()

    if (request.Method == "POST") {
        request.ParseForm()
        toHash := request.Form.Get("password")

        if (len(toHash) > 0) {
            sha512Hasherator.Write([]byte(toHash))
            fmt.Fprintf(writer, "You sent me %s\n", sha512Hasherator.Sum(nil))
        } else {
            writer.WriteHeader(http.StatusBadRequest)
            fmt.Fprintln(writer, "Bad Request: 'password' field is required")
        }
    } else {
        log.Printf("Unacceptable request received: %s\n", request)
        writer.WriteHeader(http.StatusNotAcceptable)
        writer.Write([]byte(""))
    }
}

func waiter(handler http.Handler) http.Handler {
    return http.HandlerFunc(
        func(writer http.ResponseWriter, request *http.Request) {
            time.Sleep(time.Duration(5) * time.Second)
            handler.ServeHTTP(writer, request)
        })

    }

