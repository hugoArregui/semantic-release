package release

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func RunGit(args ...string) ([]string, error) {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		fmt.Println("command failed: git", strings.Join(args, " "))
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

func RunGitOneLine(args ...string) (string, error) {
	out, err := RunGit(args...)
	if err != nil {
		return "", err
	}

	return out[0], nil
}

func GetCurrentBranch() (string, error) {
	return RunGitOneLine("rev-parse", "--abbrev-ref", "HEAD")
}

func GetLastCommit() (string, error) {
	return RunGitOneLine("rev-parse", "HEAD")
}

func GetCommitsBetween(f, t string) ([]string, error) {
	if f == "" {
		return nil, errors.New("invalid commit range, no from provided")
	}

	if t == "" {
		var err error
		t, err = GetLastCommit()
		if err != nil {
			return nil, err
		}
	}

	return RunGit("rev-list", fmt.Sprintf("%s..%s", f, t))
}

func GetCommitTitle(c string) (string, error) {
	return RunGitOneLine("show", "-s", "--pretty=%s", c)
}

func GetCommitBody(c string) (string, error) {
	return RunGitOneLine("show", "-s", "--pretty=%b", c)
}
