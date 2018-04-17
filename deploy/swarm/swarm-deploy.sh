#!/usr/bin/env bash

nodes_nbr=4
env_already_exist=false

check_environment_already_exist() {
    # nodes_names hold all nodes's name, such as node-1
    declare -a node_names
    # get all nodes
    nodes=$(docker-machine ls -q)
    # for each node, check wether the well-formed node exist
    nodes_count=0
    for ((i = 0; i < ${#nodes}; ++i))
    do
        # save node-n into nodes_names
        if [[ ${nodes[i]} =~ "node-[0-$nodes_nbr]" ]]; then
            ${node_names[$nodes_index]} = ${nodes[i]}
            ++$nodes_count
        fi
        if [[ $nodes_count -eq $nodes_nbr ]]; then
            $env_already_exit = true
        fi
    done
}

ready_environment() {
    if [[ $env_already_exist = true ]]; then
        echo "environment already exist"
    else 
        echo "readying environment..."
        for i in seq $nodes_nbr;
        do
            docker-machine create -d virtualbox node-$i;
        done
        # setup docker swarm manager (node-1's ip address is 192.168.99.100)
        eval $(docker-machine env node-1)
        token=$(docker swarm init --advertise-addr 192.168.99.100 | grep 'token' )
        #| awk '{print $$2}')
        echo "#######"
        echo $token
        # add swarm worker nodes
        for ((i = 1; i < $nodes_nbr;++i)); 
        do
            docker-machine ssh node-$i "docker swarm join --token $token 192.168.99.100:2377";
        done
        # list nodes status
        docker node ls
        # update docker image registry
        for i in `seq $nodes_nbr`; 
        do   
           docker-machine ssh node-$i \
               "echo 'EXTRA_ARGS="--registry-mirror=https://registry.docker-cn.com"' | sudo tee -a /var/lib/boot2docker/profile"; 
        done
        # restart woker node
        docker-machine restart $(docker-machine ls -q)
    fi
}

swarm_start_service() {
    check_environment_already_exist
    ready_environment
    echo "swarm starting service..."
    docker stack deploy -c swarm-deploy.yaml sentel
    docker stack ps sentel 
    docker stack service sentel 
}

swarm_stop_service() {
    echo "swarm stop service..."
}

swarm_remove_service() {
    echo "swarm remove service..."
    docker stack rm sentel 
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


