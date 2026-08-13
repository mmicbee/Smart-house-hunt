# Smart House Hunt

A modern full-stack platform for finding and listing rental houses in Kenya, with AI-powered recommendations, fraud detection, landlord verification, and more.

Built as a student project with clean architecture, focusing on real-world features for tenants, landlords, agents, and admins.

## Project Vision

Help people in Kenya (starting with areas like Kakamega and beyond) find genuine rental houses easily, while protecting them from scams through verification, AI fraud detection, and transparent pricing.

## Tech Stack

| Layer          | Technology                          |
|----------------|-------------------------------------|
| Frontend       | React + Tailwind CSS + Vite         |
| Backend        | Go (Gin framework) + JWT            |
| Database       | PostgreSQL                          |
| AI Engine      | Python + FastAPI + Scikit-Learn     |
| Storage        | MinIO (dev) / AWS S3 (prod)         |
| DevOps         | Docker, Nginx, GitHub Actions       |

## Project Structure

```
smart-house-hunt/
├── frontend/                 # React application
├── backend/                  # Go API server
├── ai-engine/                # Python AI services
├── docs/                     # Documentation (proposal, SRS, diagrams)
├── deployment/               # Docker, Nginx, Kubernetes configs
└── README.md
```

## Development Phases

### Phase 1: MVP (Minimum Viable Product) — **Current Focus**
- Authentication (Register, Login, Logout, Password Reset, JWT)
- Roles: Tenant, Landlord, Agent, Admin
- Property Module:
  - Landlord: Create / Edit / Delete listings + photo upload
  - Tenant: Search, View details, Save favorites
- Admin: Approve listings, Suspend users, View reports

### Phase 2: Verification System
- Landlord ID + Selfie upload → Admin review → Verified Badge

### Phase 3: AI Recommendation Engine
- User preferences → Match score for properties

### Phase 4: AI Chatbot
- Natural language search (e.g. "2-bedroom under 15k in Kakamega")

### Phase 5: Fraud Detection
- Duplicate image detection
- Suspicious pricing
- Complaint analysis

### Phase 6: Security
- MFA (Email/SMS OTP)
- Login monitoring + alerts
- Audit logs

### Phase 7: Rental Price Prediction
- AI model predicting expected rent based on location, bedrooms, etc.

## Getting Started

### Prerequisites
- Node.js 20+
- Go 1.22+
- Python 3.11+
- PostgreSQL 16+
- Docker (optional but recommended)

### 1. Clone the repository
```bash
git clone https://github.com/YOUR_USERNAME/smart-house-hunt.git
cd smart-house-hunt
```

### 2. Backend Setup
```bash
cd backend
go mod tidy
# Configure database connection in config/
go run cmd/api/main.go
```

### 3. Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

### 4. AI Engine (later phases)
```bash
cd ai-engine
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
uvicorn app:app --reload
```

## Core Database Tables (Planned)
- users, roles
- properties, property_images
- favorites, viewing_requests
- messages
- verification_requests
- recommendations
- fraud_reports
- audit_logs
- notifications

## Roadmap Priority
1. Authentication
2. Property Listings (CRUD + Search)
3. Favorites
4. Admin Dashboard
5. Landlord Verification
6. AI features (only after we have real data)

---

**Status**: Phase 1 in progress — Property Listing MVP

Made with ❤️ for Kenyan house hunters.
