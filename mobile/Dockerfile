FROM node:20-bullseye

ARG NODE_ENV=development
ENV NODE_ENV=${NODE_ENV}

ARG PORT=19006
ENV PORT $PORT
EXPOSE $PORT 19001 19002

ENV NPM_CONFIG_PREFIX=/home/node/.npm-global
ENV PATH /home/node/.npm-global/bin:$PATH
RUN npm i --unsafe-perm --allow-root -g npm@latest expo-cli@latest

RUN mkdir /opt/mobile
WORKDIR /opt/mobile
ENV PATH /opt/mobile/.bin:$PATH
COPY package.json package-lock.json ./
RUN npm install

WORKDIR /opt/mobile/app

COPY . .

ENTRYPOINT ["npm", "run"]
CMD ["web"]