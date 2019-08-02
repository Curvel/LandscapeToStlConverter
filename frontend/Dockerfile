FROM node:alpine

# packages
RUN apk add vim

# nodejs
RUN mkdir -p /home/node/app/node_modules && chown -R node:node /home/node/app
WORKDIR /home/node/app

COPY package*.json ./

RUN npm install

USER node
COPY --chown=node:node . .

EXPOSE 1234

CMD [ "node", "app.js" ]

# keep docker running
#CMD tail -f /dev/null