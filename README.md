# CVWO Forum

[![Ask DeepWiki](https://devin.ai/assets/askdeepwiki.png)](https://deepwiki.com/Kk120306/cvwo-2026)

A full-stack forum application built with a Go backend and React frontend. Users can engage in discussions, create posts under various topics, comment, and vote. This is a project for the CVWO 2026 Application.

**🚀 Live Demo:** [http://3.25.177.70/](#)

---

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Performance](#performance)
- [Getting Started](#getting-started)
- [Deployment](#deployment)

---

## Features

### Core Functionality

- **User Authentication** - Secure signup and login with JWT-based sessions stored in HTTP-only cookies - Username only
- **Post Management** - Full CRUD operations for posts with rich text editing
- **Voting System** - Upvote/downvote posts and comments
- **User Profiles** - View post and comment history with user statistics
- **Client-Side Filtering** - Real-time search and sort for posts and comments

### Rich Content

- **TipTap Editor** - Rich text editing with formatting, links, and embedded media
- **Image Uploads** - Image handling via pre-signed S3 URLs
- **CDN Distribution** - Fast global image delivery through CloudFront

### Admin Features

- **Pin Posts** - Highlight important posts at the top of topic feeds
- **Topic Management** - Create and organize discussion categories
- **Moderation Tools** - Delete inappropriate content

---

## 🛠️ Tech Stack

### Backend

- **[Go](https://golang.org/)**
- **[Gin](https://gin-gonic.com/)**
- **[GORM](https://gorm.io/)**
- **[PostgreSQL](https://www.postgresql.org/)** - Relational database (via [Neon](https://neon.tech/))

### Frontend

- **[React](https://reactjs.org/)**
- **[TypeScript](https://www.typescriptlang.org/)**
- **[Vite](https://vitejs.dev/)**
- **[Redux Toolkit](https://redux-toolkit.js.org/)**
- **[Material-UI](https://mui.com/)**

### Infrastructure & DevOps

- **[Docker](https://www.docker.com/)**
- **[Docker Compose](https://docs.docker.com/compose/)**
- **[Nginx](https://www.nginx.com/)** - Reverse proxy and static file serving
- **[AWS EC2](https://aws.amazon.com/ec2/)** - Cloud hosting
- **[AWS S3](https://aws.amazon.com/s3/)** - Object storage for images
- **[AWS CloudFront](https://aws.amazon.com/cloudfront/)** - Global CDN
- **[GitHub Actions](https://github.com/features/actions)**

---

## Architecture

This monorepo contains two main services:

### **Backend Service** (`backend/`)

- RESTful API built with Go and Gin framework
- GORM for PostgreSQL database interactions
- JWT authentication with HTTP-only cookies
- AWS S3 integration for image uploads
- Handles all business logic and data management

### **Frontend Service** (`frontend/`)

- Single-page application (SPA) built with React and Vite
- Redux Toolkit for centralized state management
- Material-UI component library
- Nginx serves static files and proxies API requests

### **Docker Compose** (`docker-compose.yml`)

- `backend-container` - Runs the Go application on port 4040
- `frontend-container` - Nginx serves React app on port 80 and proxies `/api/*` to backend
- Shared Docker network for inter-container communication

### **Nginx Configuration** (`nginx.conf`)

- Serves static React build files
- Handles client-side routing
- Reverse proxies API calls from `/api/*` to `http://backend:4040/*`
- Enables gzip compression and static asset caching

---

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1. **Clone the repository**

```bash
   git clone https://github.com/kk120306/cvwo-2026.git
   cd cvwo-2026
```

2. **Create environment file**

Create a `.env` file in the project root:

---

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | Backend server port | `4040` |
| `ENV` | Application environment | `development`, `production` |
| `DB_HOST` | PostgreSQL host (Neon) | `ep-*.neon.tech` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database username | `postgres` |
| `DB_PASSWORD` | Database password | `your_password` |
| `DB_NAME` | Database name | `mydb` |
| `DB_SSLMODE` | SSL mode for database | `require`, `disable` |
| `DB_CHANNEL_BINDING` | SCRAM channel binding | `require`, `prefer` |
| `SECRET` | JWT signing secret | `long_random_string` |
| `FRONTEND_URL` | Frontend URL for CORS | `http://localhost` |
| `AWS_ACCESS_KEY_ID` | AWS credentials | `AKIAIOSFODNN7EXAMPLE` |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key | `wJalrXUtnFEMI...` |
| `CLOUDFRONT_DISTRIBUTION_ID` | CloudFront distribution | `E1234567890ABC` |

---

3. **Start the application**

```bash
   docker-compose up -d
```

4. **Access the application**
   - **Frontend:** <http://localhost>
   - **Backend API:** <http://localhost:4040>

### Manual Deployment

To deploy manually on your own server:

```bash
# 1. SSH into your EC2 instance
ssh ubuntu@your-ec2-ip

# 2. Clone the repository
git clone https://github.com/kk120306/cvwo-2026.git
cd cvwo-2026

# 3. Create .env file with your secrets
nano .env

# 4. Start with Docker Compose
docker-compose up -d

# 5. Monitor logs
docker-compose logs -f
```

---

## 📧 Contact

**Author:** Kai Kameyama  
**GitHub:** [@kk120306](https://github.com/kk120306)
**Email:** <e1661384@u.nus.edu>

---

<p align="center">Built for CVWO 2026</p>
