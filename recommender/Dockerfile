# Base stage - shared
FROM python:3.10-slim AS base

WORKDIR /app
COPY requirements.txt /app/
RUN pip install --no-cache-dir --upgrade pip \
    && pip install --no-cache-dir -r requirements.txt
COPY . /app

# Development image
FROM base AS dev
EXPOSE 8002
CMD ["fastapi", "dev", "main.py", "--host", "0.0.0.0", "--port", "8002"]

# Production image
FROM python:3.10-slim AS production
EXPOSE 8002
CMD ["fastapi", "main.py", "--host", "0.0.0.0", "--port", "8002"]
