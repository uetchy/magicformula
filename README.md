# magicformula
Generate and upload Homebrew Formula.

## Install
with Homebrew:

```
$ brew install uetchy/magicformula/magicformula
```

or `go get` directly:

```
$ go get -d github.com/uetchy/magicformula
```

## Usage

```session
RELEASE_TAG=v1.0.0
GITHUB_USER=uetchy
GITHUB_TOKEN=123456789abcdefghijklmnopqrstuvwxyz
MF_PACKAGE_PATH="./bin/darwin_amd64"
magicformula push
```

## Options
There are all of available options.

option                            | description
--------------------------------- | ---------------------
GITHUB_TOKEN                      | Github access token
GITHUB_USER                       | Owner of formula repo
RELEASE_TAG                       | Release tag
MF_PACKAGE_PATH                   | binary or package
MF_PACKAGE_NAME (optional)        | Package name
MF_GIT_COMMITTER (optional)       | Commit author
MF_GIT_COMMITTER_EMAIL (optional) | Commit author email
MF_COMMIT_MESSAGE (optional)      | Commit message

## Working on [Wercker](http://wercker.com/)
See [wercker-step-homebrew](https://github.com/uetchy/wercker-step-homebrew) to get further information.

## Contributing
See [CONTRIBUTING.md](https://github.com/uetchy/magicformula/blob/master/CONTRIBUTING.md)
