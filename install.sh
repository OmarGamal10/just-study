#!/bin/bash


set -e

APP_NAME="study"
INSTALL_DIR="/usr/local/bin"
SOURCE_FILE="main.go"

if [ "$EUID" -ne 0 ]; then
  echo "Error: This script must be run as root."
  echo "Usage: sudo ./install.sh"
  exit 1
fi

if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH."
    exit 1
fi

go build -ldflags "-s -w" -o $APP_NAME $SOURCE_FILE # -s -w for a smaller binary

mv $APP_NAME $INSTALL_DIR/$APP_NAME

chmod +x $INSTALL_DIR/$APP_NAME

echo "  Usage:"
echo "    sudo $APP_NAME on  # To block sites"
echo "    sudo $APP_NAME off  # To unblock sites"
echo "    sudo $APP_NAME status  # To check the current status"