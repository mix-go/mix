#!/usr/bin/env bash
set -e

if (( "$#" != 1 ))
then
    echo "Tag has to be provided"

    exit 1
fi

NOW=$(date +%s)
CURRENT_BRANCH="master"
VERSION=$1
BASEPATH=$(cd `dirname $0`; cd ../src/; pwd)
ROOTPATH=$(cd `dirname $0`; cd ../../; pwd)

# Always prepend with "v"
if [[ $VERSION != v*  ]]
then
    VERSION="v$VERSION"
fi

repos=$(ls $BASEPATH)

for REMOTE in $repos
do
    echo ""
    echo ""
    echo "Cloning $REMOTE";
    TMP_DIR="/tmp/mix-split"
    REMOTE_URL="https://github.com/mix-go/$REMOTE.git"

    rm -rf $TMP_DIR;
    mkdir $TMP_DIR;

    (
        cd $TMP_DIR;

        git clone $REMOTE_URL .
        git checkout "$CURRENT_BRANCH";

        if [[ $(git log --pretty="%d" -n 1 | grep tag --count) -eq 0 ]]; then
            echo "github.com/mix-go/$REMOTE $VERSION"
            sed -i "" "s/github.com\/mix-go\/${REMOTE} v.*/github.com\/mix-go\/$REMOTE $VERSION/g" `find $ROOTPATH -name go.mod`
        else
            LASTAGID=`git rev-list --tags --max-count=1`
            LASTVERSION=`git describe --tags $LASTAGID`
            echo "github.com/mix-go/$REMOTE $LASTVERSION"
            sed -i "" "s/github.com\/mix-go\/${REMOTE} v.*/github.com\/mix-go\/$REMOTE $LASTVERSION/g" `find $ROOTPATH -name go.mod`
        fi
    )
done

# Update app ver

sed -i "" "s/Version = \".*/Version = \"${VERSION##*v}\";/g" ${BASEPATH}/console/application.go

TIME=$(echo "$(date +%s) - $NOW" | bc)

printf "Execution time: %f seconds\n" $TIME
