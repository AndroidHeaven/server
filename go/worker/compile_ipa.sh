#!/bin/sh -e
#
# Usage: ./compile_ipa.sh [zip archive name]
#
# Output: Stores compiled .app in the ./artifacts directory.

echo "\n**********************"
echo "*    COMPILE IPA     *"
echo "**********************"

error () { echo $1; exit; }
run () { echo; echo \$ $*; $*; }

archive=$1
dest_dir=/tmp/$(md5 -q $archive)
start_dir=$(pwd)

if ! [[ -s ~/.mjolnir/init.lua ]] ; then
	cat mjolnir.lua > ~/.mjolnir/init.lua
fi ;

run mkdir $dest_dir
run cp $archive $dest_dir/
run cd $dest_dir
run unzip -q $archive
run cd $(find . -d 1 -type d | grep -v __)
run xcodebuild -arch i386 -sdk iphonesimulator
run mkdir -p $start_dir/artifacts && \
	cp -r build/Release-iphonesimulator/*.app $start_dir/artifacts/
run cd $start_dir && rm -rf $dest_dir
run cd artifacts/ && tar czf artifact.tar.gz $(find . -d 1 -type d)
