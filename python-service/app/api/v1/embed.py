import asyncio
from typing import Annotated

from fastapi import APIRouter, File, Form, UploadFile

from app.api.v1.dto.embed_dto import (
    EmbedAssetRequest,
    EmbedResponse,
    EmbedResult,
    EmbedTextRequest,
)
from app.core.registry.audio_registry import audio_registry
from app.core.registry.models import SUPPORTED_AUDIO_MODELS
from app.core.registry.vision_registry import vision_registry
from app.core.utils import fetch_and_prepare_image, prepare_image

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


@router.post("/image/upload")
async def embed_image_upload(
    file: Annotated[UploadFile, File()], model: Annotated[str, Form()]
):
    image = prepare_image(await file.read())
    m = vision_registry.get(model)
    vector = (await asyncio.to_thread(m.encode_image, [image]))[0]

    return EmbedResult(model=model, embedding=vector.tolist())


@router.post("/text", response_model=EmbedResult)
async def embed_text(req: EmbedTextRequest):
    registry = audio_registry if req.model in SUPPORTED_AUDIO_MODELS else vision_registry

    def run():
        return registry.get(req.model).encode_text([req.text])[0]

    vector = await asyncio.to_thread(run)

    return EmbedResult(model=req.model, embedding=vector.tolist())
