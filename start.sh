#!/bin/bash

# turn on bash's job control
set -m

# Start the primary process and put it in the background
echo "Start generator"
# Start the generator process
./generator/generator &

sleep 0.5

echo

echo "Start app"
./bind/app/diploma

fg %1

