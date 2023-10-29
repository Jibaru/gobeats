#!/bin/bash

APP_NAME="gobeats"
VERSION="1.0.0"
OUTPUT_DIR="bin"
MAIN_FILE="./cmd/main/main.go"

PLATFORMS=("windows/386" "windows/amd64" "linux/386" "linux/amd64" "darwin/amd64" "linux/arm")

# Clean and create the output directory
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Iterate through platforms and compile the application
for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    OUTPUT_NAME="$OUTPUT_DIR/${APP_NAME}_${VERSION}_${GOOS}_${GOARCH}"

    # Compile for the current platform
    if [[ "$GOOS" == "windows" ]]; then
        env GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUTPUT_NAME".exe "$MAIN_FILE"
    else
        env GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUTPUT_NAME" "$MAIN_FILE"
    fi

    # Check if the compilation was successful
    if [ $? -eq 0 ]; then
        echo "Compilation successful for $PLATFORM"
    else
        echo "Error compiling for $PLATFORM"
    fi
done

echo "Compilation complete. Binaries are available in the $OUTPUT_DIR directory."
