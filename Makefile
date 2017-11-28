bindata:
	go get github.com/jteeuwen/go-bindata/...
	cd protoc-gen-nest/template; go-bindata -o ../template.go -pkg main  ./...;
cmd:
	go build cmd/main.go
