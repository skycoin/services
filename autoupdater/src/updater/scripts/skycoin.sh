#!/bin/bash
# We hardcode the service and just care about the version
service_github_url="https://downloads.skycoin.net/wallet/skycoin-"
process_name="skycoin"
binary_directory=$GOBIN
release_compilation_tag="-bin-osx-darwin-x64"

version=$2
version=$(echo ${version} | sed s/v//)

echo "fetching" 
echo "${service_github_url}${version}${release_compilation_tag}.zip -O $process_name.zip"   
# fetch new version
cd ${binary_directory}
#wget ${service_github_url}${version}${release_compilation_tag}.zip -O "${process_name}.zip"

echo "fetched"  
echo "${process_name}.zip" 

# unzip
#unzip -o -a "${process_name}.zip"
uncompressed_binary=${process_name}-${version}${release_compilation_tag}/${process_name}

if [[ -z $(diff ${uncompressed_binary} ${process_name}) ]]; then
    echo "already up to date"
    exit 0
fi

cp "./${uncompressed_binary}" ./${process_name}

# Those are two ways of restarting the service:
############################################################################################################
# 1) Launching process from bash and disowning him, so it keeps running after the script exits
############################################################################################################
# kill running previous version
mkdir /tmp/${process_name}

pid=$(pgrep -x ${process_name})
kill ${pid}

echo "restarting" 
# start new version
nohup ${process_name} -gui-dir ${uncompressed_binary}/src/gui > /tmp/${process_name}/log.txt 2>&1 &

echo "service restarted"
############################################################################################################
# 2) Configure the service to run under systemctl
############################################################################################################
# systemctl restart $process_name
