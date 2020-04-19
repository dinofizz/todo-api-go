build-push:
	docker buildx build --platform linux/amd64,linux/arm/v7 -t dinofizz/todo-api-go:latest -f Dockerfile.multi-arch --push .

