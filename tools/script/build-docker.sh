#!/usr/bin/env bash

set -eu
set -o pipefail

# Get current directory name
result=${PWD##*/}


# Set the docker registry to use
docker_registry="${DOCKER_REGISTRY:=localhost:5000}"

docker_cmd="${DOCKER_CMD:=docker}"

# Specify if the image will be push to the registry
push_to_registry="${PUSH_TO_REGISTRY:=0}"

# Tag to use
tag="${DEFAULT_TAG:=latest}"

docker_file="tools/docker/${TOOLS_SUB_FOLDER}/Dockerfile"

print_help () {
    echo "
        Usage: ${BASH_SOURCE[0]} [OPTION]

        Create the docker image and take care of the image deployment, if asked too.

        -h          Print help
        -n NAME     Set the docker image name to NAME, default to current directory name.
        -p          Publish the image to the docker registry.
        -r REGISTRY Set the url of the docker registry to use.
        -t TAG      Set the tag of the docker image to TAG.
        -f FILE     Path to the Dockerfile to use. default script/Dockerfile
    "
}

# Parse common parameter
argument_parse () {
    while getopts ":hn:pf:r:t:" opt; do
        case ${opt} in
            n)
                # Get the name to use
                result="$OPTARG"
                ;;
            p)
                push_to_registry=1
                ;;
            r)
                docker_registry="${OPTARG}"
                ;;
            t)
                tag="${OPTARG}"
                ;;
            f)
                docker_file="${OPTARG}"
                ;;
            h)
                print_help
                exit 0
                ;;
        esac
    done
}

# Get command line argument
argument_parse ${@}

# Create docker image
echo "[.] Creating docker image webtag/${result}:${tag}..."
docker build --rm=false -f ${docker_file} -t ${result}:${tag} --build-arg VERSION=`./tools/script/version.sh` .
if [[ $? -ne 0 ]]; then
    echo "[!] Docker image build got an error, exiting...'"
    exit 1
fi

docker tag ${result}:${tag} ${result}:`./tools/script/version.sh`
if [[ $? -ne 0 ]]; then
    echo "[!] Docker image tagging got an error, exiting...'"
    exit 1
fi

echo ${push_to_registry}
# Push the docker image to the registry
if [[ "${push_to_registry}" == 1 ]]; then
    echo "[.] Push image to the repository ${docker_registry} ..."
    docker tag ${result}:${tag} ${docker_registry}/${result}:${tag}
    if [[ $? -ne 0 ]]; then
        echo "[!] Docker tagging from '/${result}:${tag}' to '${docker_registry}/${result}:${tag}' got an error, exiting...'"
        exit 1
    fi

    docker tag ${result}:`./tools/script/version.sh` ${docker_registry}/${result}:`./tools/script/version.sh`
    if [[ $? -ne 0 ]]; then
        echo "[!] Docker tagging from '${result}:`./tools/script/version.sh`' to '${docker_registry}/${result}:`./tools/script/version.sh`' got an error, exiting...'"
        exit 1
    fi

    ${docker_cmd} push ${docker_registry}/${result}:${tag}
    if [[ $? -ne 0 ]]; then
        echo "[!] Docker image (${docker_registry}/${result}:${tag}) push got an error, exiting...'"
        exit 1
    fi

    ${docker_cmd} push ${docker_registry}/${result}:`./tools/script/version.sh`
    if [[ $? -ne 0 ]]; then
        echo "[!] Docker image (${docker_registry}/${result}:`./tools/script/version.sh`) push got an error, exiting...'"
        exit 1
    fi
fi

echo "[.] Done !"
