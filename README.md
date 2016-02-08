# magicformula

CLI to generate Homebrew formula.

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
magicformula create
```

```session
magicformula create \
  --token "123456789abcdefghijklmnopqrstuvwxyz" \
  --name "awesome-cli-app" \
  --owner "uetchy" \
  --tag "v1.0.0" \
  --target-path-64 "./dist/darwin_amd64.tar.gz" \
  --target-path-32 "./dist/darwin_386.tar.gz" \
  --committer "uetchy" \
  --committer-email "uetchy@randompaper.co"
```

## Options

There are all of available options.

|option |description          |
|-------|---------------------|
|token  |Github access token  |
|name   |Formula repo         |
|owner  |Owner of formula repo|
|tag    |Release tag          |
|committer|Commit author      |
|committer-email|Commit author email|
|target-path-64|binary or package(64)|
|target-path-32 (optional)|binary or package(32)|
|version (optional)|Formula's Version|
|product-owner (optional)|Owner of product repo|
|message (optional)|Commit message|

## Working with [Wercker](http://wercker.com/)

See [wercker-step-homebrew](https://github.com/uetchy/wercker-step-homebrew) to get further information.

## Contributing

This step currently focusing on Golang project.
If you have any idea of creating formula for another language's project, please feel free to submit Pull-request or create issues.
