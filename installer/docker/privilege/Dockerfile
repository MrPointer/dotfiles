# Installer privilege test image: cache dependencies for repeated local runs.
# Tests and test users are supplied at runtime by the shared Linux runner.
FROM ubuntu:24.04

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends sudo expect passwd util-linux \
    && rm -rf /var/lib/apt/lists/*
