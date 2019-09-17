run:
	go run cmd/api-server/main.go -nats-user=user -nats-secret=secret

nats:
	nats-server -m 8222 -user user -pass secret