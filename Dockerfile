FROM --platform=$BUILDPLATFORM node:22.5-alpine AS build
RUN apk add --no-cache libc6-compat
RUN apk update

RUN npm install turbo --global

WORKDIR /opt/app
ENV PATH=/opt/app/node_modules/.bin:$PATH

COPY . .
RUN cd ./hanko/frontend && npm install && npm run build:elements
RUN npm install && npm run build

CMD ["npm", "run", "dev"]