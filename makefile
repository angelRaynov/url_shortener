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
	curl http://localhost:8080/ \
	-H 'Content-Type: text/plain;charset=UTF-8' \
	--data-raw '{"long_url":"✪df.ws/1234ssaad"}'

test:
	go test -v -cover

benchmark_slow:
	go test -bench=BenchmarkSlow -benchmem -cpu=8 -benchtime=100000x

benchmark:
	go test -bench=BenchmarkShortenURLHandler -benchmem -cpu=8 -benchtime=1000000x
