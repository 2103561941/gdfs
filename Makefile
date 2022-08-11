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
.PYONY: gen.proto
gen.proto:
	@-echo build proto...
	@protoc --go_out=. --go_opt=paths=source_relative \
	 --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	 ./proto/datanode/datanode.proto
	@protoc --go_out=. --go_opt=paths=source_relative \
	 --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	 ./proto/namenode/namenode.proto
	@-echo success!
