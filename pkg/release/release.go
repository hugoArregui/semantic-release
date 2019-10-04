package release

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type versionChangeType int

const (
	versionChangeTypePatch versionChangeType = iota + 1
	versionChangeTypeMinor
	versionChangeTypeMajor
)

var ErrInvalidCommitRange = errors.New("invalid commit range")

func (c versionChangeType) String() string {
	switch c {
	case versionChangeTypePatch:
		return "patch"
	case versionChangeTypeMinor:
		return "minor"
	case versionChangeTypeMajor:
		return "major"
	}
	return "Unknown version change type"
}

var commitPattern = regexp.MustCompile(`^(feat|fix|docs|style|refactor|perf|test|chore)(?:\((.*)\))?\: (.*)$`)
var breakingPattern = regexp.MustCompile("BREAKING CHANGES?")

type Config struct {
	FromCommit   string
	ToCommit     string
	GHToken      string
	Owner        string
	Repo         string
	Branch       string
	IsPR         bool
	DebugEnabled bool
}

type Logger struct {
	debugEnabled bool
}

func (l *Logger) Debug(msg string, v ...interface{}) {
	if l.debugEnabled {
		fmt.Printf(msg, v...)
	}
}

func (l *Logger) Info(msg string, v ...interface{}) {
	fmt.Printf(msg, v...)
}

func SemanticRelease(config Config) error {
	logger := Logger{debugEnabled: config.DebugEnabled}
	commits, err := GetCommitsBetween(config.FromCommit, config.ToCommit)
	if err != nil {
		return ErrInvalidCommitRange
	}

	logger.Debug("commits in range: %s\n", strings.Join(commits, ", "))

	newReleaseType := versionChangeTypePatch
	for _, commit := range commits {
		title, err := GetCommitTitle(commit)
		if err != nil {
			return err
		}

		body, err := GetCommitBody(commit)
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

		if breakingPattern.MatchString(body) {
			newReleaseType = versionChangeTypeMajor
		} else if changeType == "feat" && newReleaseType < versionChangeTypeMinor {
			newReleaseType = versionChangeTypeMinor
		}
		logger.Debug("type: %s, scope: %s, message: %s\n", changeType, found[0][2], found[0][3])
	}

	logger.Debug("change is %s\n", newReleaseType.String())

	if config.Branch == "master" && !config.IsPR {
		ctx := context.TODO()
		oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.GHToken}))
		ghClient := github.NewClient(oauthClient)
		latestVersion, err := GetLatestVersion(ctx, ghClient, config.Owner, config.Repo)
		if err != nil {
			return err
		}

		var newVersion semver.Version
		switch newReleaseType {
		case versionChangeTypeMajor:
			newVersion = latestVersion.IncMajor()
		case versionChangeTypeMinor:
			newVersion = latestVersion.IncMinor()
		case versionChangeTypePatch:
			newVersion = latestVersion.IncPatch()
		default:
			panic("invalid release type")
		}

		logger.Info("current version: %s, change: %s, new version: %s\n", latestVersion.String(),
			newReleaseType.String(),
			newVersion.String())

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
		logger.Debug("pushed new tag %s\n", tag)
	}

	return nil
}

func GetLatestVersion(ctx context.Context, ghClient *github.Client, owner, repo string) (*semver.Version, error) {
	opts := &github.ReferenceListOptions{Type: "tags", ListOptions: github.ListOptions{PerPage: 100}}
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
