1.create virtual machine  
> for i in seq 3; do docker-machine create -d virtualbox node-$i; done  

2. setup docker swarm manager
> eval $(docker-machine env node-1)  
> docker swarm init --advertise-addr 192.168.99.100 (node-1's ip address)  

3. add swarm worker nodes  
> for i in 2 3 4; do \  
> docker-machine ssh node-$i sudo docker swarm join  \  
> --token xxxx--xxxx-xxx 192.168.99.100:2377; done  

4. confirm swarm cluster's status  
> docker info 

5.start base services  
> cd $GOPATH/src/github.com/cloustone/sentel 
> docker deploy -c swarm-cluster-base-services.yml  sentel

6. confirm service status  
> docker service ls  
> docker service ps sentel\_kafka   
> docker service logssentel\_kafka

6. test wether the broker service can be started normaly  
> docker service create \  
> --name broker  \
> --replicas 3  \
> --network sentel\_front  
> --env KAFKA\_HOST=sentel\_kafka  \  
> --env MONGO\_HOST=sentel\_mongo \  
> --env BROKER\_TENANT=hello  
> --env BROKER\_PRODUCT=world  
> sentel/broker  

7. test network's status
> docker network ls  
> docker network create --driver overlay  hello
> docker network rm hello  



