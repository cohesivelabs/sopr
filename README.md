# SOPR #

The purpose of sopr is to provide a number of convenience methods when interacting with multiple git repositories at one time.

It will allow you to expose scripts for each project/repo and provide an easy way to run those scripts against one or more of the projects at once.

Additionally it provides wrapper for a number of git operations to be performed as bulk operations, for example creating branches with the same name for one or more projects at once or performing mass git updates.

## Build from source

1. Install Go 1.11+ - [(instructions)](https://golang.org/doc/install)
2. Ensure that the ENV variable GOPATH is set and that `$GOPATH/bin` is in your path [(detailed here)](https://golang.org/doc/code.html)
5. Ensure that go module support is enabled
3. Clone this repo to your machine
4. From the root directory of the repository run `go install`

## Installation
Availabe on homebrew `brew install sopr`, for other operating systems see the releases and download the appropriate binary and place it on your cli path

## Usage

### Initialization
From the root of your project run `sopr init` this will clone all projects.

### Custom Scripts
Sopr supports adding and running custom scripts.

#### Project Level Scripts
Custom scripts that run under the context of one or more configured projects.

add the `--all` or `-a` flag to run the script for all projects that have that script configured for.
if the script is run without the `--all` flag, the user will be given a prompt to select which projects they want to run the script against

#### Top Level Scripts
Scripts that don't necessarily have an explicit tie to any configured projects but still need to run under the context of the root level directory. Any scripts names used here will override any script names defined at the project level


#### Script Descriptions
Since both project and top level scripts are added to the cli help files at run time, if you add a `descriptions` property to your `sopr.yaml` `sopr run --help` will provide additional help and context to the cli
_Please note though_ that the description key must match the name of the script exactly

```yaml
descriptions:
    test-script: 'this is a test script'

scripts:
    - name: test-script
      command: 'echo test'
```

### Git Operations
sopr supports a number of git operations that can be performed in bulk against one or more configured projects.

#### Current Support
* create a branch
* switch to a specific branch
* update repos
* list all configured repos and their current branch

#### Remotes
In order to streamline repository configuration, a user can optionally define any number of remotes that will be added at initialization.

If no remote named `origin` is provided, the first configured remoted will be used to clone the repo.

### Configuration
Sopr expects there to be a configuration file at the root level of your project named `sopr.yaml`.
* `projectDirectory` [string] (required) - base directory to clone project repositories into.
* `scripts` [array] (optional) - scripts that are not directly associated with projects, these will be executed from the context of the `projectDirectory`
* `scripts.name` [string] [required] - name of the script
* `scripts.command` [string] [required] - script to run
* `scripts.options.onInit.enable` [bool] [optional] - indicate that the script should be run after repos are cloned using `sopr init`
* `scripts.options.onInit.order` [number] [optional] - the order that the init scripts will run
* `projects` [array] (required) - list a project repositories.
* `projects.name` [string] (required) - display name of repository in cli.
* `projects.path` [string] (required) - path relative to the `repoDirectory` to clone the repository.
* `projects.remotes` [array] (required) - list of remotes for repository.
* `projects.remotes.name` [string] (required) - name of the remote
* `projects.remotes.url` [string] (required) - url of remote
* `projects.scripts.name` [string] [required] - name of the script
* `project.scripts.command` [string] [required] - script to run

#### Example Config

```yaml
projectDirectory: "projects"

descriptions:
  "deps": "install dependencies for project(s)"
  "clean": "clean and remove dependencies for project(s)"
  "test": "this is a test"

scripts:
  - name: "test"
    command: "echo 'i worked'"

projects:
  - name: "test-project"
    path: "test"
    remotes:
        - name: origin
          url: git@github.com:jmartin84/sopr.git
    scripts:
      - name: "deps"
        command: "npm install"
        options:
          onInit:
            enable: true
      - name: "clean"
        command: "rm -rf node_modules"
```

## Troubleshooting

### ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain

Sopr delegates authentication to a locally running ssh-agent if you see an error like the above, do the following:

1. If you're on windows install pagent and add your identity
2. If you're on osx or linux make sure that you have an identity set `ssh-add -l`, if not use `ssh-add` to add your identity (default identity is id_rsa)

If none of the above work then you're going to have to find someone who understands trouble shooting ssh auth better than I do.

## Special Thanks
Special thanks to [Surge Forward](https://www.surgeforward.com) whom the original iteration of `sopr` was written for.
