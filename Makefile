.PHONY: up down logs rebuild test clean

up:
	@if docker compose version >/dev/null 2>&1; then docker compose up -d --build; else docker-compose up -d --build; fi

down:
	@if docker compose version >/dev/null 2>&1; then docker compose down -v; else docker-compose down -v; fi

logs:
	@if docker compose version >/dev/null 2>&1; then docker compose logs -f --tail=100; else docker-compose logs -f --tail=100; fi

rebuild:
	@if docker compose version >/dev/null 2>&1; then docker compose build --no-cache && docker compose up -d; else docker-compose build --no-cache && docker-compose up -d; fi

test:
	@echo "Backend health:"; curl -s http://localhost:8080/healthz && echo
	@echo "Frontend health:"; curl -s http://localhost:8081/healthz && echo
	@echo "API smoke (create convo, send msg):"
	@CID=chat_smoke; \
	curl -s http://localhost:8080/api/conversations/$$CID > /dev/null; \
	curl -s -X POST "http://localhost:8080/api/conversations/$$CID/messages" \
	  -H "Authorization: Bearer bearer-demo-token" -H "Content-Type: application/json" \
	  -d '{"content":"hello from Makefile smoke"}' | jq -r .id

clean:
	@if docker compose version >/dev/null 2>&1; then docker compose down -v; else docker-compose down -v; fi
	docker system prune -f
