# WASA Project Starter

This repository is a clean starter for the **WASA – Web and Software Architecture** homeworks.

It follows the delivery page requirements (single private repo, `doc/api.yaml`, default branch `main`, deploy key (read-only)).

## Quickstart

```bash
# 1) Create a new private repo (GitHub or GitLab) and push this starter:
git init
git branch -M main
git add .
git commit -m "Init WASA starter"
# then create a new empty private repo online and run:
git remote add origin <YOUR_REMOTE_URL>
git push -u origin main
```

### Add the Deploy Key (read‑only)
1. Visit the WASA enroll dashboard and generate/copy your SSH public key.
2. **GitHub:** Settings → *Deploy keys* → **Add key** (read‑only).
   **GitLab:** Settings → *Repository* → *Deploy keys* → **Add key** (read‑only).
3. Make sure your default branch is **main** (or **master**) — grading pulls from remote HEAD.

---

## Project Layout

```
doc/
  api.yaml           # OpenAPI 3.0 (Homework 1)

backend/
  src/app.py         # FastAPI app + OpenAPI served from /openapi.json
  tests/             # pytest samples
  Dockerfile.backend # Homework 4

frontend/
  src/main.html      # ultra‑minimal vanilla app
  Dockerfile.frontend
```

You can replace the starter tech stack if your instructor prefers something else, but keep paths and filenames consistent with the rules.

---

## Develop

```bash
# Backend (FastAPI + uvicorn)
python -m venv .venv && source .venv/bin/activate
pip install -r backend/requirements.txt
uvicorn src.app:app --reload --port 8000 --app-dir backend

# Frontend (static)
python -m http.server 5173 -d frontend/src
# open http://localhost:5173
```

### Tests
```bash
pytest
```

### Docker (Homework 4)
```bash
# backend
docker build -f Dockerfile.backend -t wasa-backend:dev .
docker run -p 8000:8000 wasa-backend:dev

# frontend
docker build -f Dockerfile.frontend -t wasa-frontend:dev .
docker run -p 5173:80 wasa-frontend:dev
```

---

## Makefile

Common tasks are wrapped in a Makefile:
```bash
make setup       # install backend deps
make run-backend # run uvicorn
make run-frontend
make test
make docker-backend
make docker-frontend
```
