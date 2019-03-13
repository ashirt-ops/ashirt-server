FROM nginx:stable

COPY frontend /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

RUN sed -i 's/\/var\/run\/nginx.pid/\/var\/cache\/nginx\/nginx.pid/g' /etc/nginx/nginx.conf

CMD nginx -g "daemon off;"