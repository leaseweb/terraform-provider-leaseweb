FROM golang:1.25.5-bookworm AS base

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN wget -q -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg \
    && echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com bookworm main" | tee /etc/apt/sources.list.d/hashicorp.list \
    && apt-get update && apt-get install -y --no-install-recommends \
    terraform=1.8.3-1 \
    && rm -rf /var/lib/apt/lists/*

RUN git config --global --add safe.directory /app
