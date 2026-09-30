# Stage 1: Build the Python Application
ARG PYTHON_VERSION=3.14-slim

FROM python:${PYTHON_VERSION} AS builder

COPY --from=ghcr.io/astral-sh/uv:0.9 /uv /uvx /bin/

ENV UV_COMPILE_BYTECODE=1 \
    UV_LINK_MODE=copy \
    UV_PYTHON_DOWNLOADS=never

WORKDIR /app


# Install dependencies only (cached until pyproject.toml / uv.lock change)
RUN --mount=type=cache,target=/root/.cache/uv \
    --mount=type=bind,source=uv.lock,target=uv.lock \
    --mount=type=bind,source=pyproject.toml,target=pyproject.toml \
    uv sync --locked --no-install-project --no-dev

# Copy the source and install the project itself
COPY . .
RUN --mount=type=cache,target=/root/.cache/uv \
    uv sync --locked --no-dev

# Stage 2: Runtime
FROM python:${PYTHON_VERSION}

# Non-root user
RUN groupadd --system app && useradd --system --gid app --create-home app

WORKDIR /app

# Copy the app and its virtualenv (same path as the builder, so the venv stays valid)
COPY --from=builder /app /app

ENV PATH="/app/.venv/bin:$PATH" \
    PYTHONUNBUFFERED=1

RUN mkdir -p /home/app/.cache/huggingface && chown -R app:app /home/app/.cache

USER app

EXPOSE 8001

HEALTHCHECK --interval=30s --timeout=5s --start-period=120s --retries=3 \
    CMD python -c "import socket; socket.create_connection(('127.0.0.1', 8001), 3)" || exit 1

ENTRYPOINT ["fastapi", "run", "app/main.py", "--port", "8001"]
