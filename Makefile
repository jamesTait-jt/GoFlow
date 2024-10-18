build-workerpool:
	docker build -t workerpool -f dockerfiles/Dockerfile.workerpool .

build-pluginbuilder:
	docker build -t plugin-builder -f dockerfiles/Dockerfile.pluginbuilder .

build-goflow:
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	cmd/goflow/goflow/goflow.proto

	docker build -t goflow -f dockerfiles/Dockerfile.goflow .

build: build-goflow build-workerpool build-pluginbuilder
