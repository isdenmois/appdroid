# Install node_modules
FROM oven/bun:latest AS modules
WORKDIR /app
COPY server/package.json .
COPY server/bun.lockb .
RUN bun install

# Build the files
FROM oven/bun:latest AS builder
WORKDIR /app
COPY --from=modules /app/node_modules ./node_modules/
COPY server .
RUN bun run build

# Run the app
FROM oven/bun:latest
WORKDIR /app
ARG PORT=3000
ENV PORT=$PORT
ENV NODE_ENV=production
EXPOSE $PORT

# Copy only what's needed for running the application
COPY server/package.json .
COPY server/tsconfig.json .
COPY --from=builder /app/node_modules/aaptjs/bin/linux ./node_modules/aaptjs/bin/linux
COPY --from=builder /app/dist dist
COPY --from=builder /app/static static
COPY --from=builder /app/drizzle drizzle

CMD ["bun", "start"]
