name = shrtnr
project = deltalabs-xyz
compute_zone = europe-west2-b

.PHONY: clean test build docker-build docker-push setup clean-deploy newpod

build: test cmd/main.go
	go build -o bin/shrtnr cmd/main.go

clean:
	rm bin/shrtnr

test:
	go test ./...

docker-build: Dockerfile build
	docker build -t $(docker_image) .

docker-push: docker-build
	docker push $(docker_image)

rollout: docker-push
	kubectl apply -f kubernetes/deployment.yml

newpod: package docker-build docker-push
	kubectl delete pod $(pod_id)

clean-deploy:
	kubectl delete all -l app=$(name)

setup:
	gcloud config set project $(project)
	gcloud config set compute/zone $(compute_zone)
	gcloud container clusters create $(project) \
		--zone $(compute_zone) \
		--node-locations $(compute_zone) \
		--num-nodes 2 \
		--no-enable-basic-auth \
		--issue-client-certificate

docker_image = eu.gcr.io/$(project)/$(name)
deployment_id = $(shell kubectl get deployments -l app=$(name) --output=jsonpath='{.items[0].metadata.name}')
pod_id = $(shell kubectl get pods -l app=$(name) --output=jsonpath='{.items[0].metadata.name}')
