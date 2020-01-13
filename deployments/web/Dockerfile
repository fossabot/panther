# Can't use alpine due to some node_module deps
FROM node:13.2-alpine

LABEL description="The image that builds the web code and servers it through an express NodeJS server"

# Build argument for versioning
ARG PANTHER_VERSION


# Setup default ENV variables
ENV SERVER_PORT=8080
ENV NODE_ENV=production
ENV PANTHER_VERSION=$PANTHER_VERSION

# Update debian distro
RUN apk update && apk add --no-cache --update nasm libtool alpine-sdk autoconf automake

# Move all of our logic
RUN mkdir /code

# To a different folder so we don't pollute the root of our container
WORKDIR /code

# Mount the package.json
ADD package.json .

# Install the dependencies
RUN npm --production=false install

# Mount the code
ADD web web

# Build the project with the ENV vars present in the system & run the node server
CMD npm run build-assets && npm run serve-assets

# through the port chosen during build time (defaults to 8080)
EXPOSE ${SERVER_PORT}

