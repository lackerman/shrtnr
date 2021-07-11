name = shrtnr
image_repository = lackerman

.PHONY: setup clean-deploy newpods

release:
	goreleaser --rm-dist

rollout: docker-push
	cat kubernetes/deployment.yml | sed -e 's/LATEST_RELEASE/$(latest_release)/g' | kubectl apply -f -

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
latest_release = $(shell git tag | sort -u -r | head -n1)
