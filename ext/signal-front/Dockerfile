FROM node:18.18.2-bullseye-slim
WORKDIR /app
COPY . /app
RUN npm install
RUN npm run build && \
    ls -la /app/dist
EXPOSE 80

# Serve the app using nginx
FROM nginx:1.21.0-alpine
COPY --from=0 /app/dist/signal-one /usr/share/nginx/html