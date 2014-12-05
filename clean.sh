#!/bin/bash
# cleans the environment
source goinit.sh
go clean -i
find . -name "*~" -exec rm \{} \; 
rm -rf bin
rm -rf pkg
find . -name "*.class" -exec rm \{} \;
find . -name "*.test" -exec rm \{} \;
find . -name "*.java.orig" -exec rm \{} \;
find . -name "*flymake_*.go" -exec rm \{} \;
find . -name "y.output" -exec rm \{} \;
rm -rf build
rm -f *.jar

