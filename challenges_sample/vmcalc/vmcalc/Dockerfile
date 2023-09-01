FROM --platform=linux/amd64 node:16.13.1-bullseye-slim

RUN mkdir -p /app

WORKDIR /app

COPY package.json .

RUN npm install

COPY . .

USER node

EXPOSE 3000

CMD ["node", "index.js"]