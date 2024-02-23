IMAGE?=322456/signalonebackend
TAG?=dev

SOLUTIONS_AGENT_IMAGE?=322456/solutions-agent
SOLUTIONS_AGENT_TAG?=dev

build-backend: ## Build service image - dev version
	docker build --tag=$(IMAGE):$(TAG) .
	docker build --tag=$(SOLUTIONS_AGENT_IMAGE):$(SOLUTIONS_AGENT_TAG) ../solutionAgent

start-backend: ## Start service locally (local development only)
	docker compose -f ./docker-compose.yaml up -d

start-backend-with-init: ## Start service locally with mock data initialization
	docker compose -f ./docker-compose.yaml up -d
	python -m venv ./devenv/initenv.venv
	powershell ./devenv/initenv.venv/Scripts/activate.ps1
	pip install -r ./devenv/requirements.txt
	python ./devenv/init.py

deploy-backend: ## Deploy service to production
	echo "NotImplemented"