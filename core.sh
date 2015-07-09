#!/bin/bash

for i in {0,2}
do
  (echo "SC" | nc gophercon2015.coreos.com 4001) &
done
#echo "SC" | nc gophercon2015.coreos.com 4001
