#!/bin/bash -x

# creates the necessary docker images to run testrunner.sh locally

docker build --tag="vector/cppjit-testrunner" docker-cppjit
docker build --tag="vector/python-testrunner" docker-python
docker build --tag="vector/go-testrunner" docker-go
