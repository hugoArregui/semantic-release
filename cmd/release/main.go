package main

import (
	"github.com/hugoArregui/semantic-release/internal/release"
	"os"
	"log"
)

func main() {
	fromCommit := "0aff6e71f82ccc90697a005386f38ddc79d09cbc"

	ghToken, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		log.Fatal("missing GITHUB_TOKEN")
	}

	config := release.Config{
		FromCommit: fromCommit,
		GHToken: ghToken,
		Owner: "hugoArregui",
		Repo: "semantic-release",
	}

	if err := release.SemanticRelease(config); err != nil {
		log.Fatal(err)
	}
}
