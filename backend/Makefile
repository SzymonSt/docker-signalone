IMAGE?=322456/signalonebackend
TAG?=dev
EXTENSION_NETWORK?=322456_signalone-extension-desktop-extension_default #for local development only

build-backend: ## Build service image - dev version
	docker build --tag=$(IMAGE):$(TAG) .

start-backend: ## Start service locally (local development only)
	docker compose -f ./docker-compose.yaml up -d
	docker network connect $(EXTENSION_NETWORK) backend-backend-1

start-backend-with-init: ## Start service locally with mock data initialization
	docker compose -f ./docker-compose.yaml up -d
	python -m venv ./devenv/initenv.venv
	powershell ./devenv/initenv.venv/Scripts/activate.ps1
	pip install -r ./devenv/requirements.txt
	python ./devenv/init.py

deploy-backend: ## Deploy service to production
	echo "NotImplemented"