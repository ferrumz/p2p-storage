# p2p-storage

1. build server `go build server.go`
2. build client `go build client.go`

3. run master node `./server -port=1111`
4. run additional node `./server -port=2222 -nodeIp=127.0.0.1 -nodePort=1111`
5. run client `.\client.exe node="127.0.0.1:2222"`
   
Available commands
   
* GET NODES - get list of available nodes
* GET KEY `key` - get value by key
* GET KEYS - get all keys
* SET KEY `key` `value` - set value
