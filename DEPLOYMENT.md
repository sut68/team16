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

```
┌─────────────────────────────────────────────────────────┐
│                   Internet                              │
└────────────────────────┬────────────────────────────────┘
                         │ Port 80
┌────────────────────────▼────────────────────────────────┐
│              Frontend (Nginx)                           │
│  ┌────────────────────────────────────────────────────┐│
│  │  /          → Vue.js SPA                           ││
│  │  /api/*     → proxy to backend:8080                ││
│  │  /health    → 200 OK                               ││
│  └────────────────────────────────────────────────────┘│
└────────────────────────┬────────────────────────────────┘
                         │ Internal (app_network)
┌────────────────────────▼────────────────────────────────┐
│              Backend (Go/Gin)                           │
│  - Connection pooling enabled                          │
│  - GIN_MODE=release                                    │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│              PostgreSQL                                 │
│  - Internal only (no external port)                    │
│  - Data persisted in volume                            │
└─────────────────────────────────────────────────────────┘
```

---

## 📅 Last Updated
2026-01-10
