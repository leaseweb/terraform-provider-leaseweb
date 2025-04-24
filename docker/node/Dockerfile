FROM node:jod-slim@sha256:157c7ea6f8c30b630d6f0d892c4f961eab9f878e88f43dd1c00514f95ceded8a

RUN corepack enable \
    && corepack prepare pnpm@latest-9 --activate \
    && pnpm config set store-dir .pnpm-store
