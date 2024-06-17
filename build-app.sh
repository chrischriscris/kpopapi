#!/usr/bin/env bash

set -e

go build -o app cmd/main.go

npm install -D tailwindcss
npx tailwindcss init
npx tailwindcss -i ./public/static/styles.css -o ./public/static/output.css
