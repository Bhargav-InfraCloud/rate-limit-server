# Rate Limit Server :no_entry_sign:
Not a package to import. Not something you can integrate with your program. Just a simple Go HTTP server CLI command binary/container that can be used as rate limited calls for testing.

## How to :question:
1. Use binary:
    ```
    make run
    ```
    Or build and run directly:
    ```
    go build -o bin/rate-limit-server cmd/rate-limit-server/main.go
    ./bin/rate-limit-server
    ```
2. Use docker container:
    ```
    docker run -d --rm --name rate-limit-server -p 8080:8080 bhargav0infracloudio/rate-limit-server:latest
    ```
3. Build docker container and run:
    ```
    make docker-run
    ```
    Or build and run directly:
    ```
    docker build -t=bhargav0infracloudio/rate-limit-server:latest .
    docker run --rm --name rate-limit-server -p 8080:8080 bhargav0infracloudio/rate-limit-server:latest
    ```

---

## Examples :closed_book:
### 1. Simple Rate Limited Run
```
▶ curl localhost:8080/id/vague -H "count: 2"
{"status":"OK"}
▶ curl localhost:8080/id/vague              
{"status":"OK"}
▶ curl localhost:8080/id/vague
{"message":"rate limit reached for the specific ID","service_code":"1001","logs":["failed to add id \"vague\""]}
```

### 2. Direct Fail
```
▶ curl localhost:8080/id/vague -H "count: 0"
{"message":"rate limit reached for the specific ID","service_code":"1001","logs":["failed to add id \"vague\""]}
```

### 3. Reset Count
```
▶ curl localhost:8080/id/vague -H "count: 1"
{"status":"OK"}
▶ curl localhost:8080/id/vague              
{"message":"rate limit reached for the specific ID","service_code":"1001","logs":["failed to add id \"vague\""]}
▶ curl localhost:8080/id/vague -H "count: 2" -H "reset: true"
{"status":"OK"}
▶ curl localhost:8080/id/vague                               
{"status":"OK"}
▶ curl localhost:8080/id/vague
{"message":"rate limit reached for the specific ID","service_code":"1001","logs":["failed to add id \"vague\""]}
```

---
