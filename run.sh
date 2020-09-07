#!/bin/bash
SCRIPT_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Paths.
CONFIG_PATH="$SCRIPT_PATH/conditioning/cmd/conditioning/config"

# We need to be in the write path.
cd $SCRIPT_PATH/conditioning
[ $? -ne 0 ] && exit 1

# Test the app.
go test -p=1 -count=1 ./...
[ $? -ne 0 ] && exit 1

# Update the app.
go install ./...
[ $? -ne 0 ] && exit 1

# Format the code the app.
go fmt ./...
[ $? -ne 0 ] && exit 1

# Run the app with test configs.
$GOBIN/conditioning -config $CONFIG_PATH/config.example.json \
                    -affirm $CONFIG_PATH/affirmations.example.txt
[ $? -ne 0 ] && exit 1

# Everything is fine.
exit 0
