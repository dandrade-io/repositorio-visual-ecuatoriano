from app.db.entity.base_entity import Base
from app.db.entity.img_embedding_mixin import EmbeddingMixin


class EmbeddingJinaFull(EmbeddingMixin, Base):
    __tablename__ = "embedding_jina_full"
    __embedding_dim__ = 1024
