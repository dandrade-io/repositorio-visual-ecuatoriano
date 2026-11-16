from PIL.Image import Image
from torch import Tensor
from transformers import AutoModel

from app.core.models.base import VisionLanguageModel


class JinaClipV2(VisionLanguageModel):
    def __init__(self, truncate_dim: int = 256):
        self._model: None
        self._truncate_dim = truncate_dim

    def load(self):
        sha: str = "e10d47f5691d0454a0fb5d13f46f2199b74cb436"

        self._model = AutoModel.from_pretrained(
            "jinaai/jina-clip-v2",
            trust_remote_code=True,
            revision=sha,
        )

    def encode_text(self, texts: list[str]) -> Tensor:
        assert self._model is not None
        return self._model.encode_text(
            texts, task="retrieval.query", truncate_dim=self._truncate_dim
        )

    def encode_image(self, images: list[Image]) -> Tensor:
        assert self._model is not None
        return self._model.encode_image(images, truncate_dim=self._truncate_dim)
