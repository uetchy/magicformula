NAME=awesomeapp
OWNER=uetchy
RELEASE_TAG=v1.0.0
BRANCH=new-version-${RELEASE_TAG}
TITLE="New version: ${RELEASE_TAG}"
GITHUB_TOKEN=abcdefg12345
HUB_VERSION=2.2.1
HUB_PACKAGE=hub-linux-amd64-${HUB_VERSION}

git clone https://${GITHUB_TOKEN}@github.com/${OWNER}/homebrew-${NAME}
cd homebrew-${NAME}
git checkout -b $BRANCH
magicformula build > magicformula.rb
git commit -am "${TITLE}"
if [ $? -ne 0 ]; then
  echo "No changes"
  exit 0
fi

git push -u origin ${BRANCH}

# Download hub and set config up
curl -L https://github.com/github/hub/releases/download/v${HUB_VERSION}/${HUB_PACKAGE}.tar.gz | tar zxv
echo "---
github.com:
- protocol: https
  user: $OWNER
  oauth_token: $GITHUB_TOKEN" > $HOME/.config/hub

# Create pull-request
./${HUB_PACKAGE}/hub pull-request -m "${TITLE}" -b master -h ${BRANCH}
