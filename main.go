// Gemini Free API Bridge — gemini.google.com → OpenAI/Anthropic compatible proxy.
//
// Thin entry point: the whole bridge lives in internal/gbridge (see README
// "Project Structure"), mirroring the Qwen/GLM-Free-API layout this project
// follows. Cookie login is a separate helper under cmd/gemini-login.
//
// Build:  go build -trimpath -ldflags="-s -w" -o gemini-api .
// Run:    ./gemini-api          (or: go run .)

package main

import "github.com/Godde3s/gemini-free-api/internal/gbridge"

func main() {
	gbridge.Run()
}
