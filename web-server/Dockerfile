# -----
# nginx
# -----

FROM nginx:stable-alpine3.23-perl

WORKDIR /app/frontend

RUN chown nginx:nginx /var/cache/nginx/

# main nginx conf
COPY ./web-server/nginx.conf /etc/nginx/nginx.conf
# included nginx confs
COPY ./web-server/http_site.conf /etc/nginx/conf.d/http_site.conf
COPY ./web-server/https_site.conf /etc/nginx/conf.d/https_site.conf

# copy built frontend files
COPY ./web/static ./static
