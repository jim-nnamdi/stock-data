#!/bin/bash

func buildProj(){
    echo "building stocks-data"
    go run *.go
    echo "building project ..."
    echo "project built"
}

buildProj()
n = buildProj()
echo $(n)
