# Sample Log Files for IP Parsing Interview Practice

## 1. Apache/NGINX Access Log
```
192.168.1.100 - - [30/Sep/2025:10:23:45 +0000] "GET /api/users HTTP/1.1" 200 1234 "-" "Mozilla/5.0"
10.0.0.15 - - [30/Sep/2025:10:23:46 +0000] "POST /api/login HTTP/1.1" 401 512 "-" "curl/7.68.0"
192.168.1.100 - - [30/Sep/2025:10:23:47 +0000] "GET /api/data HTTP/1.1" 200 5678 "-" "Mozilla/5.0"
172.16.0.5 - admin [30/Sep/2025:10:23:48 +0000] "GET /admin/dashboard HTTP/1.1" 200 9876 "-" "Chrome/91.0"
10.0.0.15 - - [30/Sep/2025:10:23:49 +0000] "POST /api/login HTTP/1.1" 200 256 "-" "curl/7.68.0"
192.168.1.100 - - [30/Sep/2025:10:23:50 +0000] "GET /api/profile HTTP/1.1" 200 2048 "-" "Mozilla/5.0"
203.0.113.42 - - [30/Sep/2025:10:23:51 +0000] "GET / HTTP/1.1" 200 15234 "-" "GoogleBot/2.1"
10.0.0.15 - - [30/Sep/2025:10:23:52 +0000] "GET /api/status HTTP/1.1" 200 128 "-" "curl/7.68.0"
```

## 2. MongoDB Connection Log
```
2025-09-30T10:15:23.456+0000 I NETWORK  [listener] connection accepted from 10.0.1.50:54321 #12345 (10 connections now open)
2025-09-30T10:15:24.123+0000 I NETWORK  [conn12345] received client metadata from 10.0.1.50:54321 conn12345: { application: { name: "MongoDB Shell" }, driver: { name: "MongoDB Internal Client", version: "5.0.0" } }
2025-09-30T10:15:25.789+0000 I NETWORK  [listener] connection accepted from 192.168.100.20:43210 #12346 (11 connections now open)
2025-09-30T10:15:26.234+0000 I ACCESS   [conn12346] Successfully authenticated as principal admin on admin from client 192.168.100.20:43210
2025-09-30T10:15:27.567+0000 I NETWORK  [listener] connection accepted from 10.0.1.50:54322 #12347 (12 connections now open)
2025-09-30T10:15:28.890+0000 I NETWORK  [conn12347] received client metadata from 10.0.1.50:54322 conn12347: { application: { name: "PyMongo" }, driver: { name: "PyMongo", version: "4.0.1" } }
2025-09-30T10:15:29.123+0000 I NETWORK  [listener] connection accepted from 172.31.5.100:33445 #12348 (13 connections now open)
2025-09-30T10:15:30.456+0000 I NETWORK  [conn12345] end connection 10.0.1.50:54321 (12 connections now open)
2025-09-30T10:15:31.789+0000 I NETWORK  [listener] connection accepted from 192.168.100.20:43211 #12349 (13 connections now open)
```

## 3. Syslog Format
```
Sep 30 10:30:15 web-server sshd[12345]: Accepted password for admin from 203.0.113.10 port 52123 ssh2
Sep 30 10:30:16 web-server sshd[12346]: Failed password for root from 198.51.100.50 port 44556 ssh2
Sep 30 10:30:17 web-server kernel: [UFW BLOCK] IN=eth0 OUT= MAC=00:11:22:33:44:55 SRC=198.51.100.50 DST=10.0.0.1 LEN=60 TOS=0x00
Sep 30 10:30:18 web-server sshd[12347]: Accepted publickey for deploy from 203.0.113.10 port 52124 ssh2
Sep 30 10:30:19 web-server postfix[12348]: connect from mail-server.example.com[192.0.2.25]
Sep 30 10:30:20 web-server sshd[12349]: Failed password for admin from 198.51.100.50 port 44557 ssh2
Sep 30 10:30:21 web-server sshd[12350]: Accepted password for admin from 203.0.113.10 port 52125 ssh2
Sep 30 10:30:22 web-server kernel: [UFW BLOCK] IN=eth0 OUT= MAC=00:11:22:33:44:55 SRC=198.51.100.75 DST=10.0.0.1 LEN=60 TOS=0x00
```

## 4. AWS VPC Flow Logs
```
2 123456789012 eni-1a2b3c4d 172.31.16.5 192.0.2.100 80 49152 6 10 5000 1632988800 1632988860 ACCEPT OK
2 123456789012 eni-1a2b3c4d 192.0.2.100 172.31.16.5 49152 80 6 8 4000 1632988800 1632988860 ACCEPT OK
2 123456789012 eni-1a2b3c4d 10.0.1.25 198.51.100.200 443 54321 6 15 7500 1632988820 1632988880 ACCEPT OK
2 123456789012 eni-1a2b3c4d 198.51.100.200 10.0.1.25 54321 443 6 12 6000 1632988820 1632988880 ACCEPT OK
2 123456789012 eni-1a2b3c4d 172.31.16.5 203.0.113.75 22 55123 6 5 2500 1632988840 1632988900 REJECT OK
2 123456789012 eni-1a2b3c4d 10.0.1.25 192.0.2.100 443 49200 6 20 10000 1632988860 1632988920 ACCEPT OK
```

## 5. JSON Application Log
```json
{"timestamp":"2025-09-30T10:45:12.123Z","level":"INFO","service":"api-gateway","client_ip":"10.5.5.50","user_id":"user123","endpoint":"/api/v1/orders","method":"GET","status":200,"response_time_ms":45}
{"timestamp":"2025-09-30T10:45:13.456Z","level":"WARN","service":"api-gateway","client_ip":"192.168.50.100","user_id":"user456","endpoint":"/api/v1/users","method":"POST","status":429,"response_time_ms":12,"error":"rate_limit_exceeded"}
{"timestamp":"2025-09-30T10:45:14.789Z","level":"INFO","service":"api-gateway","client_ip":"10.5.5.50","user_id":"user123","endpoint":"/api/v1/products","method":"GET","status":200,"response_time_ms":67}
{"timestamp":"2025-09-30T10:45:15.012Z","level":"ERROR","service":"api-gateway","client_ip":"172.20.10.75","user_id":"user789","endpoint":"/api/v1/checkout","method":"POST","status":500,"response_time_ms":5000,"error":"database_timeout"}
{"timestamp":"2025-09-30T10:45:16.345Z","level":"INFO","service":"api-gateway","client_ip":"192.168.50.100","user_id":"user456","endpoint":"/api/v1/users","method":"POST","status":201,"response_time_ms":89}
{"timestamp":"2025-09-30T10:45:17.678Z","level":"INFO","service":"api-gateway","client_ip":"10.5.5.50","user_id":"user123","endpoint":"/api/v1/cart","method":"PUT","status":200,"response_time_ms":34}
```

## 6. HAProxy Load Balancer Log
```
Sep 30 11:00:01 haproxy[1234]: 203.0.113.25:54321 [30/Sep/2025:11:00:01.123] frontend-http backend-app/server1 0/0/1/25/26 200 1234 - - ---- 1/1/0/0/0 0/0 "GET /health HTTP/1.1"
Sep 30 11:00:02 haproxy[1234]: 198.51.100.100:44556 [30/Sep/2025:11:00:02.456] frontend-https backend-api/server2 0/0/2/150/152 201 5678 - - ---- 2/2/0/0/0 0/0 "POST /api/orders HTTP/1.1"
Sep 30 11:00:03 haproxy[1234]: 203.0.113.25:54322 [30/Sep/2025:11:00:03.789] frontend-http backend-app/server1 0/0/1/15/16 200 2048 - - ---- 1/1/0/0/0 0/0 "GET /status HTTP/1.1"
Sep 30 11:00:04 haproxy[1234]: 192.0.2.50:33445 [30/Sep/2025:11:00:04.012] frontend-https backend-api/server3 0/0/5/500/505 503 0 - - sHNN 3/3/0/0/0 0/0 "POST /api/payment HTTP/1.1"
Sep 30 11:00:05 haproxy[1234]: 198.51.100.100:44557 [30/Sep/2025:11:00:05.345] frontend-https backend-api/server2 0/0/1/75/76 200 9876 - - ---- 2/2/0/0/0 0/0 "GET /api/products HTTP/1.1"
```

---

## Interview Tips:

**Key things to demonstrate:**

1. **Edge Cases to Consider:**
   - Invalid IP formats (300.400.500.600)
   - IPv6 addresses (if bonus points)
   - Private vs public IPs
   - IPs in different positions (source vs destination)
   - Malformed log lines
   - Empty or null values

2. **Efficiency Considerations:**
   - How would you handle GB/TB-sized log files?
   - Streaming vs loading everything into memory
   - Using appropriate data structures (Counter, HashMap)
   - Regex performance vs string splitting

3. **Validation & Testing:**
   - Test with sample data from each format
   - Verify counts are accurate
   - Handle different log formats gracefully
   - Consider what happens with mixed log formats

4. **Discussion Points:**
   - Why might one IP appear more than others? (bot traffic, legitimate heavy user, DDoS)
   - How would you extend this to find suspicious patterns?
   - What if you needed real-time analysis?
   - How would you parallelize this for distributed systems?