#!/bin/bash

echo "Preparing files"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main .
tar czf assets.tgz assets

Echo "Uploading files"
scp main diabasis:main2
scp assets.tgz root@10.0.10.19:assets.tgz

echo "Cleaning local files"
rm assets.tgz
rm main

echo "Launching in remote computer"
ssh root@10.0.10.19 tar -xzf assets.tgz 
ssh root@10.0.10.19 rm assets.tgz
ssh root@10.0.10.19 mv main2 main
ssh root@10.0.10.19 pkill main


