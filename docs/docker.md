### docker configuration
Get mongo and kafka and zookeeper image from docker hub
> docker pull mongo
> docker pull wurstmeister/zookeeper
> docker pull wurstmeister/kafka
   
Start mongo and zookeeper and kafka with docker-compose
> docker-compose -f docker-compose-single-broker.yml up

Clear docker-compose
> docker-compose -f docker-compose-single-broker.yml rm

Start mongo and zookeeper and kafka with docker
> docker run -d --name mongo -p 27017:27017
> docker run -d --name zookeeper -p 2181 -t wurstmeister/zookeeper  
> docker run -d --name kafka -e HOST_IP=localhost -e KAFKA_ADVERTISED_PORT=9092 -e KAFKA_BROKER_ID=1 -e KAFKA_ADVERTISED_HOST_NAME=127.0.0.1 -e ZK=zk -p 9092:9092 --link zookeeper:zk -t wurstmeister/kafka

Get container id  ${CONTAINER ID} using `docker ps `
> docker exec -it ${CONTAINER ID} /bin/bash  
> cd opt/kafka\_2.11-0.10.1.1/   

Testing Kafka  
> bin/kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 1 --topic mykafka   
> bin/kafka-console-producer.sh --broker-list localhost:9092 --topic mykafka   
> bin/kafka-console-consumer.sh --zookeeper zookeeper:2181 --topic mykafka --from-beginning 
