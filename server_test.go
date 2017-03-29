package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

//////////////////// unit tests /////////////////////////

func TestRequestLoggingFormatter(t *testing.T) {
	correctLoggable := "POST from localhost; proto HTTP/1.1\n"

	request, _ := http.NewRequest("POST", "localhost", nil)
	loggable := requestLoggingFormatter(request)

	if loggable != correctLoggable {
		t.Errorf("requestLoggingFormatter created: %s", loggable)
	}
}

// in the ideal case we should be getting back a 200 and a b64encoded hash as the body of our request
func TestHasherHappyPath(t *testing.T) {
	handler := http.HandlerFunc(hasher)
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "localhost", strings.NewReader("password=angryMonkey"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("hasher returned incorrect response (%v) for correctly formatted request: %s\n", responseRecorder.Code, request.Body)
	}

	expectedResponse := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="

	encodedResponseBody, _ := responseRecorder.Body.ReadString('\n')
	stringlyResponseBody := string(encodedResponseBody)

	if expectedResponse != stringlyResponseBody {
		t.Errorf("Response body was incorrect: %s\n", stringlyResponseBody)
	}
}

//behaviour for handling non-POSTs
func TestHasherNonPOST(t *testing.T) {
	handler := http.HandlerFunc(hasher)
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "localhost", nil)
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("Hasher function returned incorrect response code for a method we don't handle: %s\n", strconv.Itoa(responseRecorder.Code))
	}

}

//behaviour for handling requests without the 'password' field
func TestHasherBadParams(t *testing.T) {
	handler := http.HandlerFunc(hasher)
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "localhost", strings.NewReader("hello"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("Hasher function returned incorrect response code for when 'password' is not present: %s\n", strconv.Itoa(responseRecorder.Code))
	}
}

////////////////////////// end to end tests  ////////////////////////////

// Prove we can handle concurrent connections. We use ListenAndServe, which calls Serve,
// which (docs: https://golang.org/src/net/http/server.go?s=57297:57357#L2614) creates
// a new goroutine for each incoming connection. Therefore, we don't need to do
// anything other than tell the router what to do with those connections -- in server.go,
// we wrap the router in a waiting middleware after telling it to use the hasher function
// to handle each of these connections.
func TestServer(t *testing.T) {
	testServer := httptest.NewServer(App())
	defer testServer.Close()

	wordsList := generateWordsList()
	clientRequestsChannel := generateRequests(wordsList)

	client := &http.Client{}

	responses := make([]string, 0, 1)

	// send requests to server
	for req := range clientRequestsChannel {
		resp, err := client.Do(req)

		if err != nil {
			log.Fatal(err)
		}

		bodyData, err := ioutil.ReadAll(resp.Body)
		responses = append(responses, string(bodyData))
	}
}
