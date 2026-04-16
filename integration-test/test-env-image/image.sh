#/bin/bash

LATEST_TAG="latest"
IMAGE_NAME="alonza0314/nsctl-test"

build_image() {
    if ! docker build -f Dockerfile -t $IMAGE_NAME:$LATEST_TAG .; then
        echo "Failed to build the Docker image."
        exit 1
    fi
}

push_image() {
    if ! docker push $IMAGE_NAME:$LATEST_TAG; then
        echo "Failed to push the Docker image."
        exit 1
    fi
}

main() {
    build_image
    push_image
}

main "$@"