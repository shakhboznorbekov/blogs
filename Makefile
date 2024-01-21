push:
	git add .
	git commit -m "some changes"
	git push origin shaxboz

run:
	go run cmd/main.go

swag-gen:
	swag init -g internal/router/router.go -o internal/controller/docs

push-main:
	git add .
	git commit -m "some changes"
	git push origin main
