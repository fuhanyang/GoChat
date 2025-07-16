@echo off
title Build Microservices

echo Building services...
start cmd /k "docker build -f rabbitmq/Dockerfile -t mq-server --build-arg ENV=online ."
start cmd /k "docker build -f app/user/Dockerfile -t user-service --build-arg ENV=online ."
start cmd /k "docker build -f app/websocket/Dockerfile -t websocket-service --build-arg ENV=online ."
start cmd /k "docker build -f app/message/Dockerfile -t message-service --build-arg ENV=online ."
start cmd /k "docker build -f app/friend/Dockerfile -t friend-service --build-arg ENV=online ."
start cmd /k "docker build -f api/Dockerfile -t my_api --build-arg ENV=online ."

echo All services built!
