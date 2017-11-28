# nest
Nest is a framework for build microserver🚀 

## install

```
go get -u github.com/nestgo/nest/protoc-gen-nest
go get -u github.com/nestgo/nest/cmd/nest
go get -u github.com/nestgo/nest
```

## usage

Automatically generate an grpc application 

```
nest gen --output=app --proto=./helloworld.proto
```
### TIPS:
protobuf file package must proto  
