#!/bin/bash

set -e          # Exit immediately if a command exits with a non-zero status.
set -o pipefail # Check the exit code of a pipe and return the error code of the failing command

# APP_NAME env var is set in the Dockerfile
APP_NAME=${APP_NAME:-""}
if [ -z "$APP_NAME" ]; then
  echo "APP_NAME env var is not set. Exiting..."
  exit 1
fi

# Run the app
echo "Starting $APP_NAME..."
exec "$@"
