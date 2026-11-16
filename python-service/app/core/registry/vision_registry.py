from app.core.models.base import VisionLanguageModel
from app.core.models.jina_clip_v2 import JinaClipV2
from app.core.models.siglip2 import SigLip2
from app.core.registry.registry import ModelRegistry

vision_registry: ModelRegistry[VisionLanguageModel] = ModelRegistry(
    model_classes={"jina-clip-v2": JinaClipV2, "siglip2": SigLip2}
)
