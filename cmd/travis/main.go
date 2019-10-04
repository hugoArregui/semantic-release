package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hugoArregui/semantic-release/pkg/release"
)

func main() {
	ghToken, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		log.Fatal("missing GITHUB_TOKEN")
	}

	branch, ok := os.LookupEnv("TRAVIS_BRANCH")
	if !ok {
		log.Fatal("missing TRAVIS_BRANCH")
	}

	commitRange, ok := os.LookupEnv("TRAVIS_COMMIT_RANGE")
	if !ok {
		log.Fatal("missing TRAVIS_COMMIT_RANGE")
	}

	slug, ok := os.LookupEnv("TRAVIS_REPO_SLUG")
	if !ok {
		log.Fatal("missing TRAVIS_REPO_SLUG")
	}

	pr, ok := os.LookupEnv("TRAVIS_PULL_REQUEST")
	if !ok {
		log.Fatal("missing TRAVIS_PULL_REQUEST")
	}

	commits := strings.Split(commitRange, "...")
	repo := strings.Split(slug, "/")

	config := release.Config{
		FromCommit:   commits[0],
		ToCommit:     commits[1],
		GHToken:      ghToken,
		Owner:        repo[0],
		Repo:         repo[1],
		Branch:       branch,
		IsPR:         pr != "0",
		DebugEnabled: true,
	}

	err := release.SemanticRelease(config)
	if err != nil {
		if err != release.ErrInvalidCommitRange {
			log.Fatal(err)
		}

		fmt.Println("invalid commit range, cannot execute semantic-release")
	}
}
