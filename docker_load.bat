@echo off
title load MicroImages

echo Loading images...
start cmd /k "docker load < mq-server.tar"
start cmd /k "docker load < user-service.tar"
start cmd /k "docker load < websocket-service.tar"
start cmd /k "docker load < message-service.tar"
start cmd /k "docker load < friend-service.tar"
start cmd /k "docker load < my_api.tar"


echo All images loaded!