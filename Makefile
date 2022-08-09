complete -W "\`grep -oE '^[a-zA-Z0-9_.-]+:([^=]|$)' ?akefile | sed 's/[^a-zA-Z0-9_.-]*$//'\`" make

help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

ecosystem_up: ## Запуск микросервисов
	docker-compose -f ecosystem.yml up --build

ecosystem_down: ## Остановка микросервисов
	docker-compose -f ecosystem.yml down

build: ## Сборка проекта
	docker-compose build 