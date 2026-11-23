from pydantic import BaseModel, Field, HttpUrl
from pydantic.functional_validators import field_validator

from app.core.registry.models import SUPPORTED_AUDIO_MODELS, SUPPORTED_MODELS


class EmbedAssetRequest(BaseModel):
    asset_url: HttpUrl
    models: list[str] = Field(default_factory=lambda: ["jina-clip-v2"], min_length=1)

    @field_validator("models")
    @classmethod
    def validate_models(cls, v: list[str]) -> list[str]:
        v = list(dict.fromkeys(v))
        unknown = [m for m in v if m not in SUPPORTED_MODELS]
        if unknown:
            raise ValueError(
                f"Unsupported models: {unknown}. Supported: {sorted(SUPPORTED_MODELS)}"
            )

        return v


class EmbedTextRequest(BaseModel):
    text: str
    model: str = "jina-clip-v2"

    @field_validator("model")
    @classmethod
    def validate_model(cls, v: str) -> str:
        supported = SUPPORTED_MODELS | SUPPORTED_AUDIO_MODELS
        if v not in supported:
            raise ValueError(
                f"Unsupported model: {v!r}. Supported: {sorted(supported)}"
            )

        return v


class EmbedResponse(BaseModel):
    results: list[EmbedResult]


class EmbedResult(BaseModel):
    model: str
    embedding: list[float]


class EmbedAudioRequest(BaseModel):
    asset_url: HttpUrl
    models: list[str] = Field(default_factory=lambda: ["clamp3"], min_length=1)

    @field_validator("models")
    @classmethod
    def validate_models(cls, v: list[str]) -> list[str]:
        v = list(dict.fromkeys(v))
        unknown = [m for m in v if m not in SUPPORTED_AUDIO_MODELS]
        if unknown:
            raise ValueError(
                f"Unsupported models: {unknown}. Supported: {sorted(SUPPORTED_AUDIO_MODELS)}"
            )

        return v
