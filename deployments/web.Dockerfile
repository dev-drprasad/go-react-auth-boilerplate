# build environment
# make sure to match with .nvmrc version
FROM node:14.14-alpine as build
WORKDIR /app

ENV NODE_ENV=production

COPY web/package.json /app/package.json
COPY web/package-lock.json /app/package-lock.json
RUN npm install


ARG REACT_APP_API_BASE_URL

COPY web/public /app/public
COPY web/jsconfig.json /app
COPY web/craco.config.js /app
COPY web/src /app/src

RUN npm run build

# production environment
FROM nginx:1.19-alpine

EXPOSE 80

COPY --from=build /app/build /usr/share/nginx/html
# RUN rm /etc/nginx/conf.d/default.conf
COPY deployments/nginx.conf.template /etc/nginx/nginx.conf.template
COPY deployments/web.docker-entrypoint.sh /

ENTRYPOINT ["sh", "/web.docker-entrypoint.sh"]

CMD ["nginx", "-g", "daemon off;"]

# https://mherman.org/blog/dockerizing-a-react-app/
