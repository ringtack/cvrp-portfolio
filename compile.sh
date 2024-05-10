#!/bin/bash

########################################
############# CSCI 2951-O ##############
########################################

# Get submodules
git submodule update --remote --merge
git submodule update --init

# Compile in each
cd cvrp && ./compile.sh && cd ../
cd cvrp-gen && ./compile.sh && cd ../

# Build wrapper
go build .
