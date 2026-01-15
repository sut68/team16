# 🚀 Deployment Guide - Scholarship System

## 📋 Overview

```
[Dev Machine] ─── Build & Push ───> [Docker Hub] ─── Pull ───> [VM]
```

---

## 🖥️ บนเครื่อง Dev (Windows)

### ⚙️ Prerequisites
- Docker Desktop ต้องเปิดและทำงานอยู่
- Login Docker Hub แล้ว (`docker login`)

### 🔨 Build & Push ทั้งหมด (ครั้งแรก)
### gowitphooang/scholarship-frontend:latest สามารถเปลี่ยนได้ แต่อย่าลืมไปเปลี่ยนใน docker-compose.hub.yml

```bash
# Build Frontend (สำคัญ: ต้องใส่ VITE_API_URL=/api)
docker build --build-arg VITE_API_URL=/api -t gowitphooang/scholarship-frontend:latest ./frontend

# Build Backend
docker build --no-cache -t gowitphooang/scholarship-backend:latest ./backend

# Push ไป Docker Hub
docker push gowitphooang/scholarship-frontend:latest
docker push gowitphooang/scholarship-backend:latest
```

### 🔄 Update แค่ Frontend

```bash
# Build & Push Frontend
docker build --build-arg VITE_API_URL=/api -t gowitphooang/scholarship-frontend:latest ./frontend
docker push gowitphooang/scholarship-frontend:latest
```

### 🔄 Update แค่ Backend

```bash
# Build & Push Backend
docker build --no-cache -t gowitphooang/scholarship-backend:latest ./backend
docker push gowitphooang/scholarship-backend:latest
```

---

## 🌐 บน VM (Linux)

### ⚙️ Prerequisites
- Docker และ Docker Compose installed
- มีไฟล์ `docker-compose.hub.yml` และ `.env`

### 📁 ไฟล์ที่ต้องมีบน VM

```
~/team16/
├── docker-compose.hub.yml
├── .env
└── uploads/              # สำหรับเก็บไฟล์ที่ upload
```

### 📝 .env (ค่าขั้นต่ำ)

```bash
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=test
DB_HOST=postgres
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=test
JWT_SECRET=xxx
```

### 🚀 Deploy ครั้งแรก

```bash
cd ~/team16

# Pull images ทั้งหมด
sudo docker compose -f docker-compose.hub.yml pull

# Start ทุก services
sudo docker compose -f docker-compose.hub.yml up -d

# ดู logs
sudo docker compose -f docker-compose.hub.yml logs -f
```

### 🔄 Update แค่ Frontend

```bash
cd ~/team16

# Pull แค่ frontend
sudo docker compose -f docker-compose.hub.yml pull frontend

# Restart frontend
sudo docker compose -f docker-compose.hub.yml up -d frontend
```

### 🔄 Update แค่ Backend

```bash
cd ~/team16

# Pull แค่ backend
sudo docker compose -f docker-compose.hub.yml pull backend

# Restart backend
sudo docker compose -f docker-compose.hub.yml up -d backend
```

### 🔄 Update ทั้งหมด

```bash
cd ~/team16

# Pull ใหม่ทั้งหมด
sudo docker compose -f docker-compose.hub.yml pull

# Restart ทั้งหมด
sudo docker compose -f docker-compose.hub.yml up -d
```

---

## 🛠️ Useful Commands บน VM

### ดู Status

```bash
# ดู containers ที่ทำงาน
sudo docker ps

# ดู logs ทั้งหมด
sudo docker compose -f docker-compose.hub.yml logs -f

# ดู logs แค่ service เดียว
sudo docker compose -f docker-compose.hub.yml logs -f frontend
sudo docker compose -f docker-compose.hub.yml logs -f backend
sudo docker compose -f docker-compose.hub.yml logs -f postgres
```

### Stop/Start

```bash
# Stop ทั้งหมด
sudo docker compose -f docker-compose.hub.yml down

# Start ทั้งหมด
sudo docker compose -f docker-compose.hub.yml up -d

# Restart service เดียว
sudo docker compose -f docker-compose.hub.yml restart frontend
```

### Cleanup

```bash
# ลบ containers และ networks (เก็บ volumes)
sudo docker compose -f docker-compose.hub.yml down

# ลบทั้งหมด รวม volumes (⚠️ ข้อมูลหาย!)
sudo docker compose -f docker-compose.hub.yml down -v

# ลบ images เก่า
sudo docker image prune -a
```

### 🗑️ Reset Data (ล้างข้อมูลทั้งหมด)

```bash
cd ~/team16

# ⚠️ ระวัง! คำสั่งด้านล่างจะลบข้อมูลทั้งหมด

# 1. หยุด services ทั้งหมด
sudo docker compose -f docker-compose.hub.yml down

# 2. ลบ PostgreSQL volume (ล้าง database)
sudo docker volume rm team16_postgres_data

# 3. ลบ MinIO volume (ล้างไฟล์ที่ upload ทั้งหมด)
sudo docker volume rm team16_minio_data

# 4. ลบทั้ง 2 volumes พร้อมกัน
sudo docker volume rm team16_postgres_data team16_minio_data

# 5. หรือลบทุกอย่างในครั้งเดียว (containers + networks + volumes)
sudo docker compose -f docker-compose.hub.yml down -v --remove-orphans
```

### 🔄 ดู Volumes ที่มี

```bash
# ดู volumes ทั้งหมด
sudo docker volume ls

# ดูรายละเอียด volume
sudo docker volume inspect team16_postgres_data
sudo docker volume inspect team16_minio_data
```

### ♻️ Fresh Start (เริ่มใหม่ทั้งหมด)

```bash
cd ~/team16

# ลบทุกอย่าง
sudo docker compose -f docker-compose.hub.yml down -v --remove-orphans

# ลบ images เก่า
sudo docker image prune -a -f

# Pull images ใหม่
sudo docker compose -f docker-compose.hub.yml pull

# Start ใหม่ทั้งหมด
sudo docker compose -f docker-compose.hub.yml up -d

# ดู logs
sudo docker compose -f docker-compose.hub.yml logs -f
```

---

## ✅ ทดสอบหลัง Deploy

```bash
# เช็ค nginx health
curl http://localhost/health

# เช็ค API ผ่าน proxy
curl http://localhost/api/scholarship

# หรือจาก browser
http://<VM_IP>/
```

---

## 🏗️ Architecture

### Tech Stack

| Layer | Technology |
|-------|------------|
| **Frontend** | Vue 3 + TypeScript + Vite + TailwindCSS + DaisyUI |
| **Web Server** | Nginx (reverse proxy + static files) |
| **Backend** | Go (Gin framework) + GORM |
| **Database** | PostgreSQL 13 |
| **File Storage** | MinIO (S3-compatible object storage) |
| **Auth** | JWT (Cookie-based) + CSRF Token |
| **Realtime** | WebSocket |
| **Container** | Docker + Docker Compose |

### System Diagram

```
┌──────────────────────────────────────────────────────────────────────┐
│                           Internet                                   │
└─────────────────────────────────┬────────────────────────────────────┘
                                  │ Port 80
┌─────────────────────────────────▼────────────────────────────────────┐
│                     Frontend (Nginx)                                 │
│  ┌─────────────────────────────────────────────────────────────────┐│
│  │   /              → Vue.js SPA (static files)                    ││
│  │   /api/*         → proxy to backend:8080                        ││
│  │   /storage/*     → proxy to minio:9000 (file access)            ││
│  │   /health        → 200 OK                                       ││
│  │   *.js,css,img   → cached 1 year (immutable)                    ││
│  └─────────────────────────────────────────────────────────────────┘│
│  Features: Gzip, Security Headers, Large file upload (50MB)         │
└──────────────────────┬──────────────────────────┬────────────────────┘
                       │                          │
         Internal (app_network)        Internal (app_network)
                       │                          │
┌──────────────────────▼────────────┐ ┌──────────▼────────────────────┐
│       Backend (Go/Gin)            │ │       MinIO                   │
│  ┌───────────────────────────────┐│ │  ┌──────────────────────────┐ │
│  │ REST API + WebSocket          ││ │  │ S3-compatible storage    │ │
│  │ - Authentication (JWT+CSRF)   ││ │  │ - uploads bucket         │ │
│  │ - Connection pooling          ││ │  │ - news images            │ │
│  │ - GIN_MODE=release            ││ │  │ - documents              │ │
│  └───────────────────────────────┘│ │  └──────────────────────────┘ │
│  Port: 8080 (internal only)       │ │  Port: 9000 (API, internal)   │
└──────────────────────┬────────────┘ │  Port: 9001 (Console, public) │
                       │              └───────────────────────────────┘
                       │                          │
         Internal (app_network)        Internal (app_network)
                       │                          │
                       └────────────┬─────────────┘
                                    │
┌───────────────────────────────────▼──────────────────────────────────┐
│                        PostgreSQL                                    │
│  ┌─────────────────────────────────────────────────────────────────┐│
│  │  - Internal only (no external port exposed)                     ││
│  │  - Data persisted in Docker volume (postgres_data)              ││
│  │  - Health check enabled                                         ││
│  └─────────────────────────────────────────────────────────────────┘│
│  Port: 5432 (internal only)                                          │
└──────────────────────────────────────────────────────────────────────┘
```

### Data Flow

```
User Browser
    │
    ├── GET /                    → Nginx → Vue.js SPA
    ├── GET /api/scholarships    → Nginx → Backend → PostgreSQL
    ├── POST /api/upload         → Nginx → Backend → MinIO
    ├── GET /storage/uploads/*   → Nginx → MinIO (direct)
    └── WS /api/ws/*             → Nginx → Backend (WebSocket)
```

### Docker Volumes

| Volume | Purpose |
|--------|---------|
| `postgres_data` | PostgreSQL database files |
| `minio_data` | MinIO object storage files |

---

## 📅 Last Updated
2026-01-15
