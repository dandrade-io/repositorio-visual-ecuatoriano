from fastapi.applications import FastAPI

from app.api.v1.router import api_v1_router

app = FastAPI()
app.include_router(api_v1_router)
