#
# Initial build stage
#
FROM node:13.2-alpine AS build-env

LABEL description="The image that builds the web code and serves it through an express NodeJS server"

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

# Build the project with the ENV vars present in the system
RUN npm run build-assets

#
# Final deployment/serving stage
#
FROM node:13.2-alpine

RUN mkdir /code

WORKDIR /code

# Copy the artifact from the previous stage
COPY --from=build-env /code/dist dist

# Install express
RUN npm install express

# Run the express server to serve the assets
CMD npm run serve-assets

# through the port chosen during build time (defaults to 8080)
EXPOSE ${SERVER_PORT}

