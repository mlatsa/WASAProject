from fastapi import FastAPI
from fastapi.responses import JSONResponse

app = FastAPI(title="WASA Sample API")

@app.get("/health")
def health():
    return JSONResponse({"status": "ok"})
