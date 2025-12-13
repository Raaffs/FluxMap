FROM node:20 AS builder

WORKDIR /app

# copy only package files first for caching
COPY client/package*.json ./
RUN npm install

# now copy the rest WITHOUT node_modules
COPY client/ .

RUN npm run build

FROM nginx:alpine

COPY docker/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/build /usr/share/nginx/html

EXPOSE 8000
