FROM node:18.18.2-bullseye-slim
WORKDIR /app
COPY . /app
RUN apt-get update && apt-get install -y bzip2 &&\
    npm install &&\ 
    npm run build-prod && \
    ls -la /app/dist
EXPOSE 80

# Serve the app using nginx
FROM nginx:1.21.0-alpine
COPY ./webserver/main.conf /etc/nginx/conf.d/default.conf
COPY --from=0 /app/dist /usr/share/nginx/html