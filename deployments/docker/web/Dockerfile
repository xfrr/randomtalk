FROM node:20-alpine

WORKDIR /app

# Copy only package.json and lockfiles for caching dependencies
COPY package.json yarn.lock* package-lock.json* pnpm-lock.yaml* ./

# Install dependencies (includes expo as a devDependency)
RUN \
    if [ -f yarn.lock ]; then yarn install --frozen-lockfile; \
    elif [ -f pnpm-lock.yaml ]; then npm install -g pnpm && pnpm install --frozen-lockfile; \
    else npm install --omit=dev; fi

# Copy app source
COPY . .

# Expose Expo ports
EXPOSE 8081 19000 19001 19002

# For CI/non-interactive environments
ENV CI=true

# Start Expo with NPX (no --no-interactive)
CMD ["npx", "expo", "start", "--web"]