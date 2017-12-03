mongo golang:
http://www.jianshu.com/p/b63e5cfa4ce5

step 1:install golang
step 2:install mongo and run
step 3:run senter-apiserver

curl -d "Description=7&Name=mytest&CategoryId=1" "http://localhost:4145/api/v1/products/3?api-version=v1"
curl -XGET "http://localhost:4145/api/v1/products/mytest?api-version=v1"

curl -XGET "http://localhost:4145/api/v1/devices/3?api-version=v1"
curl -XPUT -d "ProductId=8&DeviceStatus=1&DeviceSecret=1&ProductKey=7&DeviceName=2" "http://localhost:4145/api/v1/devices/3?api-version=v1"
curl -d "ProductId=8&DeviceStatus=0&DeviceSecret=1&ProductKey=7&DeviceName=2" "http://localhost:4145/api/v1/devices/3?api-version=v1"
curl -XDELETE "http://localhost:4145/api/v1/devices/3?api-version=v1"


