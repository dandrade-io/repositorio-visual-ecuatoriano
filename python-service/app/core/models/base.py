from abc import ABC, abstractmethod

from PIL.Image import Image
from torch import Tensor


class VisionLanguageModel(ABC):
    @abstractmethod
    def load(self):
        pass

    @abstractmethod
    def encode_text(self, texts: list[str]) -> Tensor:
        pass

    @abstractmethod
    def encode_image(self, images: list[Image]) -> Tensor:
        pass
