# CVWO Forum

[![Ask DeepWiki](https://devin.ai/assets/askdeepwiki.png)](https://deepwiki.com/Kk120306/cvwo-2026)

A full-stack forum application built with a Go backend and React frontend. Users can engage in discussions, create posts under various topics, comment, and vote. The application features complete authentication, administrative roles, and globally-distributed image hosting.

**🚀 Live Demo:** [http://3.25.177.70/](#)

---

## 📚 Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Performance](#performance)
- [Getting Started](#getting-started)
- [Deployment](#deployment)

---

## ✨ Features

### Core Functionality
- **User Authentication** - Secure signup and login with JWT-based sessions stored in HTTP-only cookies
- **Topics Management** - Admins can create, update, and delete topics for organizing discussions
- **Post Management** - Full CRUD operations for posts with rich text editing
- **Voting System** - Upvote/downvote posts and comments
- **Comments** - Threaded discussions with edit and delete capabilities
- **User Profiles** - View post and comment history with user statistics

### Rich Content
- **TipTap Editor** - Rich text editing with formatting, links, and embedded media
- **Image Uploads** - Seamless image handling via pre-signed S3 URLs
- **CDN Distribution** - Fast global image delivery through CloudFront

### Admin Features
- **Pin Posts** - Highlight important posts at the top of topic feeds
- **Topic Management** - Create and organize discussion categories
- **Moderation Tools** - Delete inappropriate content

### User Experience
- **Responsive Design** - Clean, mobile-friendly UI with Material-UI
- **Client-Side Filtering** - Real-time search and sort for posts and comments
- **Fast Performance** - Optimized with CDN and efficient backend

---

## 🛠️ Tech Stack

### Backend
- **[Go](https://golang.org/)** - High-performance backend language
- **[Gin](https://gin-gonic.com/)** - Lightweight web framework
- **[GORM](https://gorm.io/)** - Elegant ORM for database operations
- **[PostgreSQL](https://www.postgresql.org/)** - Robust relational database (via [Neon](https://neon.tech/))

### Frontend
- **[React](https://reactjs.org/)** - Component-based UI library
- **[TypeScript](https://www.typescriptlang.org/)** - Type-safe JavaScript
- **[Vite](https://vitejs.dev/)** - Lightning-fast build tool
- **[Redux Toolkit](https://redux-toolkit.js.org/)** - State management
- **[Material-UI](https://mui.com/)** - Modern component library
- **[TipTap](https://tiptap.dev/)** - Headless rich text editor

### Infrastructure & DevOps
- **[Docker](https://www.docker.com/)** - Containerization
- **[Docker Compose](https://docs.docker.com/compose/)** - Multi-container orchestration
- **[Nginx](https://www.nginx.com/)** - Reverse proxy and static file serving
- **[AWS EC2](https://aws.amazon.com/ec2/)** - Cloud hosting
- **[AWS S3](https://aws.amazon.com/s3/)** - Object storage for images
- **[AWS CloudFront](https://aws.amazon.com/cloudfront/)** - Global CDN
- **[GitHub Actions](https://github.com/features/actions)** - CI/CD pipeline

---

## 🏗️ Architecture

This monorepo contains two main services orchestrated with Docker Compose:

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
Orchestrates both services with:
- `backend-container` - Runs the Go application on port 4040
- `frontend-container` - Nginx serves React app on port 80 and proxies `/api/*` to backend
- Shared Docker network for inter-container communication

### **Nginx Configuration** (`nginx.conf`)
- Serves static React build files
- Handles client-side routing (React Router)
- Reverse proxies API calls from `/api/*` to `http://backend:4040/*`
- Enables gzip compression and static asset caching

---

## ⚡ Performance

### CloudFront CDN Impact
Implementing AWS CloudFront for image distribution dramatically improved global performance:

- **Before CloudFront:** 300+ ms response time for Germany servers
- **After CloudFront:** 0.63 ms response time for Germany servers
- **Improvement:** ~99.8% reduction in latency

This demonstrates the power of edge caching for globally distributed content delivery.

---

## 🚀 Getting Started

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
```env
   # Server Configuration
   PORT=4040
   ENV=production
   
   # Database Configuration (Neon PostgreSQL)
   DB_HOST=your-neon-host.neon.tech
   DB_PORT=5432
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=your_database
   DB_SSLMODE=require
   DB_CHANNEL_BINDING=require
   
   # Security
   SECRET=your_jwt_secret_key_here
   FRONTEND_URL=http://localhost
   
   # AWS Configuration
   AWS_ACCESS_KEY_ID=your_access_key
   AWS_SECRET_ACCESS_KEY=your_secret_key
   CLOUDFRONT_DISTRIBUTION_ID=your_distribution_id
```

3. **Start the application**
```bash
   docker-compose up -d
```

4. **Access the application**
   - **Frontend:** http://localhost
   - **Backend API:** http://localhost:4040

### Environment Variables Reference

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

## 🚢 Deployment

### CI/CD Pipeline

This project uses **GitHub Actions** for continuous deployment. The workflow (`.github/workflows/fullstack.yml`) automatically builds and deploys on every push to `main`.

#### Workflow Steps:

1. **Build Phase** (Parallel)
   - `build-backend` - Builds Go backend Docker image
   - `build-frontend` - Builds React frontend Docker image
   - Both images are pushed to Docker Hub

2. **Deploy Phase** (Self-Hosted EC2 Runner)
   - Pulls latest images from Docker Hub
   - Creates `.env` file from GitHub Secrets
   - Stops old containers with `docker-compose down`
   - Starts new containers with `docker-compose up -d`
   - Cleans up unused Docker images

#### Deployment Infrastructure:

- **Hosting:** AWS EC2 instance (self-hosted GitHub Actions runner)
- **Database:** Neon PostgreSQL (serverless)
- **Storage:** AWS S3 (image uploads)
- **CDN:** AWS CloudFront (global distribution)
- **Orchestration:** Docker Compose

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

### GitHub Secrets Configuration

Add these secrets in **Settings → Secrets and variables → Actions**:

- `DOCKER_USERNAME` - Docker Hub username
- `DOCKER_PASSWORD` - Docker Hub password
- All environment variables listed in the table above

---

## 📝 License

This project is part of the CVWO (Computing for Voluntary Welfare Organisations) assignment.

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

## 📧 Contact

**Author:** Kai Kameyama  
**GitHub:** [@kk120306](https://github.com/kk120306)

---

<p align="center">Built with ❤️ for CVWO 2026</p>