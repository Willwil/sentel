# conductor debug #
bin/zookeeper-server-start.sh config/zookeeper.properties &
bin/kafka-server-start.sh config/server.properties

kafka-topics.sh --create --replication-factor 1 --partitions 8 --topic product --zookeeper localhost:2181
kafka-topics.sh --list --zookeeper localhost:2181

kafka-console-consumer.sh --topic product --zookeeper localhost:2181 --from-beginning

kafka-console-producer.sh --broker-list localhost:9092 --topic product

