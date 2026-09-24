from io import BytesIO

import httpx
from PIL import Image, ImageOps


async def fetch_and_prepare_image(url: str, target_size: int = 512):
    async with httpx.AsyncClient(timeout=10.0) as client:
        response = await client.get(url)
        response.raise_for_status()

    img = Image.open(BytesIO(response.content))
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
