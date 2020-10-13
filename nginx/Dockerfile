FROM nginx:stable

LABEL maintainer="zuoguocai@126.com"  version="1.0" description="devops2020-demo"

WORKDIR /web
RUN  mkdir /web/daohang;mkdir /web/caidan;mkdir /web/live2d
ADD daohang ./daohang
ADD caidan     ./caidan
ADD live2d     ./live2d

COPY nginx.conf /etc/nginx/
RUN chown -R nginx:nginx /web
EXPOSE 11443
CMD ["nginx", "-g", "daemon off;"]
