## Apiserver ##


- **step 1:** Install golang
- **step 2:** Install mongo and run mongod
- **step 3:** Run apiserver

Product
-


- **Register Product**: 
 
Request:
`curl -d "Name=Hik&Description=Camera&CategoryId=1000" "http://localhost:4145/api/v1/products/3?api-version=v1"`
Response: (Id:0a0503cb-1700-4356-84d8-bf209c660d5d, write down the id for GET product)
`{"requestID":"3d5b3a42-e125-48e6-aa5b-62b0e29ca0da","success":false,"message":"","result":{"Id":"0a0503cb-1700-4356-84d8-bf209c660d5d","Name":"Hik","Description":"7","TimeCreated":"2017-12-14T11:25:46.33124225+08:00","TimeModified":"0001-01-01T00:00:00Z","CategoryId":"","ProductKey":""}}`

- **Get Product**:

Request:
`curl -XGET "http://localhost:4145/api/v1/products/0a0503cb-1700-4356-84d8-bf209c660d5d?api-version=v1"`
Response:
`{"requestID":"83554f2f-32bf-40a5-82d4-d6a3d9f9b345","success":true,"message":"","result":{"Id":"0a0503cb-1700-4356-84d8-bf209c660d5d","Name":"Hik","Description":"7","TimeCreated":"2017-12-14T11:25:46.331+08:00","TimeModified":"0001-01-01T00:00:00Z","CategoryId":"","ProductKey":""}}`

- **Get all devices under the product**

Request:
`curl -XGET "http://localhost:4145/api/v1/products/0a0503cb-1700-4356-84d8-bf209c660d5d/devices?api-version=v1"`
Response:(no devices, to add a device use Register device)
`{"requestID":"0e38855d-ebd1-4320-819c-d3f14b17cfa2","success":true,"message":"","result":[]}`

- **Get all devices under the product(Paging lookup)**

`Page 1` (Page1's LastId is 0)
`curl -XGET "http://localhost:4145/api/v1/products/devices/page?Id=fbd42c39-8844-4014-9092-d072ef2e0ed6&LastId=0&api-version=v1"`
`{"requestID":"fd385079-1db6-420a-bf85-4c2efdc3e19b","success":true,"message":"","result":[{"id":"2ae0662e-9026-46bd-8ee2-e9b8208d1356","status":""},{"id":"2b5d8452-22fa-40c5-ade8-4a63f18845d3","status":""},{"id":"4d2497d6-72c6-46cb-94ab-1c7b1b5a1e3f","status":""},{"id":"5c5d97f4-4196-45cc-a62f-86788cef13a5","status":""},{"id":"764a9544-f040-4119-b966-fa6089ccf580","status":""},{"id":"97df59f4-8616-428a-b75f-81f41d11c305","status":""},{"id":"9f94e51f-8ce5-444b-b80a-d16534994688","status":""},{"id":"a29f58c0-5389-4e86-bad1-e84c4612c262","status":""},{"id":"f492b8de-8cdc-406a-b35a-d608a4137f92","status":""},{"id":"fd385079-1db6-420a-bf85-4c2efdc3e19b","status":""}]}`

`Page 2` (LastId is page1's requestID, requestId 0 means complete.)
`curl -XGET "http://localhost:4145/api/v1/products/devices/page?Id=fbd42c39-8844-4014-9092-d072ef2e0ed6&LastId=fd385079-1db6-420a-bf85-4c2efdc3e19b&api-version=v1"`
`{"requestID":"0","success":true,"message":"","result":[{"id":"fe567fc6-a2de-4f50-812c-ce84a4d5288c","status":""}]}`

- **Get all devices by ProductName under the product**

Request:
`curl -XGET "http://localhost:4145/api/v1/products/Hik/devices?api-version=v1"`

Response:(no devices)
`{"requestID":"0e38855d-ebd1-4320-819c-d3f14b17cfa2","success":true,"message":"","result":[]}`
    
Add a device:

`curl -d "ProductKey=0a0503cb-1700-4356-84d8-bf209c660d5d&DeviceName=c1" "http://localhost:4145/api/v1/devices/3?api-version=v1"`

Query again:(OK)

`curl -XGET "http://localhost:4145/api/v1/products/Hik/devices?api-version=v1"`
`{"requestID":"4055e7b0-a088-41d8-a83c-621b89f33d83","success":true,"message":"","result":[{"id":"149870d8-cd65-4d3c-8295-9da993c213f7","status":""}]}`

- **Get Product by CategoryId** :(Each type of product has an unique Id. e.g:mobile phone)
`curl -XGET "http://localhost:4145/api/v1/products/1001/products?api-version=v1"`
`{"requestID":"98ab131f-815e-43d1-b186-a4902006ea97","success":true,"message":"","result":[{"Id":"fbd42c39-8844-4014-9092-d072ef2e0ed6","Name":"Hik","Description":"DVR","TimeCreated":"0001-01-01T00:00:00Z","TimeModified":"2017-12-15T11:47:11.932+08:00","CategoryId":"","ProductKey":""}]}`

- **Update Product**:
`curl -XPUT -d "Id=fbd42c39-8844-4014-9092-d072ef2e0ed6&Name=Hik&Description=DVR&CategoryId=1001" "http://localhost:4145/api/v1/products/?api-version=v1"`
`{"requestID":"55257036-4f4a-4bed-846e-ac8cc2ac77d2","success":true,"message":"","result":null}`

Query again:(OK)
`curl -XGET "http://localhost:4145/api/v1/products/fbd42c39-8844-4014-9092-d072ef2e0ed6?api-version=v1"`
`{"requestID":"5c8a1005-503d-4e31-9065-b5d836934581","success":true,"message":"","result":{"Id":"fbd42c39-8844-4014-9092-d072ef2e0ed6","Name":"Hik","Description":"DVR","TimeCreated":"0001-01-01T00:00:00Z","TimeModified":"2017-12-15T11:47:11.932+08:00","CategoryId":"","ProductKey":""}}`


- **Delete a product:**(Delete the product and all its devices)
`curl -XDELETE "http://localhost:4145/api/v1/products/0a0503cb-1700-4356-84d8-bf209c660d5d?api-version=v1"`
`{"requestID":"0e2f7362-e534-4f87-a8a8-e3f03fbb7fcd","success":true,"message":"","result":null}`

Quary product:(no product)

`curl -XGET "http://localhost:4145/api/v1/products/Hik/devicesname?api-version=v1"`
`{"requestID":"","success":false,"message":"not found","result":null}`

Device
-

- **Register Device**: 

Request: (ProductKey is created by Register Product)
`curl -d "ProductKey=0a0503cb-1700-4356-84d8-bf209c660d5d&DeviceName=c1" "http://localhost:4145/api/v1/devices/3?api-version=v1"`
Response:
`{"requestID":"","success":true,"message":"","result":{"DeviceId":"149870d8-cd65-4d3c-8295-9da993c213f7","DeviceName":"c1","DeviceSecret":"","DeviceStatus":"","ProductKey":"0a0503cb-1700-4356-84d8-bf209c660d5d","TimeCreated":"2017-12-14T12:09:01.251656528+08:00"}}`

- **Get Device**:
`$curl -XGET "http://localhost:4145/api/v1/devices/149870d8-cd65-4d3c-8295-9da993c213f7?api-version=v1"`
`{"requestID":"","success":true,"message":"","result":{"DeviceId":"149870d8-cd65-4d3c-8295-9da993c213f7","DeviceName":"c1","DeviceSecret":"","DeviceStatus":"","ProductKey":"149870d8-cd65-4d3c-8295-9da993c213f7","TimeCreated":"2017-12-14T14:06:37.792+08:00"}}`

- **Bulk Register Devices:**(Auto create 10 devices by named c1)
`curl -d "ProductKey=fbd42c39-8844-4014-9&Name=c1" "http://localhost:4145/api/v1/devices/10/bulk?api-version=v1"`
`{"requestID":"","success":true,"message":"","result":[{"Id":"764a9544-f040-4119-b966-fa6089ccf580","Name":"c1","ProductId":"","ProductKey":"fbd42c39-8844-4014-9092-d072ef2e0ed6","DeviceStatus":"","DeviceSecret":"","TimeCreated":"2017-12-14T14:26:35.790953255+08:00","TimeModified":"2017-12-14T14:26:35.790953296+08:00"},{"Id":"97df59f4-8616-428a-b75f-81f41d11c305","Name":"c1","ProductId":"","ProductKey":"fbd42c39-8844-4014-9092-d072ef2e0ed6","DeviceStatus":"","DeviceSecret":"","TimeCreated":"2017-12-14T14:26:35.790966088+08:00","TimeModified":"2017-12-14T14:26:35.79096611+08:00"},{"Id":"4d2497d6-72c6-46cb-94ab-1c7b1b5a1e3f","Name":"c1","ProductId":"","ProductKey":"fbd42c39-8844-4014-9092-d072ef2e0ed6","DeviceStatus":"","DeviceSecret":"","TimeCreated":"2017-12-14T14:26:35.790972053+08:00","TimeModified":"2017-12-14T14:26:35.790972074+08:00"},{"Id":"fd385079-1db6-420a-bf85-4c2efdc3e19b","Name":"c1","ProductId":"","ProductKey":"fbd42c39-8844-4014-9092-d072ef2e0ed6","DeviceStatus":"","DeviceSecret":"","TimeCreated":"2017-12-14T14:26:35.790977174+08:00","TimeModified":"2017-12-14T14:26:35.790977196+08:00"},{"Id":"fe567fc6-a2de-4f50-812c-ce84a4d5288c","Name":"c1","ProductId":"","ProductKey":"fbd42c39-8844-4014-9092-d072ef2e0ed6","DeviceStatus":"","DeviceSecret":"","TimeCreated":"2017-12-14T14:26:35.790980619+08:00","TimeModified":"2017-12-14T14:26:35.790980641+08:00"},{"Id":"f492b8de-8cdc-406a-b35a-d608a4137f92","Name":"c1","ProductId":"","ProductKey":"fbd42c39-8844-4014-9092-d072ef2e0ed6","DeviceStatus":"","DeviceSecret":"","TimeCreated":"2017-12-14T14:26:35.790985879+08:00","TimeModified":"2017-12-14T14:26:35.790985901+08:00"},{"Id":"9f94e51f-8ce5-444b-b80a-d16534994688","Name":"c1","ProductId":"","ProductKey":"fbd42c39-8844-4014-9092-d072ef2e0ed6","DeviceStatus":"","DeviceSecret":"","TimeCreated":"2017-12-14T14:26:35.790989269+08:00","TimeModified":"2017-12-14T14:26:35.79098929+08:00"},{"Id":"5c5d97f4-4196-45cc-a62f-86788cef13a5","Name":"c1","ProductId":"","ProductKey":"fbd42c39-8844-4014-9092-d072ef2e0ed6","DeviceStatus":"","DeviceSecret":"","TimeCreated":"2017-12-14T14:26:35.79099266+08:00","TimeModified":"2017-12-14T14:26:35.790992682+08:00"},{"Id":"a29f58c0-5389-4e86-bad1-e84c4612c262","Name":"c1","ProductId":"","ProductKey":"fbd42c39-8844-4014-9092-d072ef2e0ed6","DeviceStatus":"","DeviceSecret":"","TimeCreated":"2017-12-14T14:26:35.790996067+08:00","TimeModified":"2017-12-14T14:26:35.790996088+08:00"},{"Id":"2ae0662e-9026-46bd-8ee2-e9b8208d1356","Name":"c1","ProductId":"","ProductKey":"fbd42c39-8844-4014-9092-d072ef2e0ed6","DeviceStatus":"","DeviceSecret":"","TimeCreated":"2017-12-14T14:26:35.791001202+08:00","TimeModified":"2017-12-14T14:26:35.791001226+08:00"}]}`

- **Get ALL Devices by name**：(name is c1)
`curl -XGET "http://localhost:4145/api/v1/devices/c1/all?api-version=v1"`
`{"requestID":"f8498631-7273-437a-843f-6b920045f3da","success":true,"message":"","result":[{"id":"2b5d8452-22fa-40c5-ade8-4a63f18845d3","status":""},{"id":"764a9544-f040-4119-b966-fa6089ccf580","status":""},{"id":"97df59f4-8616-428a-b75f-81f41d11c305","status":""},{"id":"4d2497d6-72c6-46cb-94ab-1c7b1b5a1e3f","status":""},{"id":"fd385079-1db6-420a-bf85-4c2efdc3e19b","status":""},{"id":"fe567fc6-a2de-4f50-812c-ce84a4d5288c","status":""},{"id":"f492b8de-8cdc-406a-b35a-d608a4137f92","status":""},{"id":"9f94e51f-8ce5-444b-b80a-d16534994688","status":""},{"id":"5c5d97f4-4196-45cc-a62f-86788cef13a5","status":""},{"id":"a29f58c0-5389-4e86-bad1-e84c4612c262","status":""},{"id":"2ae0662e-9026-46bd-8ee2-e9b8208d1356","status":""}]}`

`$curl -XPUT -d "ProductId=8&DeviceStatus=1&DeviceSecret=1&ProductKey=7&DeviceName=2" "http://localhost:4145/api/v1/devices/3?api-version=v1"`
`$curl -d "ProductId=8&DeviceStatus=0&DeviceSecret=1&ProductKey=7&DeviceName=2" "http://localhost:4145/api/v1/devices/3?api-version=v1"`
`$curl -XDELETE "http://localhost:4145/api/v1/devices/3?api-version=v1"`


