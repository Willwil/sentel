### Kafka setting in docker
Get kafka and zookeeper image from docker hub
> docker pull wurstmeister/zookeeper  
> docker pull wurstmeister/kafka  
   
Start zookeeper and kafka  
> docker run -d --name zookeeper -p 2181 -t wurstmeister/zookeeper  
> docker run --name kafka -e HOST\_IP=localhost -e KAFKA\_ADVERTISED\_PORT=9092 -e KAFKA\_BROKER\_ID=1 -e ZK=zk -p 9092 --link zookeeper:zk -t wurstmeister/kafka  

Get container id  ${CONTAINER ID} using `docker ps `
> docker exec -it ${CONTAINER ID} /bin/bash  
> cd opt/kafka\_2.11-0.10.1.1/   

Testing Kafka  
> bin/kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 1 --topic mykafka   
> bin/kafka-console-producer.sh --broker-list localhost:9092 --topic mykafka   
> bin/kafka-console-consumer.sh --zookeeper zookeeper:2181 --topic mykafka --from-beginning 
