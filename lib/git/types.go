package git

import (
	git "gopkg.in/src-d/go-git.v4"
    "sopr/lib/config"
)

type Repo struct {
    FullPath string
    Ref *git.Repository
    Config config.RepoConfig
}
