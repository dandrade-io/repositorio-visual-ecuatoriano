import asyncio
from typing import Annotated

import httpx
from fastapi import APIRouter, File, Form, HTTPException, UploadFile

from app.api.v1.dto.embed_dto import (
    EmbedAudioRequest,
    EmbedResponse,
    EmbedResult,
)
from app.core.registry.audio_registry import audio_registry
from app.core.registry.models import SUPPORTED_AUDIO_MODELS
from app.core.utils import MAX_AUDIO_BYTES, fetch_audio_bytes, prepare_audio

router = APIRouter()


def _embed_audio(name: str, data: bytes) -> EmbedResult:
    model = audio_registry.get(name)
    try:
        wav = prepare_audio(data, model.sample_rate)
    except Exception as e:
        raise HTTPException(status_code=422, detail=f"Cannot decode audio: {e}")

    vector = model.encode_audio([wav])[0]
    return EmbedResult(model=name, embedding=vector.tolist())


@router.post("/audio", response_model=EmbedResponse)
async def embed_audio(req: EmbedAudioRequest):
    try:
        data = await fetch_audio_bytes(str(req.asset_url))
    except (httpx.HTTPError, ValueError) as e:
        raise HTTPException(status_code=422, detail=f"Cannot fetch audio: {e}")

    results = await asyncio.gather(
        *[asyncio.to_thread(_embed_audio, m, data) for m in req.models]
    )
    return EmbedResponse(results=list(results))


@router.post("/audio/upload", response_model=EmbedResult)
async def embed_audio_upload(
    file: Annotated[UploadFile, File()], model: Annotated[str, Form()]
):
    if model not in SUPPORTED_AUDIO_MODELS:
        raise HTTPException(status_code=422, detail=f"Unsupported model: {model}")

    # Read at most one byte over the limit so an oversized file is never fully loaded.
    data = await file.read(MAX_AUDIO_BYTES + 1)
    if len(data) > MAX_AUDIO_BYTES:
        raise HTTPException(
            status_code=413,
            detail=f"Audio file too large (max {MAX_AUDIO_BYTES // 1024 // 1024} MB)",
        )

    return await asyncio.to_thread(_embed_audio, model, data)

