# going-experimentally

A smol Go server that just returns base64-encoded hashes of client 'password' requests.

# build and run (get back a b64-encoded hash)
1. In window/tab 1: `go build && ./going-experimentally`
2. In window/tab 2: `curl -vX POST "http://localhost:8080" -d "password=frogs"`

The result in window 2 should be `3WZK2iUU8F4TjzXGlXDpO1fkjfNMPp5Pv+Mu9kAoVoPAylYukEfpXcGV5Cp5ddGsgAbaShIRYTDTXzp+QidhVw==`

Notes: server will wait 5 seconds on each request with the connection open before answering. If you choose to test the non-happy paths, you will see some logging output from the server beginning with 'Method not allowed' or 'Bad request'.

# build and run tests
`go build && go test`

# see test coverage (opens a new browser window)
`go test -coverprofile=coverage.out && go tool cover -html=coverage.out` 
