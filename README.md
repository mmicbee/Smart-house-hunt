# Smart House Hunt

Smart House Hunt is a web-based rental property platform designed to help tenants find houses and help landlords manage property listings.

The project is being developed **milestone by milestone**, starting with a simple working platform and adding advanced security and AI features only after the core system is stable.

## Project Goal

The first goal is to build a working house-hunting platform where:

- A tenant can create an account and log in.
- A landlord can create an account and manage property listings.
- Tenants can search and view available houses.
- Tenants can save favorite houses.
- Tenants can request property viewings.
- An administrator can approve listings and manage users.
- Landlords can go through a verification process.

After the core platform is complete, advanced features such as recommendations, fraud detection, chatbot search, and rental price prediction will be added.

---

## Technology Stack

This project uses technologies that are currently comfortable for the developer and are suitable for learning, building, testing, and explaining the system.

| Part | Technology |
|---|---|
| Frontend | HTML, CSS, JavaScript |
| Backend | Go |
| API | REST API |
| Authentication | JWT |
| Database | SQLite 3 |
| Password Security | Go password hashing |
| Version Control | Git + GitHub |

### Future Technologies

The following are intentionally **not part of the first MVP**:

- Python / AI services
- Machine learning models
- PostgreSQL
- Cloud storage
- Docker/Kubernetes
- Advanced deployment infrastructure

They may be introduced later when the core application requires them.

---

## Project Structure

```text
smart-house-hunt/
│
├── frontend/
│   ├── index.html
│   ├── css/
│   │   └── style.css
│   ├── js/
│   │   ├── api.js
│   │   ├── auth.js
│   │   └── properties.js
│   ├── pages/
│   │   ├── login.html
│   │   ├── register.html
│   │   └── properties.html
│   └── assets/
│       └── images/
│
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── auth/
│   │   ├── users/
│   │   ├── properties/
│   │   ├── favorites/
│   │   ├── viewing/
│   │   └── verification/
│   │
│   ├── database/
│   │   ├── migrations/
│   │   └── seeders/
│   │
│   ├── middleware/
│   ├── routes/
│   ├── config/
│   └── storage/
│       ├── property-images/
│       ├── user-documents/
│       └── logs/
│
├── ai-engine/
│   ├── recommendation/
│   ├── chatbot/
│   ├── scam-detection/
│   ├── price-prediction/
│   └── models/
│
├── docs/
│   ├── proposal/
│   ├── srs/
│   ├── diagrams/
│   └── screenshots/
│
├── deployment/
│   ├── docker/
│   ├── nginx/
│   └── kubernetes/
│
├── .gitignore
└── README.md
```

## Architecture

The application will follow a simple client-server architecture:

```text
┌─────────────────────────┐
│        Frontend         │
│     HTML/CSS/JavaScript │
└────────────┬────────────┘
             │
             │ HTTP / REST API
             ▼
┌─────────────────────────┐
│       Go Backend        │
│                         │
│ Auth | Users | Houses   │
│ Favorites | Viewing     │
│ Verification | Admin   │
└────────────┬────────────┘
             │
             ▼
┌─────────────────────────┐
│        SQLite 3         │
└─────────────────────────┘
```

Later, the AI services can be connected through the Go backend:

```text
Frontend
   │
   ▼
Go REST API
   │
   ├──────────► SQLite 3
   │
   └──────────► AI Engine
```

The Go backend remains the main API and entry point for the frontend.

---

# Development Milestones

The project will be developed incrementally. Each milestone should be working before moving to the next one.

## Milestone 0 — Project Setup

**Goal:** Establish a working development environment.

Tasks:

- [ ] Set up project structure
- [ ] Initialize Go module
- [ ] Install backend dependencies
- [ ] Connect Go to SQLite 3
- [ ] Create basic HTTP server
- [ ] Create health-check endpoint
- [ ] Create frontend base pages
- [ ] Connect frontend to backend
- [ ] Set up Git branches and commits

### Expected result

The backend runs successfully and:

```text
GET /health
```

returns:

```json
{
  "status": "ok",
  "message": "Smart House Hunt API is running"
}
```

---

## Milestone 1 — Authentication

**Goal:** Allow users to securely create accounts and log in.

### Features

- [ ] User registration
- [ ] User login
- [ ] Password hashing
- [ ] JWT generation
- [ ] JWT validation middleware
- [ ] Logout
- [ ] Protected routes
- [ ] Role-based authorization

### Roles

- Tenant
- Landlord
- Agent
- Admin

### Expected result

A user can register, log in, receive a JWT, and access protected endpoints according to their role.

---

## Milestone 2 — Property Listings

**Goal:** Allow landlords to manage houses.

### Landlord

- [ ] Create property
- [ ] Edit property
- [ ] Delete property
- [ ] Upload property images
- [ ] View own listings

### Tenant

- [ ] View approved properties
- [ ] View property details

### Property status

```text
pending
approved
rejected
suspended
```

---

## Milestone 3 — Search and Favorites

**Goal:** Make it easy for tenants to find suitable houses.

### Search

- [ ] Search by location
- [ ] Filter by price
- [ ] Filter by bedrooms
- [ ] Filter by bathrooms
- [ ] Filter by amenities

### Favorites

- [ ] Save property
- [ ] Remove property
- [ ] View favorite properties

---

## Milestone 4 — Viewing Requests

**Goal:** Allow tenants to request a property viewing.

Tenant:

```text
Select Property
      ↓
Request Viewing
      ↓
Select Preferred Date
      ↓
Landlord Receives Request
```

Features:

- [ ] Create viewing request
- [ ] View requests
- [ ] Accept request
- [ ] Reject request
- [ ] Track request status

---

## Milestone 5 — Admin Dashboard

**Goal:** Give administrators control over the platform.

Admin features:

- [ ] View users
- [ ] Suspend users
- [ ] View property listings
- [ ] Approve listings
- [ ] Reject listings
- [ ] View reports

---

## Milestone 6 — Landlord Verification

**Goal:** Increase trust between tenants and landlords.

Process:

```text
Landlord
   ↓
Upload ID
   ↓
Upload Selfie
   ↓
Verification Request
   ↓
Admin Review
   ↓
Approved
   ↓
Verified Badge
```

Planned table:

```text
verification_requests

id
user_id
id_document
selfie
status
reviewed_by
created_at
updated_at
```

---

## Milestone 7 — Security Improvements

After the core platform is stable:

- [ ] Email OTP
- [ ] Multi-factor authentication
- [ ] Login monitoring
- [ ] IP logging
- [ ] Device information
- [ ] Audit logs
- [ ] Security alerts

Example audit events:

```text
LOGIN
CREATE_PROPERTY
UPDATE_PROPERTY
DELETE_PROPERTY
APPROVE_PROPERTY
SUSPEND_USER
APPROVE_VERIFICATION
```

---

# Future AI Milestones

AI will **not** be built before the core platform is working.

## AI Recommendation Engine

Use information such as:

- Location
- Budget
- Bedrooms
- Bathrooms
- Amenities
- Favorites
- Search history

Example:

```text
House A → 95% match
House B → 87% match
House C → 75% match
```

## AI Chatbot

Example:

```text
User:
I need a 2-bedroom house below 15,000 in Kakamega.
```

The system can eventually convert this into search filters:

```text
location = Kakamega
bedrooms = 2
max_price = 15000
```

## Fraud Detection

Potential checks:

- Duplicate property images
- Suspicious rental prices
- Repeated complaints
- Suspicious landlord activity

Example:

```text
Average rent: 15,000 KES
Listing price: 4,000 KES

→ Potentially suspicious
```

## Rental Price Prediction

The future model may use:

- Location
- Bedrooms
- Bathrooms
- Amenities
- Historical rental prices

Output:

```text
Expected rent: KES 18,000
```

---

# Database

SQLite 3 will be used during development and for the initial project.

Planned core tables:

```text
users
properties
property_images
favorites
viewing_requests
messages
verification_requests
fraud_reports
audit_logs
notifications
```

The database schema will be introduced gradually as each milestone requires it.

---

# Getting Started

## Prerequisites

Install:

- Go 1.22 or newer
- Git
- A modern web browser

No Node.js, React, PostgreSQL, Python, or Docker is required for the initial version.

## Clone the Repository

```bash
git clone https://github.com/YOUR_USERNAME/smart-house-hunt.git
cd smart-house-hunt
```

## Start the Backend

```bash
cd backend
go mod tidy
go run ./cmd/api
```

The API will run on:

```text
http://localhost:8080
```

Test the health endpoint:

```text
http://localhost:8080/health
```

## Open the Frontend

The frontend is currently plain HTML/CSS/JavaScript.

You can open:

```text
frontend/index.html
```

directly in your browser during the early stages.

For a better local development experience, you can also use a simple static server.

---

# Git Workflow

Development will use feature branches.

Example:

```bash
git checkout -b feature/authentication
```

After completing the work:

```bash
git add .
git commit -m "feat: add user authentication"
git push origin feature/authentication
```

Recommended branch types:

```text
feature/authentication
feature/property-listings
feature/property-search
feature/favorites
feature/viewing-requests
feature/admin-dashboard
feature/landlord-verification
feature/security
feature/ai-recommendations
```

---

# Development Principle

The project follows one main rule:

> **Build the simple version first. Improve it after it works.**

We will not introduce complicated technologies just because they are popular.

The initial stack is intentionally simple:

```text
HTML
CSS
JavaScript
   ↓
Go
   ↓
JWT
   ↓
SQLite 3
```

Once the foundation is stable, advanced technologies can be introduced where they provide real value.

---

# Current Status

**Milestone 0 — Project Setup**

The project structure has been prepared for the new stack.

Next development target:

1. Connect Go to SQLite 3
2. Create the initial database schema
3. Build the authentication module
4. Implement registration
5. Implement login
6. Add JWT authentication

---

## Project Vision

Smart House Hunt aims to become a trusted digital platform for discovering rental properties, connecting tenants and landlords, and eventually using intelligent systems to make property discovery safer and easier.

**Build small. Test. Understand. Improve.**
