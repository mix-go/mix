#!/usr/bin/env bash
set -e

if (("$#" != 1)); then
  echo "Tag has to be provided"

  exit 1
fi

NOW=$(date +%s)
CURRENT_BRANCH="master"
VERSION=$1
EXAMPLESPATH=$(cd $(dirname $0); cd ../examples/; pwd)
SRCPATH=$(cd $(dirname $0); cd ../src/; pwd)

# Always prepend with "v"
if [[ $VERSION != v* ]]; then
  VERSION="v$VERSION"
fi

repos=$(ls $EXAMPLESPATH)

for SKELETON in $repos; do
  echo $SKELETON
  (
    cd $(dirname $0)
    cd ../examples/$SKELETON
    go get -u all && go mod tidy
  )
done

# Update devtool version
sed -i "" "s/SkeletonVersion = \".*/SkeletonVersion = \"${VERSION##*v}\"/g" ${SRCPATH}/mixcli/commands/version.go
sed -i "" "s/CLIVersion      = \".*/CLIVersion      = \"${VERSION##*v}\"/g" ${SRCPATH}/mixcli/commands/version.go

TIME=$(echo "$(date +%s) - $NOW" | bc)
printf "Execution time: %f seconds\n" $TIME
