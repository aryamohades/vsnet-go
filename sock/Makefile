TAG=sock

build:
	docker build -t $(TAG) -f Dockerfile ..

clean:
	docker rm -vf $$(docker ps -aq -f label=$(TAG))

run:
	docker run -it -l $(TAG) --ulimit nofile=2147483584:2147483584 --sysctl net.netfilter.nf_conntrack_max=2147483584 $(TAG)
