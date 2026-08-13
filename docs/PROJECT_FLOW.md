# Smart House Hunt — Project Flow

## High-Level Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Frontend  │────▶│   Backend   │────▶│  PostgreSQL │
│  (React)    │     │   (Go/Gin)  │     │             │
└─────────────┘     └──────┬──────┘     └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │  AI Engine  │
                    │  (Python)   │
                    └─────────────┘
```

## User Flows

### Tenant Flow
1. Register / Login
2. Browse / Search properties
3. View property details
4. Save to favorites
5. Contact landlord / Request viewing
6. (Later) Chat with AI assistant
7. (Later) Get personalized recommendations

### Landlord Flow
1. Register / Login
2. Complete verification (ID + Selfie)
3. Create property listing + upload photos
4. Manage listings (edit / delete)
5. Receive inquiries
6. (Later) See price prediction suggestions

### Admin Flow
1. Login
2. Review & approve listings
3. Review verification requests
4. Suspend fraudulent users
5. View reports & audit logs

## Development Order (Recommended)

**Week 1-2**
- Project setup + GitHub
- Authentication (JWT)
- Basic property CRUD

**Week 3-4**
- Search & filters
- Property details page
- Favorites
- Admin basic dashboard

**Week 5-6**
- Landlord verification system
- Image upload (MinIO/S3)

**Later**
- AI Recommendation
- Chatbot
- Fraud detection
- Price prediction
- MFA & advanced security
