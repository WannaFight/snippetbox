help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

.PHONY: mysql-connect-root
mysql-connect-root: ## Connect to MySQL as root
	@docker compose exec mysql mysql --user=root --password=snippetboxadmin

.PHONY: mysql-connect-web
mysql-connect-web: ## Connect to MySQL as web
	@docker compose exec mysql mysql --user=web --password=webpass
