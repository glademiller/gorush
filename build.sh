#!/bin/bash

IMAGE_FILE="netscale-technologies_gorush"

# Load a docker image from a .tar.gz file
gunzip < $IMAGE_FILE.tar.gz | docker load
