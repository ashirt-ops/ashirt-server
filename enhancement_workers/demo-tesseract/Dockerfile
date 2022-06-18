FROM node:16-buster-slim

WORKDIR /app

RUN apt-get update && apt-get upgrade -y \
    && apt-get install -y tesseract-ocr

COPY . .
RUN yarn install
RUN yarn build

CMD [ "node", "dist/src/main.js"]
