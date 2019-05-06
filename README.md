# SOPR #

> This tool is for organizing mutliple git repositories into a monorepo like experience and provides a number of helpers to streamline that experience

## Build from source

1. Install Go 1.11+ - [(instructions)](https://golang.org/doc/install)
2. Ensure that the ENV variable GOPATH is set and that `$GOPATH/bin` is in your path [(detailed here)](https://golang.org/doc/code.html)
5. Ensure that go module support is enabled
3. Clone this repo to your machine
4. From the root directory of the repository run `go install`

## Usage

### Configuration
Sopr expects there to be a configuration file at the root level of your project named `sopr.yaml`.
* `repoDirectory` [string] (required) - base directory to clone project repositories into.
* `repos` [array] (required) - list a project repositories.
* `repos.name` [string] (required) - display name of repository in cli.
* `repos.path` [string] (required) - path relative to the `repoDirectory` to clone the repository.
* `repos.remotes` [array] (required) - list of remotes for repository. If none are named origin then the first one will be used to clone the repository
* `repo.remotes.name` [string] (required) - name of the remote
* `repo.remotes.url` [string] (required) - url of remote
* `repo.installDeps` [string] (optional) - command that can be run to initialize any dependencies in the project repository, will be run relative to the repository
* `repo.removeDeps` [string] (optional) - command that can be run to remove any dependencies in the project repository, will be run relative to the repository

### Initialization
From the root of your project run `sopr init`

## Troubleshooting

### ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain

Sopr delegates authentication to a locally running ssh-agent if you see an error like the above, do the following:

1. If you're on windows install pagent and add your identity
2. If you're on osx or linux make sure that you have an identity set `ssh-add -l`, if not use `ssh-add` to add your identity (default identity is id_rsa)

If none of the above work then you're going to have to find someone who understands trouble shooting ssh auth better than I do.

## Special Thanks
Special thanks to [Surge Forward](https://www.surgeforward.com) whom the original iteration of `sopr` was written for.
