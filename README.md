# solver

Keep your Homebrew's formula fresh.

## Install

```session
$ go get -d github.com/uetchy/solver
```

## Usage

```session
solver push \
  --token "123456789abcdefghijklmnopqrstuvwxyz" \
  --name "awesome-cli-app" \
  --owner "uetchy" \
  --version "1.0.0" \
  --target-path-64 "./dist/darwin_amd64.tar.gz" \
  --target-path-32 "./dist/darwin_386.tar.gz" \
  --committer "uetchy" \
  --committer-email "uetchy@randompaper.co"
```

## Working with [Wercker](http://wercker.com/)

Using solver with Wercker is awesome combination!
