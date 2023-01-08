#!/bin/bash

buildProj(){
    echo "building stocks-data from docker"
    docker build -t metro-stock-data . 
    echo "building project ..."
    echo "project built"
}

buildProj
