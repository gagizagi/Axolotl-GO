# Build the following image from
# https://github.com/nginxinc/docker-nginx-amplify
FROM docker-nginx-amplify

COPY default.conf /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/nginx.conf
COPY /ssl-certs /etc/nginx/certs
