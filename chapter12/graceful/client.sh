#!/bin/bash -xu

function get() {
    curl localhost:8080
}

function get5() {
	get &
	get &
	get &
	get &
	get &
}

get5
sleep 1
get5
sleep 1
get5
sleep 1
get5
sleep 1
get5
pkill -HUP start_server &
sleep 1

get5
sleep 1
get5
sleep 1
get5
sleep 1
get5
sleep 1
get5
sleep 1

# 全部終わるの待つ
sleep 5
pkill -HUP start_server

echo done
