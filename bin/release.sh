#!/usr/bin/env bash
set -e

if (("$#" != 1)); then
  echo "Tag has to be provided"

  exit 1
fi

NOW=$(date +%s)
TOKEN=""
CURRENT_BRANCH="master"
VERSION=$1
BASEPATH=$(
  cd $(dirname $0)
  cd ../src/
  pwd
)

# Always prepend with "v"
if [[ $VERSION != v* ]]; then
  VERSION="v$VERSION"
fi

repos=$(ls $BASEPATH)
repos="$repos cli-skeleton web-skeleton api-skeleton grpc-skeleton mix"

for REMOTE in $repos; do
  echo ""
  echo ""
  echo "Cloning $REMOTE"
  TMP_DIR="/tmp/mix-split"
  REMOTE_URL="https://$TOKEN@github.com/mix-go/$REMOTE.git"

  rm -rf $TMP_DIR
  mkdir $TMP_DIR

  (
    cd $TMP_DIR

    git clone $REMOTE_URL .
    git checkout "$CURRENT_BRANCH"

    if [[ $(git log --pretty="%d" -n 1 | grep tag --count) -eq 0 ]]; then
      echo "Releasing $REMOTE"
      git tag $VERSION || true
      git push origin --tags
    fi
  )
done

TIME=$(echo "$(date +%s) - $NOW" | bc)
printf "Execution time: %f seconds\n" $TIME
