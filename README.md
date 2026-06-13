# Purgatorio

To run with file logging, use the command inside of the backend directory
```bash
LOG_FILE="logs/purg.$(date +"%F_%T").log" go run ./cmd/server/main.go
```

To run the frontend, run `npm run dev` inside of the frontend directory

### Setup

Run the setup.sh file to create the necessary keys for JWT Tokens and environment variables. Then just run `docker compose up`
