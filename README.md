# p2p-storage

1. build server `go build server.go`
2. build client `go build client.go`

3. run master node `./server -port=1111`
4. run additional node `./server -port=2222 -nodeIp=127.0.0.1 -nodePort=1111`
5. run client `./client.exe -node="127.0.0.1:2222"`
   
Available commands
   
* GET NODES - get list of available nodes
* GET KEY `key` - get value by key
* GET KEYS - get all keys
* SET KEY `key` `value` - set value


Description: Implement a set of two apps - node and client communicating using any
layer 4 protocol. Node application should act as a storage node which is able to store
and receive string key-value pairs. In case of multiple nodes launched, all nodes
should share all keys and values they have including newly added keys on any of the
nodes. Node should be able to bootstrap itself with a single ip:port combination of any
node already running. Client should be able to connect to any node with IP:port and
have an ability to retrieve all values or one by a key.
