package release

import (
	"context"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
)

type Release struct {
	SHA     string
	Version *semver.Version
}

type Releases []*Release

func (r Releases) Len() int {
	return len(r)
}

func (r Releases) Less(i, j int) bool {
	return r[j].Version.LessThan(r[i].Version)
}

func (r Releases) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func getLatestRelease(ctx context.Context, ghClient *github.Client, owner, repo, vrange string) (*Release, error) {
	allReleases := make(Releases, 0)
	opts := &github.ReferenceListOptions{"tags", github.ListOptions{PerPage: 100}}
	for {
		refs, resp, err := ghClient.Git.ListRefs(ctx, owner, repo, opts)
		if resp != nil && resp.StatusCode == 404 {
			return &Release{"", &semver.Version{}}, nil
		}
		if err != nil {
			return nil, err
		}
		for _, r := range refs {
			version, err := semver.NewVersion(strings.TrimPrefix(r.GetRef(), "refs/tags/"))
			if err != nil {
				continue
			}
			allReleases = append(allReleases, &Release{r.Object.GetSHA(), version})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	sort.Sort(allReleases)

	var lastRelease *Release
	for _, r := range allReleases {
		if r.Version.Prerelease() == "" {
			lastRelease = r
			break
		}
	}

	if vrange == "" {
		if lastRelease != nil {
			return lastRelease, nil
		}
		return &Release{"", &semver.Version{}}, nil
	}

	constraint, err := semver.NewConstraint(vrange)
	if err != nil {
		return nil, err
	}
	for _, r := range allReleases {
		if constraint.Check(r.Version) {
			return r, nil
		}
	}

	nver, err := semver.NewVersion(vrange)
	if err != nil {
		return nil, err
	}

	splitPre := strings.SplitN(vrange, "-", 2)
	if len(splitPre) == 1 {
		return &Release{lastRelease.SHA, nver}, nil
	}

	npver, err := nver.SetPrerelease(splitPre[1])
	if err != nil {
		return nil, err
	}
	return &Release{lastRelease.SHA, &npver}, nil
}
