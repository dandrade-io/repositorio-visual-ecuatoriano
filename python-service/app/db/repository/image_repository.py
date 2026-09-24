from app.db.entity.image import Image
from sqlalchemy.ext.asyncio import AsyncSession


async def save_image(image: Image, session: AsyncSession):
    session.add(image)
    await session.flush()
