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

	branch, ok := os.LookupEnv("TRAVIS_BRANCH")
	if !ok {
		log.Fatal("missing TRAVIS_BRANCH")
	}

	commit, ok := os.LookupEnv("TRAVIS_COMMIT")
	if ok {
		log.Println("TRAVIS COMMIT", commit)
	}

	commitRange, ok := os.LookupEnv("TRAVIS_COMMIT_RANGE")
	if ok {
		log.Println("TRAVIS COMMIT RANGE", commitRange)
	}

	// pr, ok := os.LookupEnv("TRAVIS_PULL_REQUEST")

	config := release.Config{
		FromCommit: fromCommit,
		GHToken: ghToken,
		Owner: "hugoArregui",
		Repo: "semantic-release",
Branch: branch,
	}

	if err := release.SemanticRelease(config); err != nil {
		log.Fatal(err)
	}
}
