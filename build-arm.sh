#!/bin/bash

if ! docker build --platform linux/arm64 -f build-arm.Dockerfile -t mongodb-sqlite-versus .
then
    echo "Build failed"
    exit 1
fi
if ! ctr_id=$(docker create mongodb-sqlite-versus --name mongodb-sqlite-versus)
then
    echo "Create failed"
    exit 1
fi
if ! docker cp "$ctr_id:/app/mongodb-sqlite-versus" .
then
    echo "Copy failed"
    exit 1
fi