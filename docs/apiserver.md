# apiserver debug #
**step 1:** Install golang

**step 2:** Install mongo and run mongodb

**step 3:** Run senter-apiserver

**step 4:** test
>curl -d "Description=7&Name=mytest&CategoryId=1" "http://localhost:4145/api/v1/products/3?api-version=v1"

>curl -XGET "http://localhost:4145/api/v1/products/mytest?api-version=v1"

>curl -XGET "http://localhost:4145/api/v1/devices/3?api-version=v1"

>curl -XPUT -d "ProductId=8&DeviceStatus=1&DeviceSecret=1&ProductKey=7&DeviceName=2" "http://localhost:4145/api/v1/devices/3?api-version=v1"

>curl -d "ProductId=8&DeviceStatus=0&DeviceSecret=1&ProductKey=7&DeviceName=2" "http://localhost:4145/api/v1/devices/3?api-version=v1"

>curl -XDELETE "http://localhost:4145/api/v1/devices/3?api-version=v1"


