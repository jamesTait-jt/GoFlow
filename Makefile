build-workerpool:
	docker build -t workerpool -f dockerfiles/Dockerfile.workerpool .

build-pluginbuilder:
	docker build -t plugin-builder -f dockerfiles/Dockerfile.pluginbuilder .

build-goflow:
	docker build -t goflow -f dockerfiles/Dockerfile.goflow .

build-all:
	docker build -t workerpool -f dockerfiles/Dockerfile.workerpool .
	docker build -t plugin-builder -f dockerfiles/Dockerfile.pluginbuilder .
	docker build -t goflow -f dockerfiles/Dockerfile.goflow .
