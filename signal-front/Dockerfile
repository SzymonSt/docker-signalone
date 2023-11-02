FROM node:14.17.0-alpine3.13
WORKDIR /app
COPY . /app
RUN npm install
RUN npm run build
EXPOSE 80

# Serve the app using nginx
FROM nginx:1.21.0-alpine
COPY --from=0 /app/dist/app /usr/share/nginx/html