gen:    
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/*.proto
clean:  
	rm  proto/*.go
infra:

up:
	docker-compose -f docker-compose.yml up & $(c)
start:
	 docker-compose -f docker-compose.yml start $(c)
run:   
	start
	go  run cmd/client/main.go
down:
	docker-compose -f docker-compose.yml down -v $(c)

