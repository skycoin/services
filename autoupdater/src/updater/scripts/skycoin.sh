#!/bin/bash
# We hardcode the service and just care about the version
service_github_url="https://downloads.skycoin.net/wallet/skycoin-"
process_name="skycoin"
binary_directory=$GOBIN
release_compilation_tag="-bin-osx-darwin-x64"

version=$2
version=$(echo $version | sed s/v//)

echo "fetching" 
echo "${service_github_url}${version}${release_compilation_tag}.zip -O $process_name.zip"   
# fetch new version
cd $binary_directory
wget ${service_github_url}${version}${release_compilation_tag}.zip -O "${process_name}.zip"

echo "fetched"  
echo "${process_name}.zip" 

# unzip
unzip -o -a "${process_name}.zip" 
uncompressed_dir=${process_name}-${version}${release_compilation_tag}/${process_name}
cp "./$uncompressed_dir" ./$process_name

# Those are two ways of restarting the service:
############################################################################################################
# 1) Launching process from bash and disowning him, so it keeps running after the script exits
############################################################################################################
# kill running previous version
pid=pgrep -x $process_name 
kill $pid

echo "restarting" 
# start new version
nohup $process_name -gui-dir $uncompressed_dir/src/gui &

echo "service restarted"
############################################################################################################
# 2) Configure the service to run under systemctl
############################################################################################################
# systemctl restart $process_name
