import asyncio

from fastapi import APIRouter

from app.api.v1.dto.embed_dto import (
    EmbedAssetRequest,
    EmbedResponse,
    EmbedResult,
    EmbedTextRequest,
)
from app.core.registry.vision_registry import vision_registry
from app.core.utils import fetch_and_prepare_image

router = APIRouter()


@router.post("/image", response_model=EmbedResponse)
async def embed_image(req: EmbedAssetRequest):
    image = await fetch_and_prepare_image(str(req.asset_url))

    async def embed_one(name: str) -> EmbedResult:
        model = vision_registry.get(name)
        vector = (await asyncio.to_thread(model.encode_image, [image]))[0]
        return EmbedResult(model=name, embedding=vector.tolist())

    results = await asyncio.gather(*[embed_one(m) for m in req.models])
    return EmbedResponse(results=results)


@router.post("/text", response_model=EmbedResult)
def embed_text(req: EmbedTextRequest):
    model = vision_registry.get(req.model)
    vector = model.encode_text([req.text])[0]

    return EmbedResult(model=req.model, embedding=vector.tolist())
