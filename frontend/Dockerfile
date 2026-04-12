FROM node:20-alpine AS build

ARG VITE_API_URL=http://localhost:8080
ENV VITE_API_URL=${VITE_API_URL}

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

FROM node:20-alpine

WORKDIR /app

COPY --from=build /app/dist ./dist
RUN npm install -g serve

EXPOSE 80
CMD ["serve", "-s", "dist", "-l", "80"]