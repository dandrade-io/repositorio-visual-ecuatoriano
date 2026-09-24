from pydantic import BaseModel, HttpUrl


class EmbedRequest(BaseModel):
    asset_url: HttpUrl
    model: str = "jina-clip-v2"


class EmbedTextRequest(BaseModel):
    text: str
    model: str = "jina-clip-v2"


class EmbedResponse(BaseModel):
    model: str
    embedding: list[float]
