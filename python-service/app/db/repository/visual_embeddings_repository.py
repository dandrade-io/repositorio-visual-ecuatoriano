from sqlalchemy.ext.asyncio.session import AsyncSession

from app.db.entity.img_embedding_mixin import EmbeddingMixin


async def save_embedding(embedding: EmbeddingMixin, session: AsyncSession):
    session.add(embedding)
    await session.flush()
