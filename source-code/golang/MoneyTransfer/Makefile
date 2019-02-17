build_read:
	echo "build query service"
	go build -o ./money-transfer-view-service/query ./money-transfer-view-service/query-service.go
build_write:
	echo "build command service"
	go build -o ./money-transfer-service/command ./money-transfer-service/money-transfer-controller.go
run: build_read build_write
	docker-compose up