FROM node:20

WORKDIR /app

# Install dependencies first for better caching
COPY client/package*.json ./
RUN npm install

# We do NOT copy the source code here; 
# the volume in docker-compose will handle it.

EXPOSE 8000
# Force the dev server to run on port 8000
CMD ["npm", "start", "--", "--port", "8000"]