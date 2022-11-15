FROM node:18.12-alpine

WORKDIR /app

COPY . .
RUN yarn install
RUN yarn build

CMD [ "node", "dist/src/main.js"]
