bindata:
	go get github.com/jteeuwen/go-bindata/...
	cd protoc-gen-nest/template; go-bindata -o ../template.go -pkg main  ./...;
install:
	go get -u github.com/nestgo/nest/protoc-gen-nest
	go get -u github.com/nestgo/nest/cmd/nest
	go get -u github.com/nestgo/nest
