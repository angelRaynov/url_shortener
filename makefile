build:
	docker-compose up -d && docker exec -it url_shortener bash

up:
	docker-compose up -d

down:
	docker-compose down -v

exec:
	docker exec -it url_shortener bash

run:
	go run main.go

redirect:
	curl http://localhost:1234/$(URL_ID)

shorten:
	@URL=$(URL) \
	curl -X POST http://localhost:8080/shorten \
	-H 'Content-Type: text/plain;charset=UTF-8' \
	-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFsaWNlIiwiZXhwIjoxNjg4MjI1NzQ3fQ.UWq27-0q8Xg5_kgtgkU6RH2MaErjAn7PyFgBbG7RC_Q' \
	--data-raw '{"long_url":"$(URL)"}'

expand:
	@URL=$(URL) \
	curl -X POST http://localhost:8080/expand \
	-H 'Content-Type: text/plain;charset=UTF-8' \
	--data-raw '{"short_url":"$(URL)"}'

auth_alice:
	curl -X POST localhost:8081/authenticate -d '{"username":"alice", "password":"alice134312"}'

test:
	go test -v -cover

benchmark_slow:
	go test -bench=BenchmarkSlow -benchmem -cpu=8 -benchtime=100000x

benchmark:
	go test -bench=BenchmarkShortenURLHandler -benchmem -cpu=8 -benchtime=1000000x
