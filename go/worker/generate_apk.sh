#!/bin/sh -e
#
# Usage: ./generate_apk.sh [base Android project] [ipa id] [app name]
#
# Output: Puts generated .apk into ./artifacts
#
# Warning: Only works with OS X

echo "\n**********************"
echo "*    GENERATE APK    *"
echo "**********************"

error () { echo $1; exit; }
run () { echo; echo \$ $*; $*; }

base_proj=$1
ipa_id=$2
dest_dir=/tmp/apk/$ipa_id
start_dir=$(pwd)
app_name=$3

run mkdir -p $dest_dir
run cp -r $base_proj/ $dest_dir
run cd $dest_dir
find . -type f -name "*.java" -exec sed -i '' -e "s/{{IPA_ID}}/$ipa_id/" {} \;
find . -type f -name "AndroidManifest.xml" -exec sed -i '' -e "s/{{APP_NAME}}/$app_name/" {} \;
run ./gradlew assembleDebug
run mkdir -p $start_dir/artifacts && \
	cp app/build/outputs/apk/app-debug.apk $start_dir/artifacts/output.apk
