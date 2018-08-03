#!/bin/bash
# We hardcode the service and just care about the version
process_name=$1
version=$2
shift 2
arguments=$@
parsed_args=$(eval echo ${arguments})

echo "process name: ${process_name}"

binary="${GOPATH}/src/github.com/skycoin/skywire/cmd/${process_name}"

socksc_binary="${GOPATH}/src/github.com/skycoin/skywire/cmd/socks/socksc"
sockss_binary="${GOPATH}/src/github.com/skycoin/skywire/cmd/socks/sockss"
sshc_binary="${GOPATH}/src/github.com/skycoin/skywire/cmd/ssh/sshc"
sshs_binary="${GOPATH}/src/github.com/skycoin/skywire/cmd/ssh/sshs"

service_github_url="github.com/skycoin/skywire"
binary_directory=${GOBIN}

build_and_copy_if_different () {
    cd $1
    go build

    if [[ -z $(diff ${GOBIN}/${2} ./${2}) ]]; then
        echo "already up to date"
        exit 0
    fi

    cp ${2} ${GOBIN}/${2}
}

build_and_copy () {
    cd $1
    go build

    cp ${2} ${GOBIN}/${2}
}

echo "fetching"
echo "go get -d -u ${service_github_url}"

# fetch new version
go get -d -u ${service_github_url}
exit_status=$?

if [ ${exit_status} != 0 -a ${exit_status} != 1 ]; then
    exit 1
fi

echo "fetched"

echo "updating..."

build_and_copy ${socksc_binary} "socksc"
build_and_copy ${sockss_binary} "sockss"
build_and_copy ${sshc_binary} "sshc"
build_and_copy ${sshs_binary} "sshs"

build_and_copy_if_different ${binary} ${process_name}

echo "updated"
echo "restarting..."

cd ${GOBIN}
pkill -9 -F ${process_name}.pid

if [ ${process_name} == "manager" ]; then
    echo ${parsed_args}
    #nohup ./manager ${arguments} > /dev/null 2>&1 &sleep 3
    nohup ./manager ${parsed_args} > ${process_name}.log 2>&1 &echo $! > ${process_name}.pid &sleep 3
else
    #nohup ./node ${arguments} > /dev/null 2>&1 &sleep 3
    nohup ./node ${parsed_args} > ${process_name}.log 2>&1 &echo $! > ${process_name}.pid &sleep 3
fi
