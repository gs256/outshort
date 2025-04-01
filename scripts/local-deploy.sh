#!/bin/sh

set -ex
cd "$(dirname "$0")"/..

docker compose build

mkdir -p ./deploy

cp ./docker-compose.yml ./deploy
cp -r ./server/env ./deploy

docker compose up
