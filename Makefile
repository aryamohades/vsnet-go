format:
	gofmt -w .

client:
	docker build -t sock-client -f Dockerfile.client .

clean-client:
	docker rm -vf $$(docker ps -aq --filter label=sock-client)

docker-client:
	docker run -it -l sock-client sock-client /client -conn=10000 -ramp=0 -ip=172.17.0.2
