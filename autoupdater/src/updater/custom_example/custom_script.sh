#!/bin/bash
# We hardcode the service and just care about the version
service_github_url="https://github.com/me/my_service"
github_service_name="my_service"
process_name="my_service"
binary_directory=$GOBIN
release_compilation_tag="-darwin-amd64"

service_location="$GOPATH/src/$service_github_name"
version=$2

# fetch new version
cd $binary_directory
wget $service_github_name/releases/download/$version/${github_service_name}${release_compilation_tag} -O $process_name

# kill running previous version
pid=pgrep -x $process_name 
kill $pid

# start new version
$process_name &