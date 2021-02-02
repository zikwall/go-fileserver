#!/bin/bash

docker run -d -p 1337:1337 \
-v $HOME/tmp:/app/tmp \
-e BIND_ADDRESS='0.0.0.0:1337' \
-e TOKEN='123456' \
--name go-fileserver-example qwx1337/go-fileserver:latest