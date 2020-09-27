#! /bin/bash

python ../tools/genxmlgo.py -p config -o ../config/appconfig ../xmlparsers/appconfig.xml
go build -o zt2-webservice

