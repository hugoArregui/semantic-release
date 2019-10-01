package main

import (
	"log"
	"os"

	"github.com/hugoArregui/semantic-release/internal/release"
)

func main() {
	ghToken, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		log.Fatal("missing GITHUB_TOKEN")
	}

	branch := "master"

	config := release.Config{
		FromCommit: "0a12dc1f848a83bc2962e3a78a2f8a29bf98a6c6",
		GHToken:    ghToken,
		Owner:      "hugoArregui",
		Repo:       "semantic-releae",
		Branch:     branch,
	}

	if err := release.SemanticRelease(config); err != nil {
		log.Fatal(err)
	}
}
