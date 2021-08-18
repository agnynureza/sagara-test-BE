docker.build:
	docker build -t fiber .

docker.run:
	docker run --rm -d \
	--name dev-fiber \
	-p 5000:5000 \
	fiber

docker.fiber: docker.build docker.run

docker.stop:
	docker stop dev-fiber