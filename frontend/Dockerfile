# develop stage
FROM node:20.11-alpine as build-stage
ARG ENVIRONMENT
ENV ENVIRONMENT=$ENVIRONMENT

WORKDIR /app
COPY . .
RUN yarn
RUN yarn build

# production stage
# FROM node:20.11-alpine
# WORKDIR /app
# COPY --from=build-stage /app/dist .
# RUN yarn && yarn cache clean --all

# EXPOSE 80
# CMD ["node", "index.js"]
FROM nginx:1.25-alpine as production-stage
COPY --from=build-stage /app/nginx_conf/nginx.conf /etc/nginx/nginx.conf
COPY --from=build-stage /app/dist/ /usr/share/nginx/html
RUN sed 's/root   \/usr\/share\/nginx\/html;/root   \/usr\/share\/nginx\/html;\n        try_files $uri $uri\/  \/index.html;/g' /etc/nginx/conf.d/default.conf -i
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
