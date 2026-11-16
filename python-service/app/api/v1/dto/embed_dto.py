from pydantic import BaseModel, Field, HttpUrl
from pydantic.functional_validators import field_validator

from app.core.registry.models import SUPPORTED_MODELS


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


class EmbedResponse(BaseModel):
    results: list[EmbedResult]


class EmbedResult(BaseModel):
    model: str
    embedding: list[float]
