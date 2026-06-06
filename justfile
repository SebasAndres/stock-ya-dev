default: dev

# Start the StockYa backend (serves frontend on :8080)
dev:
    cd backend && go run .
