# Campus Forum

[Short project description - max ~70 characters]

Campus Forum is a feature-rich campus community forum system built with a frontend-backend separated architecture. The backend is developed in Go with the Gin framework, providing efficient and stable RESTful API services; the frontend is a single-page application built with Vue 3, delivering a smooth user experience. It supports complete forum functionalities including article posting, comment interaction, and user authentication.

## Technology Stack

### Backend
- **Language**: Go 1.21
- **Framework**: Gin
- **Database**: MySQL + GORM
- **Authentication**: JWT
- **File Storage**: WebDAV

### Frontend
- **Framework**: Vue 3
- **Routing**: Vue Router
- **Build Tool**: Vite

### Deployment
- **Containerization**: Docker + Docker Compose
- **Web Server**: Nginx

## Features

### User System
- User registration and login
- JWT authentication
- Profile management
- Avatar upload

### Article System
- Create, edit, and delete articles
- Categorized article browsing
- Markdown-rich text support
- Like/unlike articles
- Comment interaction

### Content Management
- Category management
- Site announcement publishing
- Sidebar configuration
- Deletion request approval workflow

### File Management
- Avatar upload
- Image upload
- Video upload
- WebDAV cloud storage support

## Project Structure

```
campus-forum/
├── backend/              # Go backend service
│   ├── controllers/      # Controller layer
│   ├── database/         # Database initialization
│   ├── middleware/       # Middleware
│   ├── models/           # Data models
│   ├── utils/            # Utility functions
│   ├── main.go           # Entry file
│   └── config.json       # Configuration file
├── frontend/             # Vue frontend application
│   ├── src/
│   │   ├── api/          # API interfaces
│   │   ├── components/   # Shared components
│   │   ├── views/        # Page views
│   │   ├── App.vue
│   │   └── main.js
│   └── index.html
├── nginx/                # Nginx configuration
├── docker-compose.yml    # Container orchestration
└── build.sh              # Build script
```

## Quick Start

### Prerequisites

- Docker >= 20.10
- Docker Compose >= 2.0

### Deploy with Docker Compose

```bash
# Clone the project
git clone https://gitee.com/byusistudio/campus-forum.git
cd campus-forum

# Start all services
docker-compose up -d
```

After startup, services will be available on:
- Frontend: http://localhost
- Backend API: http://localhost:8080

### Manual Deployment

#### Start Backend

```bash
cd backend

# Modify config.json
# Configure database connection, WebDAV, and other parameters

# Run the service
go run main.go
```

#### Start Frontend

```bash
cd frontend

# Install dependencies
npm install

# Development mode
npm run dev

# Production build
npm run build
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST   | /api/auth/register | User registration |
| POST   | /api/auth/login | User login |
| GET    | /api/auth/profile | Get user profile |
| PUT    | /api/auth/profile | Update user profile |
| GET    | /api/articles | Get article list |
| POST   | /api/articles | Create article |
| GET    | /api/articles/:id | Get article details |
| PUT    | /api/articles/:id | Update article |
| DELETE | /api/articles/:id | Delete article |
| POST   | /api/articles/:id/like | Like article |
| DELETE | /api/articles/:id/like | Unlike article |
| POST   | /api/articles/:id/comments | Add comment |
| DELETE | /api/comments/:id | Delete comment |
| GET    | /api/categories | Get category list |
| POST   | /api/categories | Create category |
| GET    | /api/announcements | Get announcements |
| PUT    | /api/announcements | Update announcement |

## Configuration

### Backend Configuration (backend/config.json)

```json
{
  "server": {
    "port": "8080"
  },
  "database": {
    "dsn": "user:password@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  },
  "jwt": {
    "secret": "your-jwt-secret-key"
  },
  "webdav": {
    "url": "https://your-webdav-server.com",
    "username": "your-username",
    "password": "your-password"
  }
}
```

## License

This project is open-sourced under the MIT License. See [LICENSE](./LICENSE) for details.

## Contributors

Contributions via Issues and Pull Requests are welcome!