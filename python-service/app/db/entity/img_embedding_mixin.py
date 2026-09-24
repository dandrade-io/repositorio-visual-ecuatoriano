import uuid

from pgvector.sqlalchemy import Vector
from sqlalchemy.orm import Mapped, mapped_column
from sqlalchemy.orm.decl_api import declared_attr
from sqlalchemy.sql.schema import ForeignKey
from sqlalchemy.sql.sqltypes import UUID


class EmbeddingMixin:
    __embedding_dim__: int

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)

    @declared_attr
    def asset_id(cls) -> Mapped[uuid.UUID]:
        return mapped_column(UUID(as_uuid=True), ForeignKey("asset.id"), nullable=False)

    @declared_attr
    def embedding(cls) -> Mapped[list[float]]:
        return mapped_column(Vector(cls.__embedding_dim__), nullable=False)
