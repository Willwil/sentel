#!/usr/bin/env bash

usage() {
    cat << USAGE >&2
Usage:
    $cmdname k8s|dockercompose|dockerswarm args
    args:
    start start the deployment 
    stop  stop the deployment
    rm    remove the deployment
USAGE
    exit 1
}

dockerswarm() {
    echo "starting services in docker-swarm mode..."
    case $1 in
        "start")./deploy/swarm/swarm-deploy.sh $1;;
        "stop")./deploy/swarm/swarm-deploy.sh $1;;
        "rm")./deploy/swarm/swarm-deploy.sh $1;;
        *)
            echo "invalid action '$1'"
            usage;;
    esac

}

kubernetes() {
    echo "starting services in kubernetes mode..."
    case $1 in
        "start")
            kubectl apply -f ./deploy/kubernetes/zk_deploy.yaml
            kubectl apply -f ./deploy/kubernetes/kafka_deploy.yaml
            kubectl apply -f ./deploy/kubernetes/mongo_deploy.yaml
            kubectl apply -f ./deploy/kubernetes/iothub_deploy.yaml
            kubectl apply -f ./deploy/kubernetes/iothubmanager_deploy.yaml
            kubectl apply -f ./deploy/kubernetes/mns_deploy.yaml
            kubectl apply -f ./deploy/kubernetes/whaler_deploy.yaml
            kubectl apply -f ./deploy/kubernetes/apiserver_deploy.yaml
            ;;
        "stop")docker-compose down
            echo "kubernetes deploymode doesn't support stop action"
            ;;
        "rm")
            kubectl delete deploy whaler-deployment 
            kubectl delete deploy mns-deployment 
            kubectl delete deploy iotmanager-deployment 
            kubectl delete deploy iothub-deployment 
            kubectl delete deploy kafka-deployment 
            kubectl delete deploy zk-deployment 
            kubectl delete deploy mongo-deployment 
            ;;
        *)
            echo "invalid action '$1'"
            usage;;
    esac

}

dockercompose() {
    echo "starting services in docker-compose mode..."
    case $1 in
        "start")docker-compose up;;
        "stop")docker-compose down;;
        "rm")docker-compose -f -a;;
        *)
            echo "invalid action '$1'"
            usage;;
    esac
}

if [ "$#" = "2" ];then
    case $1 in
        "kubernetes") kubernetes $2;;
        "swarm") dockerswarm $2;;
        "dockercompose") dockercompose $2;;
        *) 
            echo "invalid mode parameter '$1'"
            usage
            ;;
    esac
else
    echo "invalid parameter"
    usage
fi
