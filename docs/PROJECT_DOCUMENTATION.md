# Smart House Hunt — Project Documentation

The project follows one main principle:
> **Build the simple version first. Improve it after it works.**

We intentionally begin with a clean, lightweight architecture using plain HTML/CSS/JavaScript on the frontend, Go (Gin) on the backend, JWT authentication, and SQLite 3 for data persistence. This ensures that every component is easy to understand, test, and debug before any advanced technologies or AI services are introduced.

---

## Project Structure

```text
smart-house-hunt/
│
├── frontend/
│   ├── css/
│   │   └── style.css
│   ├── js/
│   │   ├── api.js
│   │   ├── auth.js
│   │   └── properties.js
│   ├── pages/
│   │   ├── login.html
│   │   ├── register.html
│   │   ├── properties.html
│   │   ├── property-details.html
│   │   ├── my-listings.html
│   │   └── create-property.html
│   ├── assets/
│   │   └── images/
│   └── index.html
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
│   │   ├── favorites/        (Milestone 3)
│   │   ├── viewing/          (Milestone 4)
│   │   └── verification/     (Milestone 6)
│   │
│   ├── database/
│   │   ├── migrations/
│   │   │   ├── 001_initial.sql
│   │   │   ├── 002_properties.sql
│   │   │   └── 003_property_images.sql
│   │   └── database.go
│   │
│   ├── middleware/
│   │   ├── auth.go
│   │   └── cors.go
│   ├── routes/
│   ├── config/
│   ├── storage/
│   │   ├── property-images/
│   │   ├── user-documents/
│   │   └── logs/
│   │
│   ├── go.mod
│   └── go.sum
│
├── ai-engine/ (Future Plan)
│   ├── recommendation/
│   ├── price-prediction/
│   ├── scam-detection/
│   ├── face-verification/
│   ├── chatbot/
│   └── models/
│
├── docs/
│   ├── PROJECT_DOCUMENTATION.md
│   └── PROJECT_FLOW.md
│
├── deployment/ (Future Plan)
│   ├── docker/
│   ├── nginx/
│   └── kubernetes/
│
└── README.md
```

---

## Core Architecture & Stack

### Current Active Stack

| Component | Technology | Rationale |
| --- | --- | --- |
| **Frontend** | HTML5, CSS3, Vanilla JS (Fetch API) | Zero build step, fast loading, accessible, easy to inspect. |
| **Backend** | Go 1.22+, Gin Framework | High performance, simple concurrency, clean idiomatic Go. |
| **Auth** | JWT (`golang-jwt/jwt/v5`), bcrypt | Stateless authentication, role-based authorization gates. |
| **Database** | SQLite 3 (`mattn/go-sqlite3`) | Embedded, serverless, zero-configuration, robust SQL transactions. |
| **File Storage** | Local filesystem (`backend/storage/`) | Direct file serving via `/uploads`, git-ignored directories. |

### Planned / Future Evolutions (Post-MVP)

| Component | Future Technology | Trigger for Adoption |
| --- | --- | --- |
| **Frontend** | React / Tailwind CSS | When complex client-side state or rich interactivity warrants a SPA build pipeline. |
| **Database** | PostgreSQL | When migrating to multi-instance server deployments or heavy concurrent write loads. |
| **Storage** | MinIO (local) / AWS S3 (production) | When moving image and document storage off the host server to cloud object storage. |
| **AI Engine** | Python (FastAPI, Scikit-Learn, PyTorch) | Once real property and user interaction datasets exist to train models. |

---

## System Modules & Roadmap

### Phase 1: MVP (Milestones 0 – 2)

#### Authentication (Milestone 1 — Complete)
- User registration with roles (`tenant`, `landlord`, `agent`, `admin`)
- User login with password verification (bcrypt)
- JWT token generation and validation middleware
- Role-based authorization middleware
- Current user profile endpoint (`GET /api/me`)
- User logout

#### Property Listings (Milestone 2 — Complete)
- **Landlord features:**
  - Create listing (`POST /api/properties`)
  - Edit own listing (`PUT /api/properties/:id`)
  - Delete own listing (`DELETE /api/properties/:id`)
  - Upload listing photos to local storage (`POST /api/properties/:id/images`)
  - View own listings (`GET /api/properties/mine` & `my-listings.html`)
  - Ownership enforcement (403 for unauthorized landlord edits/deletes)
- **Tenant / Public features:**
  - Browse approved properties (`GET /api/properties` & `properties.html`)
  - View property details with full photo gallery (`GET /api/properties/:id` & `property-details.html`)
  - Non-approved and missing properties return 404

---

### Phase 2: Discovery & Interaction (Milestones 3 – 4)

#### Search & Favorites (Milestone 3)
- Search by location, rent budget, bedrooms, and bathrooms.
- Save properties to favorites for quick reference.

#### Viewing Requests & Messaging (Milestone 4)
- Tenant schedules house viewing visits.
- Basic messaging between tenants and landlords.

---

### Phase 3: Administration & Trust (Milestones 5 – 6)

#### Admin Dashboard (Milestone 5)
- Review and approve pending property listings.
- Suspend accounts violating community guidelines.
- System metrics and audit trail.

#### Landlord Verification (Milestone 6)
- Landlord uploads national ID and selfie.
- Admin reviews documents in `storage/user-documents/`.
- Verified badge issued upon confirmation to reduce rental fraud.

---

### Phase 4: Advanced & AI Features (Future Plans)

These features will be built in the `ai-engine/` Python service only after the core platform has real user activity and listing data:

#### AI Recommendation Engine (Planned)
- Learns tenant search habits, budget, and location preferences.
- Calculates listing match scores (e.g. 95% match).

#### AI Natural Language Chatbot (Planned)
- Converts tenant search prompts (e.g. "2-bedroom under 15k in Kakamega") into structured API queries.

#### Fraud & Scam Detection (Planned)
- Perceptual image hashing to flag duplicate images across different landlords.
- Outlier price detection (e.g. listing 4k in a 20k neighborhood).
- Complaint volume and pattern monitoring.

#### Rental Price Prediction (Planned)
- Machine learning regression model estimating fair market rent based on location, size, and amenities.

---

## Database Schema (Core Tables)

### Active Tables (SQLite 3)

```sql
-- roles
CREATE TABLE roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

-- users
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

-- properties
CREATE TABLE properties (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    landlord_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    location TEXT NOT NULL,
    price INTEGER NOT NULL,
    bedrooms INTEGER NOT NULL DEFAULT 1,
    bathrooms INTEGER NOT NULL DEFAULT 1,
    amenities TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (landlord_id) REFERENCES users(id)
);

-- property_images
CREATE TABLE property_images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    property_id INTEGER NOT NULL,
    file_path TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (property_id) REFERENCES properties(id)
);
```

### Planned Milestone Tables

- `favorites`: `(id, user_id, property_id, created_at)`
- `viewing_requests`: `(id, user_id, property_id, preferred_date, status, created_at)`
- `verification_requests`: `(id, user_id, id_document, selfie, status, reviewed_by, created_at)`
- `audit_logs`: `(id, user_id, action, ip_address, created_at)`

---

## Summary of Completed Work

1. **Milestone 0 — Project Setup:** Go module, SQLite connection, Gin server, health check.
2. **Milestone 1 — Authentication:** Register, login, JWT middleware, roles, protected routes.
3. **Milestone 2 — Property Listings:** Full landlord management (CRUD, photos, my listings), tenant listing view, and details page.