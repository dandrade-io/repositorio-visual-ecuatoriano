import torch
from PIL.Image import Image
from transformers import AutoModel, AutoProcessor

from app.core.models.base import VisionLanguageModel


class SigLip2(VisionLanguageModel):
    def load(self):
        model_id: str = "google/siglip2-large-patch16-512"
        sha: str = "49488218e80259885f3be61d7a9455faf833b7a8"

        self._model = AutoModel.from_pretrained(
            model_id, revision=sha, device_map="auto"
        ).eval()
        self._processor = AutoProcessor.from_pretrained(model_id, revision=sha)

    def encode_text(self, texts: list[str]) -> torch.Tensor:
        inputs = self._processor(
            text=[t.lower() for t in texts], return_tensors="pt"
        ).to(self._model.device)

        with torch.no_grad():
            text_features = self._model.get_text_features(**inputs)

        return text_features / text_features.norm(p=2, dim=-1, keepdim=True)

    def encode_image(self, images: list[Image]) -> torch.Tensor:
        inputs = self._processor(images=images, return_tensors="pt").to(
            self._model.device
        )

        with torch.no_grad():
            embeds = self._model.get_image_features(**inputs)

        return embeds / embeds.norm(dim=-1, keepdim=True)
