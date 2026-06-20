# Tahap membuat API

## API Structure

```
/api
/api/systems
/api/servers
/api/notes
/api/documents
/api/backups
/api/users
```

## Example API

### Get Systems

```http
GET /api/systems
```

**Response:**

```json
[
    {
        "id": 1,
        "name": "QA Qualification",
        "type": "web",
        "status": "active"
    }
]
```

### Create System

```http
POST /api/systems
```

**Body:**

```json
{
    "name": "QA System",
    "type": "web",
    "server_id": 1
}
```

### Get System Detail

```http
GET /api/systems/{id}
```

**Response:**

```json
{
    "system": {},
    "notes": [],
    "documents": []
}
```