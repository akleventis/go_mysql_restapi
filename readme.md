## 🛸 GO REST API 🛸
---

Features
- MySql database integration
- RESTful architecture 
- Test files (this was a doozy)
- Hit all endpoints with Postman (correct HTTP codes)

### Run locally
- Navigate to this struct @ src/db/dbconnect.go
```
const (
	username = "gosql"
	password = "gosql"
	hostname = "127.0.0.1:3306"
	dbname   = "pets"
)
```
- Either create a new MySql user (grant all permissions), or throw in your username/password
- Verify MySql is running on local machine
- The database and tables will be created on first run

### Runs on localhost:8000
---
Endpoints:

- `GET /dogs` | `GET /dogs/{id}` | `POST /dogs` | `PATCH /dogs/{id}` | `DELETE /dogs/{id}`
- `GET /cats` | `GET /cats/{id}` | `POST /cats` | `PATCH /cats/{id}` | `DELETE /cats/{id}`
