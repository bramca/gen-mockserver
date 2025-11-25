
# syntax=docker/dockerfile:1

ARG NODE_VERSION=21.5.0

FROM node:${NODE_VERSION}-alpine

ENV NODE_ENV production

WORKDIR /app

# Copy the rest of the source files into the image.
COPY . .
RUN chmod 777 db.json

RUN npm install

# Run the application as a non-root user.
USER node

# Expose the port that the application listens on.
EXPOSE 5000

# Run the application.
CMD npm start
