service-dependencies:
	docker run --network=host --env STORAGE=elasticsearch --env ES_NODES=http://0.0.0.0:9200 --env ES_USERNAME=admin --env ES_PASSWORD=admin jaegertracing/spark-dependencies