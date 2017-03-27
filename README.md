# going-experimentally

A smol Go server that just returns base64-encoded hashes of client 'password' requests.

# build and run (get back a b64-encoded hash)
1. In window/tab 1: `go build && ./going-experimentally`
2. In window/tab 2: `curl -vX POST "http://localhost:8080" -d "password=frogs"`

The result in window 2 should be `3WZK2iUU8F4TjzXGlXDpO1fkjfNMPp5Pv+Mu9kAoVoPAylYukEfpXcGV5Cp5ddGsgAbaShIRYTDTXzp+QidhVw==`

Note: server will wait 5 seconds on each request with the connection open before answering.

# build and run tests
`go build && go test -v`

# see test coverage (opens a new browser window)
`go test -coverprofile=coverage.out && go tool cover -html=coverage.out` 
