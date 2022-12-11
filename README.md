# PubSub_v2
PubSub_v2 is an websocket microservice. It accepts messages from internal services through its inbuilt gRPC server. 

Message body includes, 
* uid - Destination/Target identification id
* Topic - Since I plan to make it an pubsub server and on client side to classify message type, the topic of message
* Message - message boody, typically json
