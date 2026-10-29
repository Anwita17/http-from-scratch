CLUSTER := myhttp
IMAGE := httpserver:dev

.PHONY: image deploy-local

image:
	docker build -t $(IMAGE) .

deploy-local: image
	kind load docker-image $(IMAGE) --name $(CLUSTER)
	kubectl apply -f deploy/deployment.yaml -f deploy/ingress.yaml -f deploy/service.yaml -f deploy/vm-scrape.yaml
	kubectl rollout restart deployment/httpserver
	kubectl rollout status deployment/httpserver