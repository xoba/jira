#!/bin/bash
export GO=`go run lib/files.go src/xoba | xargs`
echo README.md lib/*.sh *.sh *.java $GO


