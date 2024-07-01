#!/usr/bin/env bash

docker compose up -d
./scrape_images
docker compose stop
