package main

import (
	"log"
	"os"
	"flag"

	"github.com/hugoArregui/semantic-release/pkg/release"
)

func main() {
	ghToken, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		log.Fatal("missing GITHUB_TOKEN")
	}

	config := release.Config{
		GHToken:      ghToken,
		DebugEnabled: true,
	}

	branch, err := release.GetCurrentBranch()
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&config.FromCommit, "fromCommit", "", "range start")
	flag.StringVar(&config.ToCommit, "toCommit", "", "range end. Default: latest commit in the branch")
	flag.StringVar(&config.Branch, "branch", branch, "branch. Default: current branch")
	flag.StringVar(&config.Owner, "owner", "", "owner of the repo")
	flag.StringVar(&config.Repo, "repo", "", "repo name")
	flag.BoolVar(&config.IsPR, "isPR", true, "are we building a PR?")
	flag.Parse()

	if config.FromCommit == "" {
		log.Fatal("missing fromCommit")
	}

	if config.Owner == "" {
		log.Fatal("missing repo owner")
	}

	if config.Repo == "" {
		log.Fatal("missing repo name")
	}

	if err := release.SemanticRelease(config); err != nil {
		log.Fatal(err)
	}
}
