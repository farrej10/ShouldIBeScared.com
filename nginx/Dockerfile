FROM nginx

RUN rm /etc/nginx/conf.d/default.conf
#mk ssl folder for later
RUN mkdir -p /etc/ssl

COPY static/ /var/www/static
COPY nginx.conf /etc/nginx/conf.d
COPY nginx.conf /etc/nginx/conf.d