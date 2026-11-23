from abc import ABC, abstractmethod

import numpy as np
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


class AudioTextModel(ABC):
    """Model that embeds audio and text into one shared space."""

    sample_rate: int

    @abstractmethod
    def load(self):
        pass

    @abstractmethod
    def encode_text(self, texts: list[str]) -> Tensor:
        pass

    @abstractmethod
    def encode_audio(self, waveforms: list[np.ndarray]) -> Tensor:
        """Encode mono float32 waveforms already resampled to `sample_rate`."""
