# SMS Bird

Simple SMS client with request throttling

## Installing

```
brew install go
go get glide
glide install
```

## Running
```
go build && ./SmsBird -key <your_key>
```

## Using
```
curl http://localhost:8080 -X POST -d "{\"recipient\":31612345678,\"originator\":\"MessageBird\",\"message\":\"This is a test message.\"}"
```

## Testing
```
go test -v $(go list ./... | grep -v /vendor/)
```
