#!/bin/bash

docker-compose down --remove-orphans
docker volume rm $(docker volume ls -qf dangling=true)
rm -rf persistent_volumes/db/*
docker-compose up --build -d
