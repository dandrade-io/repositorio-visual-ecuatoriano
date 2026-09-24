from fastapi import APIRouter

from app.api.v1.dto.embed_dto import EmbedRequest, EmbedResponse, EmbedTextRequest
from app.core.registry.vision_registry import vision_registry
from app.core.utils import fetch_and_prepare_image

router = APIRouter()


@router.post("/image", response_model=EmbedResponse)
async def embed_image(req: EmbedRequest):
    model = vision_registry.get(req.model)
    image = await fetch_and_prepare_image(str(req.asset_url))
    vector = model.encode_image([image])[0]

    return EmbedResponse(model=req.model, embedding=vector.tolist())


@router.post("/text", response_model=EmbedResponse)
async def embed_text(req: EmbedTextRequest):
    model = vision_registry.get(req.model)
    vector = model.encode_text([req.text])[0]

    return EmbedResponse(model=req.model, embedding=vector.tolist())
