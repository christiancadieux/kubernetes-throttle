
IMAGE=rdei-throttle
RELEASE=v1.0.0

bin:
	 CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o throttle_server ./cmd/...


image: bin
	docker build -t ${IMAGE}:${RELEASE} .
	docker push ${IMAGE}:${RELEASE} 

