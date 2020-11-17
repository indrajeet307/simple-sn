DIR=$(PWD)
IMAGE_TAG=huddle-image:latest

test: docker
	docker run --rm -it -v ${DIR}:/app -w /app ${IMAGE_TAG} go test . -v

build: docker
	docker run --rm -it -v ${DIR}:/app -w /app ${IMAGE_TAG} go build .

docker:
	docker build . -t ${IMAGE_TAG}

clean: clean-cache
	rm ./huddle-app

clean-cache:
	docker run --rm -it -v ${DIR}:/app -w /app ${IMAGE_TAG} rm -rf ./.cache
