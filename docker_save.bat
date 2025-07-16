@echo off
title save MicroImages

echo Saving images...
start cmd /k "docker save mq-server > mq-server.tar"
start cmd /k "docker save user-service > user-service.tar"
start cmd /k "docker save websocket-service > websocket-service.tar"
start cmd /k "docker save message-service > message-service.tar"
start cmd /k "docker save friend-service > friend-service.tar"
start cmd /k "docker save my_api > my_api.tar"


echo All images saved!
