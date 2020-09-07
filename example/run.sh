#!/bin/bash
SCRIPT_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Paths.
CONFIG_PATH="$SCRIPT_PATH"

# Subpath.
SUB_PATH=$1
if [ -z "${SUB_PATH}" ]; then
  echo "Subpath required."
  exit 1
fi

# Which Slideshow.
SLIDE_PATH="${SCRIPT_PATH}/${SUB_PATH}"
if [ ! -d "${SLIDE_PATH}" ]; then
  echo "Path does not exist: $SLIDE_PATH"
  exit 1
fi

# Run the app with test configs.
$GOBIN/conditioning -config $CONFIG_PATH/config.json \
                    -affirm $SLIDE_PATH/affirmations.txt
[ $? -ne 0 ] && exit 1

# Everything is fine.
exit 0
