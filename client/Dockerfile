FROM node:lts-alpine
WORKDIR /app
COPY . .
RUN npm i
RUN npm run build
CMD ["npx", "serve", "-s", "dist/outshort/browser", "-l", "4200"]
EXPOSE 4200
