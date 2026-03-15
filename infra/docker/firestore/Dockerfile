FROM eclipse-temurin:21-jre-jammy

RUN apt-get update \
  && apt-get install -y --no-install-recommends curl ca-certificates gnupg \
  && rm -rf /var/lib/apt/lists/*

RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
  && apt-get update \
  && apt-get install -y --no-install-recommends nodejs \
  && node --version \
  && npm --version \
  && npm install -g firebase-tools \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace
