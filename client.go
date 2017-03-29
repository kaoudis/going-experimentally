package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

// used as part of end-to-end test for server.

func generateWordsList() []string {
	words := make([]string, 0, 1)

	wordsfile, err := os.Open("/usr/share/dict/words")

	if err != nil {
		log.Fatal(err)
	}

	defer wordsfile.Close()

	scanner := bufio.NewScanner(wordsfile)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err1 := scanner.Err(); err1 != nil {
		log.Fatal(err1)
	}

	return words
}

func formatAsPassword(pass string) []byte {
	return []byte(fmt.Sprintf("password=%s", pass))
}

//following the pattern in https://blog.golang.org/pipelines
func generateRequests(passwords []string) <-chan *http.Request {
	outlet := make(chan *http.Request)

	go func() {
		for _, pass := range passwords {
			request, err := http.NewRequest("POST", "http://localhost:8080", bytes.NewBuffer(formatAsPassword(pass)))
			if err == nil {
				request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
				outlet <- request
			} else {
				log.Print("Oh no! Building a hash request for %s failed!\n", pass)
				continue
			}
		}

		//close channel 'outlet' after all password requests added to it
		close(outlet)
	}()

	return outlet
}
