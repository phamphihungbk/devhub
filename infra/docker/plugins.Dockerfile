
ARG MODE

FROM python:3.11-slim AS base

WORKDIR /app

# COPY plugins/ ./
RUN apt-get update && apt-get install -y build-essential gcc python3-dev libffi-dev libssl-dev
RUN pip install uv
# RUN pip install --no-cache-dir -r requirements.txt

# Final image
FROM base AS final

EXPOSE 5000

# CMD ["sh", "-c", "if [ '$MODE' = 'dev' ]; then python main.py; else python main.py; fi"]
