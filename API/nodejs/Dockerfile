FROM node:latest
WORKDIR /usr/src/nodejs
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 8080
RUN mkdir -p /elements/procs

CMD ["node", "index.js"]
