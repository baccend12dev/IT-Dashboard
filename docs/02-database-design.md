# Tahap Perancangan Database

## ERD

### Entities Utama
- users
- servers
- systems
- notes
- documents
- backups

## Table Schema

### users
| Column | |
|--------|---|
| id | |
| name | |
| email | |
| password | |
| role | |
| created_at | |
| updated_at | |

### servers
| Column | |
|--------|---|
| id | |
| name | |
| ip_address | |
| os | |
| location | |
| note | |
| status | |
| created_at | |
| updated_at | |

### systems
| Column | |
|--------|---|
| id | |
| name | |
| type | |
| link | |
| server_id | |
| status | |
| description | |
| created_at | |
| updated_at | |

**type:** web, desktop

**status:** active, maintenance, deprecated

### notes
| Column | |
|--------|---|
| id | |
| system_id | |
| user_id | |
| note | |
| created_at | |
| updated_at | |

### documents
| Column | |
|--------|---|
| id | |
| system_id | |
| file_path | |
| version | |
| description | |
| uploaded_by | |
| created_at | |

### backups
| Column | |
|--------|---|
| id | |
| system_id | |
| last_backup_date | |
| method | |
| location | |
| note | |
| created_at | |