#!/bin/bash

VERSION=$(git describe --tags --always --dirty --match=v* 2> /dev/null || \
            cat $(CURDIR)/.version 2> /dev/null || echo 0.1.0-0)
RELEASE=$(date '+%Y%m%d.%s')
pushd cmd
    mv version.go version.bak.${RELEASE}
    sed -e "s/dirty/${RELEASE}/g" -e "s/dev/${VERSION}/g" version.bak.${RELEASE} > version.go
popd
