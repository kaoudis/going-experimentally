package main

import (
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

////////////////////////// end to end tests ////////////////////////////

// prove we can handle concurrent connections
func TestServer(t *testing.T) {
	testServer := httptest.NewServer(App())
	defer testServer.Close()

        // we'll use a client here instead of NewRequest I guess and throw multiple reqs
}
