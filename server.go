package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	err := http.ListenAndServe(":8080", App())

	if err != nil {
		log.Fatal(err)
	}
}

// note: could also log request body, but I'm not sure if that's actually helpful.
func requestLoggingFormatter(request *http.Request) string {
	return fmt.Sprintf("%s from %s; proto %s\n", request.Method, request.URL, request.Proto)
}

// hash anything passed as the value to the 'password' parameter
func hasher(writer http.ResponseWriter, request *http.Request) {
	sha512Hasherator := sha512.New()

	if request.Method == "POST" {
		toHash := request.FormValue("password")

		if len(toHash) > 0 {
			sha512Hasherator.Write([]byte(toHash))
			hash := base64.StdEncoding.EncodeToString(sha512Hasherator.Sum(nil))
			fmt.Fprintf(writer, "%s", hash)
		} else { //400
			log.Printf("Bad request: %s", requestLoggingFormatter(request))
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(""))
		}
	} else { //405
		log.Printf("Method not allowed: %s", requestLoggingFormatter(request))
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte(""))
	}
}

// 'middleware': wait 5 seconds with connection open before doing anything
func waitingHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			time.Sleep(time.Duration(5) * time.Second)
			handler.ServeHTTP(writer, request)
		})
}

// split Handler out into a function callable in main so we can test our work
func App() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/", hasher)
	return waitingHandler(router)
}
