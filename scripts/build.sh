#!/usr/bin/env bash

cd ../src;
go install UP2P;
cd client; go install UP2P/client;
cd ../server; go install UP2P/server;
cd ../node; go install UP2P/node;
cd ..; go build;
cd ../scripts;