#!/bin/bash

########################################
############# CSCI 2951-O ##############
########################################

# Get submodules
git pull --recurse-submodules

# Compile in each
cd cvrp && ./compile.sh && cd ../
cd cvrp-gen && ./compile.sh && cd ../

# Build wrapper
go build .
