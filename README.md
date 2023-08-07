# file-pubsub home assignment
This is a Golang solution for task: "Queue between a file reader and writer".

# How to run
## test
```make test```
## container
```make build-container && make run-container```
## Go
```go run cmd/file-pubsub/main.go```

# Notes
Service is contenarized and can be deployed on some cluster, like kubernetes.
One need to just build container and apply the image to deployment.

Service uses REST API, but a gRPC one could be considered as well.

It's using an InMemory PubSub implementation to queue the messages. On each call, a single file is written and consumed.
I tried to abstract this implementation, so that a service could scale with another Queuing technology (like Kafka).
Another approach would be to have some sort of a service cache and store the files with a expiration date.
However, I feel like the InMem PubSub is more suitable for this task, and it can be replaced later with more robust implementation.


# TODO (run out of time)
* Test more thoroughly --- as of now it covers happy paths and some error checking
* Document services and components better
* Handle context cancellation, timeouts, etc. (It's just passed along right now)