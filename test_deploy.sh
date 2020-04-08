#!/usr/bin/sh

shDir=/dvols/shree

generate_keys(){

    ssh-keygen -t rsa -b 4096 -P "" -f $shDir/keys/ca/ca_host_key
    ssh-keygen -t rsa -b 4096 -P "" -f $shDir/keys/ca/ca_user_key
    ssh-keygen -t rsa -b 4096 -P "" -f $shDir/keys/serv/id_host


}

mkdirs(){
    mkdir -p $shDir/ca $shDir/backend $shDir/serv
    mkdir -p $shDir/redis
    mkdir -p $shDir/keys/ca $shDir/keys/serv
}

deploy(){
    docker-compose start shree-stack
}
