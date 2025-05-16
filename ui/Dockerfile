# ---------- DEVELOPMENT STAGE ----------
FROM node:20-alpine AS dev
WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

EXPOSE 5173
CMD ["npm", "run", "dev", "--", "--host"]

# ---------- BUILD STAGE ----------
FROM node:20-alpine AS builder
WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# ---------- PRODUCTION STAGE ----------
FROM nginx:stable-alpine AS prod
COPY --from=builder /app/dist /usr/share/nginx/html

# Optional: configure nginx for SPA routing
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

