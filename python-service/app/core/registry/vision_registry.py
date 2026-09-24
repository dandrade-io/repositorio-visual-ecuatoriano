from app.core.models.base import VisionLanguageModel
from app.core.models.jina_clip_v2 import JinaClipV2
from app.core.registry.registry import ModelRegistry
from app.db.entity.visual_embeddings import EmbeddingJinaFull

vision_registry: ModelRegistry[VisionLanguageModel] = ModelRegistry(
    model_classes={
        "jina-clip-v2": JinaClipV2,
    }
)

embedding_entity_classes = {
    "jina-clip-v2": EmbeddingJinaFull,
    # "siglip2-so400m": EmbeddingSiglip2So400m,
}
