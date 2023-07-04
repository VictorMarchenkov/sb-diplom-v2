#!/bin/sh
cd C:\Users\marvik\go\src\sb-diplom-v2
docker build -t diploma . -f Dockerfile
docker run --rm -p 8282:8282 -p 8383:8383 diploma