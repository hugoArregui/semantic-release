package release

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func runGit(args ...string) ([]string, error) {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Println(string(exitErr.Stderr))
		} else {
			fmt.Println(err.Error())
		}
		return nil, err
	}
	r := strings.Split(strings.TrimSpace(string(out)), "\n")
	return r, nil
}

func runGitOneLine(args ...string) (string, error) {
	out, err := runGit(args...)
	if err != nil {
		return "", err
	}

	return out[0], nil
}

func getCurrentBranch() (string, error) {
	return runGitOneLine("rev-parse", "--abbrev-ref", "HEAD")
}

func getLastCommit() (string, error) {
	return runGitOneLine("rev-parse", "HEAD")
}

func getCommitsBetween(f, t string) ([]string, error) {
	if f == "" {
		return nil, errors.New("invalid commit range, no from povided")
	}

	if t == "" {
		var err error
		t, err = getLastCommit()
		if err != nil {
			return nil, err
		}
	}

	return runGit("rev-list", fmt.Sprintf("%s..%s", f, t))
}

func getCommitTitle(c string) (string, error) {
	return runGitOneLine("show", "-s", "--pretty=%s", c)
}

func getCommitBody(c string) (string, error) {
	return runGitOneLine("show", "-s", "--pretty=%b", c)
}
