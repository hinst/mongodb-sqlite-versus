docker build --platform linux/arm64 -f build-arm.Dockerfile -t mongodb-sqlite-versus . &&^
    docker create mongodb-sqlite-versus --name mongodb-sqlite-versus