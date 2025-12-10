
ARG MODE

FROM python:3.11-slim AS base

WORKDIR /app

# COPY plugins/ ./

# RUN pip install --no-cache-dir -r requirements.txt

# Final image
FROM base AS final

EXPOSE 5000

# CMD ["sh", "-c", "if [ '$MODE' = 'dev' ]; then python main.py; else python main.py; fi"]
