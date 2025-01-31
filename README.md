# ClearingHouse

## Getting Started

### 1. Start Development  
Ensure Docker Engine is running, then start the application:  
```bash
docker compose up --build -d
```

## Updating Swagger Documentation  

To regenerate Swagger documentation, follow these steps:  

1. Install `swag` if you haven't already:  
```bash
   go install github.com/swaggo/swag/cmd/swag@latest
```  

2. Generate the Swagger docs:  
```bash
   swag init -g cmd/api/main.go -o docs --parseInternal  
```

📌 **Run the command from the project's root directory.**
