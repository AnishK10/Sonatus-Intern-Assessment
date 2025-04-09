#  Distributed Log Aggregator

This project is a lightweight, concurrent log aggregation service built using Go. It allows ingestion of logs from various microservices and supports querying logs based on service name and time ranges. Logs are held in memory and expire after a period (default: 1 hour).

##  Features

- ✅ Log ingestion via REST API (`POST /logs`)
- ✅ Log retrieval by service and time window (`GET /logs`)
- ✅ Thread-safe with in-memory storage using `sync.RWMutex`
- ✅ Automatic expiration of old logs (every 1 minute)
- ✅ Structured logging using [Zap](https://github.com/uber-go/zap)
- ✅ Clean code architecture (`cmd/`, `internal/`, `test/` structure)

## Project Structure
```text
.
├── cmd/logaggregator         # Main entry point
├── internal
│   ├── handler               
│   ├── logger                
│   ├── model                 
│   └── store                 
├── test                      
├── go.mod / go.sum           
└── README.md                
```

## 🏁 Getting Started
### 1. Clone the repository
```bash
git clone https://github.com/your-username/log-aggregator.git
cd log-aggregator
```

### 2. Set up dependencies

```bash
go mod tidy
```

### 3. Run the service

```bash
go run cmd/logaggregator/main.go
```
The server will start at:
http://localhost:8080

---

###  Test the endpoints using Postman

You can use Postman or any API client to test the service:

- **POST** `http://localhost:8080/logs` – to ingest a log
- **GET** `http://localhost:8080/logs?service=SERVICE_NAME&start=START_TIME&end=END_TIME` – to query logs


###  Sample Request & Response

####  Ingest a Log (POST)

**Endpoint:**
POST http://localhost:8080/logs


**Body:**
```json
{
  "service_name": "auth-service",
  "timestamp": "2025-04-08T23:20:00Z",
  "message": "User login successful"
}
```
**Response:**
```json
{
  "status": "success",
  "message": "Log added successfully"
}
```

####  Query Logs (GET)

**Endpoint:**
GET `/logs?service=<service_name>&start=<timestamp>&end=<timestamp>`


**Response:**
```json
[
  {
    "timestamp": "2025-04-08T23:20:00Z",
    "message": "User login successful"
  }
]

```
