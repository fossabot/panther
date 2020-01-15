#
# ****************     BUILD STAGE     *******************
#
FROM node:13.2-alpine AS build-env

LABEL description="The image that builds the web code and serves it through an express NodeJS server"

# Update debian distro
RUN apk update && apk add --no-cache --update nasm libtool alpine-sdk autoconf automake

# Move all of our logic
RUN mkdir /code

# To a different folder so we don't pollute the root of our container
WORKDIR /code

# Mount the package.json
ADD package.json package-lock.json ./

# Install the dependencies
RUN npm --production=false install

# Mount the necessary source code
ADD .babelrc tsconfig.json webpack.config.js ./
ADD enterprise/web enterprise/web
ADD web web

# Build argument for versioning (needed for the build process only)
ARG PANTHER_VERSION
ARG AWS_ACCOUNT_ID
ARG AWS_REGION
ARG GRAPHQL_ENDPOINT
ARG AWS_COGNITO_USER_POOL_ID
ARG AWS_COGNITO_APP_CLIENT_ID

# Create ENV from build ARGs
ENV PANTHER_VERSION=$PANTHER_VERSION
ENV AWS_ACCOUNT_ID=$AWS_ACCOUNT_ID
ENV AWS_REGION=$AWS_REGION
ENV GRAPHQL_ENDPOINT=$GRAPHQL_ENDPOINT
ENV AWS_COGNITO_USER_POOL_ID=$AWS_COGNITO_USER_POOL_ID
ENV AWS_COGNITO_APP_CLIENT_ID=$AWS_COGNITO_APP_CLIENT_ID

# Setup default ENV variables
ENV NODE_ENV=production

# Build the project with the ENV vars present in the system
RUN npm run build

#
# ****************     DEPLOYMENT & SERVE STAGE     *******************
#
FROM node:13.2-alpine

# Setup default ENV variables
ENV SERVER_PORT=8080
ENV NODE_ENV=production

# Make sure we add a folder to not pollute our root
RUN mkdir /code

# And make it the default WD
WORKDIR /code

# Copy the artifacts from the previous stage
COPY --from=build-env /code/web/dist dist

# Add the script that runs the express server
ADD  web/scripts/serve.js .

# Install express
RUN npm install express

# Run the express server to serve the assets
CMD node serve.js

# through the port chosen during build time (defaults to 8080)
EXPOSE ${SERVER_PORT}

