PY=python
PIP=pip

setup:
	$(PY) -m venv .venv && . .venv/bin/activate && $(PIP) install -r backend/requirements.txt pytest

run-backend:
	uvicorn src.app:app --reload --port 8000 --app-dir backend

run-frontend:
	python -m http.server 5173 -d frontend/src

test:
	pytest -q

docker-backend:
	docker build -f Dockerfile.backend -t wasa-backend:dev .

docker-frontend:
	docker build -f Dockerfile.frontend -t wasa-frontend:dev .
