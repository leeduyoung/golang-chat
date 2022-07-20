.DEFAULT_GOAL := start

.PHONY: redis-up
redis-up:
	kubectl port-forward -n redis svc/my-redis-master 6379:6379

.PHONY: test
test: test_setup test_run test_clean	

.PHONY: start
start:
	go run main.go

.PHONY: test_setup
test_setup:
	@echo "building the environment.."
	docker-compose -f docker/docker-compose.yml up --build -d
	@sleep 5
	@echo "environment build is done."

.PHONY: test_run
test_run:
	@echo "started run the all tests."
	go test ./... -v -count 1 -p 1 -timeout 600s ./...
	@echo "all tests were completed."

.PHONY: test_clean
test_clean:
	@echo "cleaning the environment.."
	docker-compose -f docker/docker-compose.yml down
	@echo "environment cleaned up."