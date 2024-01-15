IMAGE?=322456/signalonebackend
TAG?=dev

build-backend: ## Build service image - dev version
	docker build --tag=$(IMAGE):$(TAG) .

start-backend: ## Start service locally
	docker-compose -d -f ./docker-compose.yaml up
	python ./devenv/init.py

deploy-backend: ## Deploy service to production
	echo 