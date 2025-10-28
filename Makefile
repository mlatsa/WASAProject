.PHONY: up down logs rebuild test clean

up:
	docker compose up -d --build

down:
	docker compose down -v

logs:
	docker compose logs -f --tail=100

rebuild:
	docker compose build --no-cache && docker compose up -d

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
	docker compose down -v
	docker system prune -f
