name = shrtnr
project = deltalabs-xyz
compute_zone = europe-west2-b
image_repository = lackerman

.PHONY: clean docker-build docker-push setup clean-deploy newpods

build: test main.go
	CGO_ENABLED=0 GOOS=linux go build -o bin/shrtnr main.go

clean:
	rm bin/shrtnr

test:
	go test ./...

docker-build: Dockerfile build
	docker build -t $(docker_image):$(latest_commit) .

docker-push: docker-build
	docker push $(docker_image):$(latest_commit)

rollout: docker-push
	cat kubernetes/deployment.yml | sed -e 's/LATEST_COMMIT/$(latest_commit)/g' | kubectl apply -f -

newpods: docker-push
	kubectl delete pod -l app=$(name)

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

docker_image = $(image_repository)/$(name)
deployment_id = $(shell kubectl get deployments -l app=$(name) --output=jsonpath='{.items[0].metadata.name}')
pod_id = $(shell kubectl get pods -l app=$(name) --output=jsonpath='{.items[0].metadata.name}')
latest_commit = $(shell git log -n 1 --pretty=format:"%H")
