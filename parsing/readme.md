# SRE Practice Dataset

## 1. Web Server Access Logs (access.log)

```
192.168.1.100 - - [25/Sep/2025:10:15:23 +0000] "GET /api/users HTTP/1.1" 200 1534 "-" "Mozilla/5.0"
10.0.0.5 - - [25/Sep/2025:10:15:24 +0000] "POST /api/login HTTP/1.1" 200 89 "-" "curl/7.68.0"
192.168.1.101 - - [25/Sep/2025:10:15:25 +0000] "GET /api/orders HTTP/1.1" 500 0 "-" "Mozilla/5.0"
malformed log entry without proper format
172.16.0.10 - admin [25/Sep/2025:10:15:26 +0000] "DELETE /api/users/123 HTTP/1.1" 404 142 "-" "PostmanRuntime/7.29.0"
192.168.1.100 - - [25/Sep/2025:10:15:27 +0000] "GET /health HTTP/1.1" 200 15 "-" "kube-probe/1.21"
10.0.0.5 - - [25/Sep/2025:10:16:01 +0000] "POST /api/orders HTTP/1.1" 201 456 "-" "mobile-app/2.1.0"
192.168.1.102 - - [25/Sep/2025:10:16:02 +0000] "GET /api/users HTTP/1.1" 429 78 "-" "bot-scanner/1.0"
172.16.0.11 - - [25/Sep/2025:10:16:03 +0000] "GET /api/metrics HTTP/1.1" 200 2048 "-" "prometheus/2.30.0"
192.168.1.100 - - [25/Sep/2025:10:16:04 +0000] "POST /api/logout HTTP/1.1" 500 0 "-" "Mozilla/5.0"
192.168.1.103 - - [25/Sep/2025:10:16:05 +0000] "GET /api/orders HTTP/1.1" 200 890 "-" "Mozilla/5.0"
```

## 2. Application Error Logs (app.log)

```
2025-09-25 10:15:23 INFO [user-service] User login successful for user_id=12345
2025-09-25 10:15:24 ERROR [order-service] Database connection timeout
	at com.example.OrderService.getConnection(OrderService.java:45)
	at com.example.OrderService.createOrder(OrderService.java:89)
	at com.example.OrderController.handlePost(OrderController.java:23)
2025-09-25 10:15:25 WARN [user-service] Rate limit exceeded for IP: 192.168.1.102
2025-09-25 10:15:26 ERROR [payment-service] Payment gateway returned error: insufficient_funds
2025-09-25 10:15:27 INFO [health-check] All services healthy
Invalid log format here
2025-09-25 10:16:01 INFO [order-service] Order created successfully order_id=67890
2025-09-25 10:16:02 ERROR [user-service] Authentication failed for token: expired
2025-09-25 10:16:03 ERROR [order-service] Database connection timeout
	at com.example.OrderService.getConnection(OrderService.java:45)
	at com.example.OrderService.updateOrder(OrderService.java:112)
2025-09-25 10:16:04 ERROR [auth-service] JWT token validation failed
	org.springframework.security.jwt.InvalidTokenException: Token expired
	at org.springframework.security.jwt.JwtHelper.decode(JwtHelper.java:89)
2025-09-25 10:16:05 INFO [order-service] Order retrieved successfully order_id=67891
```

## 3. System Metrics (metrics.csv)

```
timestamp,service,cpu_percent,memory_mb,disk_io_ops,network_bytes_in,network_bytes_out,response_time_ms
2025-09-25T10:15:00Z,user-service,45.2,512,150,1024,2048,120
2025-09-25T10:16:00Z,user-service,67.8,589,180,1536,3072,145
2025-09-25T10:17:00Z,user-service,23.1,445,90,768,1024,95
2025-09-25T10:15:00Z,order-service,78.9,1024,250,2048,4096,350
2025-09-25T10:16:00Z,order-service,89.5,1156,380,3072,5120,450
malformed,data,here,missing,fields
2025-09-25T10:17:00Z,order-service,34.2,892,120,1024,2048,180
2025-09-25T10:15:00Z,payment-service,12.3,256,45,512,1024,200
2025-09-25T10:16:00Z,payment-service,15.7,289,67,768,1536,250
2025-09-25T10:17:00Z,payment-service,9.8,234,34,256,512,180
2025-09-25T10:15:00Z,auth-service,56.4,678,89,1536,2048,90
2025-09-25T10:16:00Z,auth-service,72.1,745,145,2048,3072,110
2025-09-25T10:17:00Z,auth-service,41.8,623,76,1024,1536,85
```

## 4. Configuration Files

### config1.yaml

```yaml
database:
  host: localhost
  port: 5432
  name: production_db

services:
  user-service:
    replicas: 3
    memory: 512Mi
    cpu: 0.5
  order-service:
    replicas: 2
    memory: 1Gi
    cpu: 1.0

monitoring:
  enabled: true
  interval: 30s
```

### config2.yaml

```yaml
database:
  host: db.example.com # Override
  ssl: true # New setting

services:
  user-service:
    replicas: 5 # Override
  payment-service: # New service
    replicas: 2
    memory: 256Mi
    cpu: 0.25

alerts:
  cpu_threshold: 80
  memory_threshold: 90
```

## 5. Alert Events (alerts.json)

```json
[
  {
    "timestamp": "2025-09-25T10:15:30Z",
    "service": "order-service",
    "type": "high_cpu",
    "value": 89.5,
    "threshold": 80,
    "severity": "warning"
  },
  {
    "timestamp": "2025-09-25T10:15:35Z",
    "service": "order-service",
    "type": "database_error",
    "message": "Connection timeout",
    "severity": "critical"
  },
  {
    "timestamp": "2025-09-25T10:15:40Z",
    "service": "user-service",
    "type": "rate_limit",
    "value": 150,
    "threshold": 100,
    "severity": "warning"
  },
  {
    "timestamp": "2025-09-25T10:16:15Z",
    "service": "payment-service",
    "type": "external_api_error",
    "message": "Gateway timeout",
    "severity": "error"
  },
  {
    "timestamp": "2025-09-25T10:16:20Z",
    "service": "order-service",
    "type": "high_memory",
    "value": 92.3,
    "threshold": 90,
    "severity": "critical"
  },
  {
    "timestamp": "2025-09-25T10:16:25Z",
    "service": "auth-service",
    "type": "token_validation_error",
    "message": "High failure rate",
    "severity": "error"
  }
]
```

# Practice Questions Using This Dataset:

1. **Log Parsing**: Parse access.log and calculate error rate per minute, identify top error-generating IPs
2. **Error Analysis**: Parse app.log to find all stack traces and group by service and error type
3. **Metrics Processing**: Calculate rolling averages from metrics.csv and identify services exceeding thresholds
4. **Config Merging**: Merge config1.yaml and config2.yaml with override rules
5. **Alert Correlation**: Find related alerts that occurred within 5-minute windows
6. **Anomaly Detection**: Identify unusual patterns in the metrics data
7. **Data Cleaning**: Handle malformed entries in each dataset
8. **Report Generation**: Create a service health summary from all data sources
