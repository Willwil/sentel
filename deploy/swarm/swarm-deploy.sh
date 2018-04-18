#!/usr/bin/env bash

nodes_nbr=4
env_already_exist=false

check_environment_already_exist() {
    echo "checking swarm environment..."
    # get all nodes
    nodes=($(docker-machine ls -q))
    # for each node, check wether the well-formed node exist
    count=0
    for node in ${nodes[@]}; do
        # save node-n into nodes_names
        if [[ $node =~ node-[0-$nodes_nbr] ]]; then
            ((count=count+1))
        fi
    done
    if [[ $nodes_nbr -eq $count ]]; then
        env_already_exist=true
        echo "iot swarm environment already exist with $count nodes"
    fi
}

ready_environment() {
    if [[ $env_already_exist = false ]]; then
        echo "no environment exist, making environment..."
        # create node
        for i in `seq $nodes_nbr`; do
            docker-machine create -d virtualbox node-$i
        done
        # setup docker swarm manager 
        eval $(docker-machine env node-1)
        ip=$(docker-machine ip node-1)
        token=$(docker swarm init --advertise-addr $ip | grep 'token' | head -n 1 | awk '{print $5}')
        echo "swarm token:$token"
        # add swarm worker nodes
        for i in `seq $nodes_nbr`; do
            docker-machine ssh node-$i "docker swarm join --token $token $ip:2377";
        done
        # list nodes status
        docker node ls
        # update docker image registry
        for i in `seq $nodes_nbr`; do   
           docker-machine ssh node-$i \
               "echo 'EXTRA_ARGS="--registry-mirror=https://registry.docker-cn.com"' | sudo tee -a /var/lib/boot2docker/profile"; 
        done
    else
        for i in `seq $nodes_nbr`;do
            node_status=$(docker-machine status node-$i)
            echo "node-$i status is $node_status"
            if [[ $node_status != "Running" ]]; then
                echo "starting node-$i..."
                docker-machine start node-$i
            fi
        done
    fi
}

swarm_start_service() {
    check_environment_already_exist
    ready_environment
    echo "swarm starting service..."
    docker stack deploy -c ./deploy/swarm/swarm-deploy.yaml sentel
    docker stack ps sentel 
    docker stack services sentel 
}

swarm_stop_service() {
    echo "swarm stop service..."
    docker stack rm sentel 
}

swarm_remove_service() {
    echo "swarm remove service..."
    for i in `seq $nodes_nbr`;do
            docker-machine rm node-$i
    done
}

if [ "$#" = "1" ];then
    case $1 in
        "start") swarm_start_service $1;;
        "stop") swarm_stop_service ;;
        "rm") swarm_remove_service $1;;
        *) 
            echo "invalid parameter '$1'"
            ;;
    esac
else
    echo "invalid parameter"
fi


