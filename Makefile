###  =================================== Env  ==================================================
bin = bin #binary file directory
cmd = cmd #entry file direcotory


###  =================================== Command  ==================================================
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



## start rpc server
.PYONY: run.client
run.client: bin/client
	@-echo run client...
	@./bin/client


port = 50051
.PYONY: run.namenode
run.namenode: bin/namenode
	@-echo run namenode...
	@./bin/namenode --port=${port}


.PYONY: run.datanode
run.datanode: bin/datanode 
	@-echo run datanode...
	@./bin/datanode




###  =================================== Build File ==================================================

# create client binary file
bin/client: cmd/client/main.go
	@-echo build client...
	@go build -o bin/client cmd/client/main.go
	@-echo success!

# create namenode binary file
bin/namenode: cmd/namenode/main.go
	@-echo build namenode...
	@go build -o bin/namenode cmd/namenode/main.go
	@-echo success!

# create datandoe binary file
bin/datanode: cmd/datanode/main.go
	@-echo build datanode...
	@go build -o bin/datanode cmd/datanode/main.go
	@-echo success!


.PYONY: clean
clean:
	@-echo start clean...
	@rm -rf bin/*
	@-echo clean over!