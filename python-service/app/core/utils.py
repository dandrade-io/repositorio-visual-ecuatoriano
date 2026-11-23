import math
from io import BytesIO

import httpx
import numpy as np
import soundfile as sf
from PIL import Image, ImageOps
from scipy.signal import resample_poly


async def fetch_and_prepare_image(url: str, target_size: int = 512):
    async with httpx.AsyncClient(timeout=10.0) as client:
        response = await client.get(url)
        response.raise_for_status()

    return prepare_image(response.content, target_size)


def prepare_image(data: bytes, target_size: int = 512):
    img = Image.open(BytesIO(data))
    img = ImageOps.exif_transpose(img)
    img = img.convert("RGB")

    img.thumbnail((target_size, target_size), Image.Resampling.LANCZOS)

    delta_w = target_size - img.width
    delta_h = target_size - img.height
    padding = (
        delta_w // 2,
        delta_h // 2,
        delta_w - delta_w // 2,
        delta_h - delta_h // 2,
    )

    return ImageOps.expand(img, padding, fill=(255, 255, 255))


MAX_AUDIO_BYTES = 100 * 1024 * 1024
MAX_AUDIO_SECONDS = 600


async def fetch_audio_bytes(url: str) -> bytes:
    chunks: list[bytes] = []
    size = 0
    async with httpx.AsyncClient(timeout=30.0, follow_redirects=True) as client:
        async with client.stream("GET", url) as response:
            response.raise_for_status()
            async for chunk in response.aiter_bytes():
                size += len(chunk)
                if size > MAX_AUDIO_BYTES:
                    raise ValueError("Audio file too large")
                chunks.append(chunk)

    return b"".join(chunks)


def prepare_audio(data: bytes, target_sr: int) -> np.ndarray:
    """Decode to mono float32 at `target_sr`, capped to MAX_AUDIO_SECONDS."""
    wav, sr = sf.read(BytesIO(data), dtype="float32", always_2d=True)
    wav = wav.mean(axis=1)

    wav = wav[: int(MAX_AUDIO_SECONDS * sr)]
    if sr != target_sr:
        g = math.gcd(sr, target_sr)
        wav = resample_poly(wav, target_sr // g, sr // g).astype(np.float32)

    if wav.size == 0:
        raise ValueError("Empty audio")

    return wav
