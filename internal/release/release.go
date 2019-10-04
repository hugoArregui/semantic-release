package release

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var commitPattern = regexp.MustCompile("^(feat|fix|docs|style|refactor|perf|test|chore)(?:\\((.*)\\))?\\: (.*)$")
var breakingPattern = regexp.MustCompile("BREAKING CHANGES?")

type Config struct {
	FromCommit string
	ToCommit   string
	GHToken    string
	Owner      string
	Repo       string
	Branch     string
	IsPR       bool
}

func SemanticRelease(config Config) error {
	commits, err := getCommitsBetween(config.FromCommit, config.ToCommit)
	if err != nil {
		return err
	}

	fmt.Println("commits", commits)

	newReleaseType := "patch"
	for _, commit := range commits {
		title, err := getCommitTitle(commit)
		if err != nil {
			return err
		}

		body, err := getCommitBody(commit)
		if err != nil {
			return err
		}

		if len(title) == 0 {
			return fmt.Errorf("invalid empty commit message, commit: %s", commit)
		}

		if len(title) > 70 {
			return fmt.Errorf("commit title too long, commit: %s", commit)
		}

		found := commitPattern.FindAllStringSubmatch(title, -1)
		if len(found) < 1 {
			return fmt.Errorf(`commit title did not follow semantic versioning: %s.
Please see https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#commit-message-format`, title)
		}

		changeType := strings.ToLower(found[0][1])
		changeScope := found[0][2]
		changeMessage := found[0][3]

		if breakingPattern.MatchString(body) {
			changeType = "breaking-change"
		}

		switch changeType {
		case "breaking-change":
			newReleaseType = "major"
		case "feat":
			if newReleaseType == "patch" {
				newReleaseType = "minor"
			}
		}

		fmt.Printf("commit: %s, title: %s, body: %s\n", commit, title, body)
		fmt.Printf("type: %s, scope: %s, message:%s \n", changeType, changeScope, changeMessage)
	}

	fmt.Println("new release type is", newReleaseType)

	if config.Branch == "master" && !config.IsPR {
		ctx := context.TODO()
		oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.GHToken}))
		ghClient := github.NewClient(oauthClient)
		latestVersion, err := getLatestVersion(ctx, ghClient, config.Owner, config.Repo)
		if err != nil {
			return err
		}

		fmt.Println("latest version:", latestVersion.String())

		var newVersion semver.Version
		switch newReleaseType {
		case "major":
			newVersion = latestVersion.IncMajor()
		case "minor":
			newVersion = latestVersion.IncMinor()
		case "feat":
			newVersion = latestVersion.IncPatch()
		}

		fmt.Println("new version:", newVersion.String())

		tag := fmt.Sprintf("v%s", newVersion.String())
		ref := "refs/tags/" + tag
		tagOpts := &github.Reference{
			Ref:    &ref,
			Object: &github.GitObject{SHA: &commits[0]},
		}
		_, _, err = ghClient.Git.CreateRef(ctx, config.Owner, config.Repo, tagOpts)
		if err != nil {
			return err
		}
	}

	return nil
}

func getLatestVersion(ctx context.Context, ghClient *github.Client, owner, repo string) (*semver.Version, error) {
	opts := &github.ReferenceListOptions{"tags", github.ListOptions{PerPage: 100}}
	lastVersion := &semver.Version{}
	for {
		refs, resp, err := ghClient.Git.ListRefs(ctx, owner, repo, opts)
		if resp != nil && resp.StatusCode == 404 {
			break
		}
		if err != nil {
			return nil, err
		}
		for _, r := range refs {
			version, err := semver.NewVersion(strings.TrimPrefix(r.GetRef(), "refs/tags/"))

			if lastVersion.LessThan(version) {
				lastVersion = version
			}

			if err != nil {
				continue
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return lastVersion, nil
}
