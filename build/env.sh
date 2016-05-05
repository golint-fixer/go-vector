#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
vecdir="$workspace/src/github.com/vector"
if [ ! -L "$vecdir/go-vector" ]; then
    mkdir -p "$vecdir"
    cd "$vecdir"
    ln -s ../../../../../. go-vector
    cd "$root"
fi

# Set up the environment to use the workspace.
# Also add Godeps workspace so we build using canned dependencies.
GOPATH="$vecdir/go-vector/Godeps/_workspace:$workspace"
GOBIN="$PWD/build/bin"
export GOPATH GOBIN

# Run the command inside the workspace.
cd "$vecdir/go-vector"
PWD="$vecdir/go-vector"

# Launch the arguments with the configured environment.
exec "$@"
