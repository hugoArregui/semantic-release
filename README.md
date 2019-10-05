# semantic-release

I was very frustrated with the current solutions for semantic release in go, both for linting a commit message to match the angular commit message format and for actually releasing a new version. Most of them require node, the ones in go assume too much about the underline CI tool and in general do more than what I needed. That's why I decided to create this very simple tool to solve my personal needs, I'm aiming at keeping it very short and simple. I'm happy to review problems you may have and I'm open to PRs but also feel free to just copy paste the code here into your own solution or just import this one as a library. I feel most of the basic use cases are very simple and can be solved with a custom script instead of a huge version release library, specially if you can program it in go instead of hacking bash or installing a second programming language.

- [Angular commit message format reference](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#commit-message-format)

## Usage

### As a library

```go
import (
	"github.com/hugoArregui/semantic-release/pkg/release"
)

func main() {
    config := release.Config{
        FromCommit   string // commit range start
        ToCommit     string // commit range end
        GHToken      string // the GH token
        Owner        string // owner of the repo
        Repo         string // name of the repo
        Branch       string // the branch to build
        IsPR         bool   // is this a PR build?
        DebugEnabled bool   // print debug info
    }
    release.SemanticRelease(config)
}
```

### In travis

in `cmd/travis/main.go` there is a special command line interface for travis, since it's what I'm going to be using. Remember to set the env var `GITHUB_TOKEN` with your GH token.

#### Caveats

Sometimes TRAVIS_COMMIT_RANGE is invalid, for example, if you add a new commit (sha: "A") to master, push it, and then `--append` (resulting in sha: "B") and push it again, TRAVIS_COMMIT_RANGE will be `A..B`, but A is not in the repo anymore. This doesn't happen for PRs though, so when the commit range is invalid I decided to simply abort the operation.

### Generic command line

in `cmd/custom/main.go` there is a generic command line in which you can build the config by using command line flags.

Remember to set the env var `GITHUB_TOKEN` with your GH token.
