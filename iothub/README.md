# sentel-gateway
Sentel iothub 

# pull base image
docker pull centos

# build nginx-mqtt
docker build -t nginx-mqtt .

# run container
docker run -itd -P --name mynginx nginx-mqtt
