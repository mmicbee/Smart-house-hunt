# Smart House Hunt — Project Flow

## High-Level Architecture

```text
┌────────────────────┐     ┌─────────────────┐     ┌──────────────┐
│     Frontend       │────▶│     Backend     │────▶│   SQLite 3   │
│  (HTML / CSS / JS) │     │    (Go / Gin)   │     │   (Database) │
└────────────────────┘     └────────┬────────┘     └──────────────┘
                                    │
                                    ├────▶ Local Storage (`storage/`)
                                    │
                                    ▼ (Future Plan)
                             ┌──────────────┐
                             │  AI Engine   │
                             │   (Python)   │
                             └──────────────┘
```

## User Flows

### Tenant Flow
1. Register / Login
2. Browse approved properties (`GET /api/properties`)
3. View property details & photos (`GET /api/properties/:id`)
4. (Milestone 3) Search & filter properties
5. (Milestone 3) Save to favorites
6. (Milestone 4) Contact landlord / Request viewing
7. (Future) Chat with AI assistant
8. (Future) Get personalized recommendations

### Landlord Flow
1. Register / Login
2. Create property listing (`POST /api/properties`)
3. Upload listing photos (`POST /api/properties/:id/images`)
4. Manage listings (edit / delete / view own via `/properties/mine`)
5. (Milestone 4) Receive viewing inquiries
6. (Milestone 6) Complete verification (ID + Selfie)
7. (Future) See price prediction suggestions

### Admin Flow
1. Login
2. Review & approve listings (Milestone 5)
3. Review verification requests (Milestone 6)
4. Suspend fraudulent users
5. View reports & audit logs

## Development Milestones

**Milestone 0: Project Setup** (Completed)
- Clean project structure, Go module, Gin HTTP server, SQLite 3 connection, health check endpoint.

**Milestone 1: Authentication** (Completed)
- User registration, login, logout, password hashing (bcrypt), JWT generation & validation middleware, role-based authorization.

**Milestone 2: Property Listings** (Completed)
- Landlord: create, edit, delete properties, upload photos to local storage, view own listings.
- Tenant: browse approved properties, view property details and images.

**Milestone 3: Search and Favorites** (Next)
- Search by location, price range, bedrooms.
- Save and view favorite properties.

**Milestone 4: Viewing Requests & Messaging**
- Book house visit, basic tenant-landlord communication.

**Milestone 5: Admin Dashboard**
- Review & approve listings, manage users.

**Milestone 6: Landlord Verification & Security**
- Upload ID and selfie, manual admin verification, secure document handling.

**Future Plans (After Core Platform is Stable)**
- AI Recommendations (Python / FastAPI)
- AI Chatbot for natural search
- Fraud / scam detection (duplicate images, suspicious pricing)
- Rental price prediction
- MinIO / AWS S3 for object storage
- PostgreSQL migration for production scale
