@echo off
title Start Microservices

echo Starting services...
start cmd /k "redis-server.exe"
start cmd /k "etcd.exe"
start cmd /k "cd app/friend && go run main.go"
start cmd /k "cd app/message && go run main.go"
start cmd /k "cd app/user && go run main.go"
start cmd /k "cd app/websocket && go run main.go"
start cmd /k "cd rabbitmq && go run main.go"
start cmd /k "cd api && go run main.go"


echo All services started!