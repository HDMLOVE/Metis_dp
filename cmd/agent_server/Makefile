.PHONLY: server server_client clean
all: clean server server_client

server:
	go build server.go
server_client:
	go build server_client.go
clean:
	rm -rf server server_client
