tests:
	go get .
	go test -cover ./...
	
start-kapacitor:
	docker-compose -f infra/docker-compose.yml up -d

sample1:
	go run kapacitorunit.go -dir ./sample/tick_scripts -tests ./sample/test_cases/test_case.yaml

sample1_debug:
	go run kapacitorunit.go -dir ./sample/tick_scripts -tests ./sample/test_cases/test_case.yaml -stderrthreshold=INFO

sample1_batch:
	go run kapacitorunit.go -dir ./sample/tick_scripts -tests ./sample/test_cases/test_case_batch.yaml

sample1_batch_debug:
	go run kapacitorunit.go -dir ./sample/tick_scripts -tests ./sample/test_cases/test_case_batch.yaml -stderrthreshold=INFO

sample_dir:
	go run kapacitorunit.go -dir ./sample/tick_scripts -tests ./sample/test_cases