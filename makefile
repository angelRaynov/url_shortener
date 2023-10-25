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

shorten:
	@URL=$(URL) \
	curl -X POST http://localhost:8080/shorten \
	-H 'Content-Type: text/plain;charset=UTF-8' \
	-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJvYiIsImVtYWlsIjoidGVzdGJvYkBtYW5zaW9uLmNvbSIsInVpZCI6IjJiYTNhZjVmLTNhMWItNGUyYi1hZjczLWJiNzkzM2NjM2Y3YSIsImV4cCI6MTY5ODIyOTY1Nn0.HcyTaoSYhAYiqnbPH7VozSOMr5QlmS_a52IjOToRBcc' \
	--data-raw '{"long_url":"$(URL)", "title":"testTitle"}'

my:
	@URL=$(URL) \
	curl http://localhost:8080/my \
	-H 'Content-Type: text/plain;charset=UTF-8' \
	-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJvYiIsImVtYWlsIjoidGVzdGJvYkBtYW5zaW9uLmNvbSIsInVpZCI6IjJiYTNhZjVmLTNhMWItNGUyYi1hZjczLWJiNzkzM2NjM2Y3YSIsImV4cCI6MTY5ODIyOTY1Nn0.HcyTaoSYhAYiqnbPH7VozSOMr5QlmS_a52IjOToRBcc' \

redirect:
	@URL=$(URL) \
	curl -X GET $(URL) \
	-H 'Content-Type: text/plain;charset=UTF-8' \
	-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJvYiIsImVtYWlsIjoidGVzdGJvYkBtYW5zaW9uLmNvbSIsInVpZCI6IjJiYTNhZjVmLTNhMWItNGUyYi1hZjczLWJiNzkzM2NjM2Y3YSIsImV4cCI6MTY5ODIyOTY1Nn0.HcyTaoSYhAYiqnbPH7VozSOMr5QlmS_a52IjOToRBcc' \
#	--data-raw '{"short_url":"$(URL)"}'

register_alice:
	curl -X POST localhost:8080/register -d '{"username":"alice", "password":"alice134312", "email":"test@mansion.com"}'

register_bob:
	curl -X POST localhost:8080/register -d '{"username":"bob", "password":"bob123456", "email":"testbob@mansion.com"}'

auth_alice:
	curl -X POST localhost:8080/authenticate -d '{"username":"alice", "password":"alice134312"}'

auth_bob:
	curl -X POST localhost:8080/authenticate -d '{"username":"bob", "password":"bob123456"}'

update_bob:
	curl -X PATCH localhost:8080/edit -d '{"username":"bob", "email":"bob@bob.com"}' -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJvYiIsImVtYWlsIjoidGVzdGJvYkBtYW5zaW9uLmNvbSIsInVpZCI6IjQ5YzZmYzI0LWRkYTYtNGM5Ni1hODZiLTYzMWI4M2E3ZWU4YiIsImV4cCI6MTY5NzYyNjA5NX0.9fgfbYwa4vMT7bzL8Od1Ajn5G95l_n_2XozUVUULT2c'

update_alice:
	curl -X PATCH localhost:8080/edit -d '{"username":"alice", "email":"alice@wonderland.com"}' -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFsaWNlIiwiZW1haWwiOiJ0ZXN0QG1hbnNpb24uY29tIiwidWlkIjoiMTYzOTMzYWYtN2MxOS00MTE3LTlkZjEtZjMwMTZlMDY2NTAzIiwiZXhwIjoxNjk3NjE5MjMwfQ.gh7exK9H276jq91guRGlomR0keMx1uKu5W0K4qGQ2mQ'
