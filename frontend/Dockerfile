# syntax=docker/dockerfile:1

FROM node:20-alpine AS build
WORKDIR /app
COPY package.json package.json
COPY tsconfig.json tsconfig.json
COPY tsconfig.node.json tsconfig.node.json
COPY vite.config.ts vite.config.ts
COPY postcss.config.cjs postcss.config.cjs
COPY tailwind.config.cjs tailwind.config.cjs
COPY .eslintrc.cjs .eslintrc.cjs
COPY .prettierrc .prettierrc
COPY index.html index.html
COPY openapi.yaml openapi.yaml
COPY src src
COPY .husky .husky
RUN npm install
RUN npm run generate || true
RUN npm run build

FROM nginx:1.25-alpine AS runtime
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
