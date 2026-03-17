#!/bin/bash

platforms=("linux/amd64" "darwin/amd64" "windows/amd64" "linux/arm64" "windows/arm64" "darwin/arm64")

for platform in "${platforms[@]}"; do
  OS=$(echo $platform | cut -d'/' -f1)
  ARCH=$(echo $platform | cut -d'/' -f2)
  output_name="ndoujin-cli-${OS}-${ARCH}"
  if [ "$OS" = "darwin" ]; then
    output_name="ndoujin-cli-mac-${ARCH}"
  fi
  if [ "$OS" = "windows" ]; then
    output_name+=".exe"
  fi

  GOOS=$OS GOARCH=$ARCH go build -o $output_name main.go
  echo "Built $output_name"
done
