push:
	git add .
	git commit -m "some changes"
	git push origin shaxboz

run:
	go run cmd/main.go

swag-gen:
	swag init -g internal/router/router.go -o internal/controller/docs

local-deploy:
	GOOS=linux GOARCH=amd64 go build -o main cmd/main.go && scp main root@192.168.0.118:/var/www/blogs/website/backend

push-main:
	git add .
	git commit -m "some changes"
	git push origin main
