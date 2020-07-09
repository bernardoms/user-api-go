.PHONY: dependency unit-test integration-test docker-up docker-down clear
export mongoURI=mongodb://local:local@localhost:27017/local?authSource=admin
export database=local
export sns_topic=arn:aws:sns:us-east-1:000000000000:user_update_notify
export aws_region=us-east-1
export endpoint=http://localhost:4575

dependency:
	@cd cmd/; go mod download


integration-test: docker-up dependency
	@cd test/integration; go test ./... -v

unit-test: dependency
	@cd test/unit; go test ./... -v -short

docker-up:
	@cd deps/; docker-compose up -d mongo localstack
	@sleep 10

docker-down:
	@cd deps/; docker-compose stop

clear: docker-down
