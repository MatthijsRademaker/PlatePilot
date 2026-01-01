# Production Dockerfile for PlatePilot Frontend
# Multi-stage build for optimized production image

# Build stage
FROM node:22-alpine AS builder

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm ci

# Copy source and build
COPY . .
RUN npm run build

# Production stage - serve with nginx
FROM nginx:alpine AS production

# Copy built assets
COPY --from=builder /app/dist/spa /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
