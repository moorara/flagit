.PHONY: run
run:
	go build -o app && \
	./app \
	  -verbose \
	  -port 8080 \
	  -log-level=info \
	  --config-timeout 60s \
	  --config-endpoints=service-1:8080,service-2:8080 && \
	rm app
