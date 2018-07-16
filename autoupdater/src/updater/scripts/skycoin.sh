#!/bin/bash
# We hardcode the service and just care about the version
service_github_url="https://downloads.skycoin.net/wallet/skycoin-"
process_name="skycoin"
binary_directory=$GOBIN
release_compilation_tag="-bin-osx-darwin-x64.zip"

version=$2
version=$(echo $version | sed s/v//)

echo "fetching" > tmp.txt
echo "${service_github_url}${version}${release_compilation_tag} -O $process_name.zip" > tmp.txt  
# fetch new version
cd $binary_directory
wget ${service_github_url}${version}${release_compilation_tag} -O $process_name.zip
# unzip
unzip -a $process_name.zip

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