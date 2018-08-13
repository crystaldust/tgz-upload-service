SHELL := /bin/bash
#	go build -ldflags '-w -extldflags "-static"' -o testbin
all:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-w' -o tgz-upload-service
	docker build -t servicecomb/jmeter-collector:v1 ./
	distribute-image.sh servicecomb/jmeter-collector:v1

	# istioctl kube-inject -f ./testbin.yaml > testbin.injected.yaml
	kubectl apply -f ./jmeter-collector.yaml
