### location-history

A toy implementation of in-memory location history server.

Assumptions:
1. The client is authorized/authenticated to make the request
2. Design client to send dynamic queries (based on CLI args, for example: go run client/main.go -method GET -order_id abc123 -max 4)

To Run:
1. Run server/main.go to start the server
2. Run client/main.go to start the client

Improvements:
1. Implement authorization before processing a request
2. Writing location data to a database (NoSQL db such as MongoDB) instead of a json file
3. More test coverage