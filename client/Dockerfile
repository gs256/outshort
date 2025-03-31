FROM node:lts-alpine AS build
WORKDIR /app
COPY . .
RUN npm i
RUN npm run build

FROM node:lts-alpine
WORKDIR /app
COPY --from=build /app/dist/outshort/browser ./dist
RUN npm install -g serve
CMD ["serve", "-s", "dist", "-l", "4200"]
EXPOSE 4200
