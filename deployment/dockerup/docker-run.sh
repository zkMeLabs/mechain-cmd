#!/bin/bash
docker run --rm --network mechain-network -v ./deployment/dockerup/:/root/.mechain-cmd zkmelabs/mechain-cmd $1
