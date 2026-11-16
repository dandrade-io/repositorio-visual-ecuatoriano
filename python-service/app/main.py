import asyncio
from contextlib import asynccontextmanager

from fastapi.applications import FastAPI

from app.api.v1.router import api_v1_router
from app.core.registry.vision_registry import vision_registry

preload_models: list[str] = ["jina-clip-v2", "siglip2"]


@asynccontextmanager
async def lifespan(app: FastAPI):
    for name in preload_models:
        await asyncio.to_thread(vision_registry.get, name)

    yield


app = FastAPI(lifespan=lifespan)
app.include_router(api_v1_router)
