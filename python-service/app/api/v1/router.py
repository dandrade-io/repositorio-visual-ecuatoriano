from fastapi.param_functions import Depends
from fastapi.routing import APIRouter

from app.api.security import require_api_key
from app.api.v1 import embed

api_v1_router = APIRouter(prefix="/api/v1")

api_v1_router.include_router(
    embed.router, prefix="/embed", dependencies=[Depends(require_api_key)]
)
