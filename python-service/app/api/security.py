import secrets

from fastapi import Security, status
from fastapi.security import APIKeyHeader
from starlette.exceptions import HTTPException

from app.config import settings

api_key_header = APIKeyHeader(name="X-API-Key", auto_error=False)


async def require_api_key(provided: str | None = Security(api_key_header)) -> None:
    if provided is None or not secrets.compare_digest(provided, settings.embed_service_api_key):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid or missing API key",
        )
