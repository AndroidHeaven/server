#!/bin/sh -e
#
# Usage: ./generate_apk.sh [ipa id] [app name]
#
# Output: Puts generated .apk into ./artifacts
#
# Warning: Only works with OS X

error () { echo $1; exit; }
run () { echo; echo \$ $*; $*; }

base_proj=android_base
ipa_id=$1
dest_dir=/tmp/apk/$ipa_id
start_dir=$(pwd)
app_name=$2

run mkdir -p $dest_dir
run cp -r $base_proj/ $dest_dir
run cd $dest_dir
run find . -type f -name "*.java" -exec sed -i '' -e "s/{{IPA_ID}}/$ipa_id/" {} \;
run find . -type f -name "AndroidManifest.xml" -exec sed -i '' -e "s/{{APP_NAME}}/$app_name/" {} \;
run ./gradlew assembleDebug
run mkdir -p $start_dir/artifacts && \
	cp app/build/outputs/apk/app-debug.apk $start_dir/artifacts/output.apk
run cd $start_dir
run rm -rf $dest_dir
