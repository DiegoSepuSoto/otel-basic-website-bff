build:
	@docker build . -t diegosepusoto/otel-basic-website-bff:local

start:
	@docker run --platform linux/amd64 \
	-d \
	-p 8082:8082 \
	--name otel-basic-website-bff \
	-e OTEL_TRACES_EXPORTER=otlp \
	-e ORDER_API_HOST=http://host.docker.internal:8081 \
	-e OTEL_EXPORTER_OTLP_ENDPOINT=http://host.docker.internal:4317 \
	-e OTEL_RESOURCE_ATTRIBUTES=service.name=bff,service.version=1.0,deployment.environment=local \
	diegosepusoto/otel-basic-website-bff:local

stop:
	@docker rm -f otel-basic-website-bff

.PHONY: build start