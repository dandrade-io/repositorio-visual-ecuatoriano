from app.core.models.base import AudioTextModel
from app.core.models.clamp3 import Clamp3
from app.core.registry.registry import ModelRegistry

audio_registry: ModelRegistry[AudioTextModel] = ModelRegistry(
    model_classes={"clamp3": Clamp3}
)
