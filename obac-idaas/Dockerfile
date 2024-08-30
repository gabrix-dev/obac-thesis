FROM node:lts-alpine3.17 as build

RUN mkdir /app
WORKDIR /app

RUN apk update && apk upgrade
RUN apk add git

COPY package.json /app/package.json
RUN yarn config set network-timeout 600000 -g
RUN yarn install

COPY .env app.vue nuxt.config.ts tailwind.config.js tsconfig.json store.js  /app/
COPY assets /app/assets
COPY components /app/components
COPY pages /app/pages
COPY plugins /app/plugins
COPY public /app/public

RUN yarn build

FROM node:lts-alpine3.17
WORKDIR /app
COPY --from=build /app/.output /app/.output

CMD node .output/server/index.mjs

