1. create virtual machine  
> for i in `seq 3`; do docker-machine create -d virtualbox node-$i; done  

2. confirm node-1's ip address
> docker-machine ip node-1

3. setup docker swarm manager (node-1's ip address is 192.168.99.100)
> eval $(docker-machine env node-1)  
> docker swarm init --advertise-addr 192.168.99.100

4. add swarm worker nodes  
> for i in 2 3; do \  
> docker-machine ssh node-$i "docker swarm join  \  
> --token <token> 192.168.99.100:2377"; done

4. confirm swarm cluster's status  
> docker node ls 

5. update docker image registry
> for i in `seq 3`; do \  
> docker-machine ssh node-$i "echo 'EXTRA_ARGS=\"--registry-mirror=https://registry.docker-cn.com\"' | sudo tee -a /var/lib/boot2docker/profile"  

6. restart worker node
> for i in `seq 3`; do docker-machine restart node-$i; done  

6. start base services  
> cd $GOPATH/src/github.com/cloustone/sentel    
> docker stack deploy -c swarm-cluster-base-services.yml sentel

7. confirm service status  
> docker service ls  
> docker service ps sentel\_kafka   
> docker service logssentel\_kafka

8. test wether the broker service can be started normaly  
> docker service create \  
> --name broker  \  
> --replicas 3  \  
> --network sentel\_front  \  
> --env KAFKA\_HOST=sentel\_kafka  \  
> --env MONGO\_HOST=sentel\_mongo \  
> --env BROKER\_TENANT=hello  \  
> --env BROKER\_PRODUCT=world  \  
> sentel/broker  

9. test network's status
> docker network ls  
> docker network create --driver overlay  hello
> docker network rm hello  



