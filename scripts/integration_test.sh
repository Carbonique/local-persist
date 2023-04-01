#!/bin/bash

set -e

TAG=$1
PLUGIN=ghcr.io/carbonique/local-persist:${TAG}

function create-volume {
  VOLUME_ONE=`docker volume create --driver=${PLUGIN} --opt mountpoint=/docker-data/local-persist-integration/ --name=test-data-one`
  VOLUME_TWO=`docker volume create --driver=${PLUGIN} --name=test-data-two`
}

function create-containers {
    ONE=`docker run -d -v test-data-one:/app/data/ alpine sleep 30`
    TWO=`docker run -d -v test-data-one:/src/data/ alpine sleep 30`
    THREE=`docker run -d -v test-data-two:/app/data/ alpine sleep 30`
}

function check-containers {
    (docker exec $ONE cat /app/data/test.txt | grep 'Cameron Spear') || exit 111
    (docker exec $TWO cat /src/data/test.txt | grep 'Cameron Spear') || exit 222
    (docker exec $THREE cat /app/data/test.txt | grep 'Cameron Spear') || exit 111 
}

function clean {
    docker rm -f $ONE
    docker rm -f $TWO
    docker rm -f $THREE
    docker volume rm $VOLUME_ONE
    docker volume rm $VOLUME_TWO 
}
# setup
create-volume
create-containers

# copy a test file (note how this subtly breaks integration tests if my name is removed from the LICENSE ;-))
docker cp LICENSE $ONE:/app/data/test.txt
docker cp LICENSE $THREE:/app/data/test.txt 

# check that the file exists in all 
check-containers

# delete everything (start over point)
clean

# do it all again, but this time, DON'T manually copy a file... it should have persisted from before!
create-volume
create-containers

# if we were just using the `local` driver, this step would fail
check-containers

clean


echo -e "\nSuccess!"
