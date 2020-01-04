#!/bin/bash
## This script helps regenerating multiple client instances in different network namespaces using Docker
## This helps to overcome the ephemeral source port limitation
## Usage: ./setup <connections> <number of instances> <ramp> <server ip>
## Example ./client 50 3 10 172.17.0.2

CONNECTIONS=$1
INSTANCES=$2
RAMP=$3
IP=$4

for (( i=0; i<${INSTANCES}; i++ ))
do
    docker run -d -l sock-client sock-client /client -conn=${CONNECTIONS} -ramp=${RAMP} -ip=${IP}
done
