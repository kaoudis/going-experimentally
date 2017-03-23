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
	router := http.NewServeMux()
	router.HandleFunc("/", hasher)

	waitingRouter := waiter(router)

	err := http.ListenAndServe(":8080", waitingRouter)
	log.Fatal(err)
}

func badRequestFormatter(request *http.Request) string {
	return fmt.Sprintf(
		"%s from %s \n\t(proto %s, host %s)\n",
		request.Method,
		request.URL,
		request.Proto,
		request.Host)
}

func hasher(writer http.ResponseWriter, request *http.Request) {
	sha512Hasherator := sha512.New()

	if request.Method == "POST" {
		request.ParseForm()
		toHash := request.Form.Get("password")

		if len(toHash) > 0 {
			sha512Hasherator.Write([]byte(toHash))
			hash := base64.StdEncoding.EncodeToString(sha512Hasherator.Sum(nil))
			fmt.Fprintf(writer, "%s", hash)
		} else {
			writer.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(writer, "Bad Request: 'password' field is required")
		}
	} else {
		log.Printf("Unacceptable request: %s\n", badRequestFormatter(request))
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
