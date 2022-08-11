# build go exec
# create gdfs client exec 
.PYONY: build.client
build.client:
	@-echo build client...
	@go build -o gdfs cmd/client/main.go
	@-echo success!


# create cobra command doc files
.PYONY: cmd.gendocs
cmd.gendocs: 
	@-echo gendocs...
	@go run cmd/gendocs/main.go
	@-echo success!



# gen xx.pb.go file from xx.proto file
# gen all protobuf to golang file
.PYONY: gen.proto.all
gen.proto.all: proto/datanode/datanode.proto proto/namenode/namenode.proto


proto=proto
.PYONY: gen.proto 
gen.proto: proto/${proto}/proto