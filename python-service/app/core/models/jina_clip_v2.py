from PIL.Image import Image
from torch._prims_common import Tensor
from transformers import AutoModel

from app.core.models.base import VisionLanguageModel


class JinaClipV2(VisionLanguageModel):
    def __init__(self, truncate_dim: int = 256):
        self._model: None
        self._truncate_dim = truncate_dim

    def load(self):
        self._model = AutoModel.from_pretrained(
            "jinaai/jina-clip-v2", trust_remote_code=True
        )

    def encode_text(self, texts: list[str]) -> Tensor:
        assert self._model is not None
        return self._model.encode_text(
            texts, task="retrieval.query", truncate_dim=self._truncate_dim
        )

    def encode_image(self, images: list[Image]) -> Tensor:
        assert self._model is not None
        return self._model.encode_image(images, truncate_dim=self._truncate_dim)
