#!/bin/bash
# We hardcode the service and just care about the version
echo "hey"

service_github_url="https://github.com/me/my_service"
github_service_name="my_service"
process_name="my_service"
binary_directory=$GOBIN
release_compilation_tag="-darwin-amd64"

service_location="$GOPATH/src/$service_github_name"
version=$2

echo "fetching" > tmp.txt
echo "$service_github_url/releases/download/$version/${github_service_name}${release_compilation_tag} -O $process_name" > tmp.txt  
# fetch new version
cd $binary_directory
wget $service_github_url/releases/download/$version/${github_service_name}${release_compilation_tag} -O $process_name
echo "fetched" > tmp.txt
# Those are two ways of restarting the service:
############################################################################################################
# 1) Launching process from bash and disowning him, so it keeps running after the script exits
############################################################################################################
# kill running previous version
pid=pgrep -x $process_name 
kill $pid

echo "restarting" >> tmp.txt
# start new version
nohup $process_name &

############################################################################################################
# 2) Configure the service to run under systemctl
############################################################################################################
systemctl restart $process_name