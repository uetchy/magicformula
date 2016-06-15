# magicformula

Generate and upload Homebrew Formula like magic.

## Install

`go get` directly:

```
$ go get -d github.com/uetchy/magicformula
```

## Usage

with git repository:

```bash
magicformula build \
  --kind golang \
  --url "https://github.com/uetchy/awesomeapp.git" \
  --tag v1.2.0 > awesomeapp.rb
```

with archived package:

```bash
magicformula build \
  --path=./v1.2.0.tar.gz \
  --url="https://github.com/uetchy/awesomeapp/archive/v1.2.0.tar.gz" \
  --tag v1.2.0 > awesomeapp.rb
```

## Options

There are all of available options.

option        | description
------------- | ----------------------------------
--description | Package description [$DESCRIPTION]
--name        | Package name [$PACKAGE_NAME]
--tag         | Release tag [$RELEASE_TAG]
--homepage    | Homepage [$HOMEPAGE]
--path        | Package path [$PACKAGE_PATH]
--url         | Package url [$PACKAGE_URL]

## [Wercker](http://wercker.com/) step

See [wercker-step-homebrew](https://github.com/uetchy/wercker-step-homebrew) to get further information.

## Contributing

See [CONTRIBUTING.md](https://github.com/uetchy/magicformula/blob/master/CONTRIBUTING.md)
