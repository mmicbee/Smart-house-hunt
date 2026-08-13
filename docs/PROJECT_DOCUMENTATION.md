Since you're studying IT and now have Go experience from Zone01, I'd recommend building this as a modern full-stack project with a clean architecture. That will make development easier and impress your lecturers if they review the codebase.

Project Structure

smart-house-hunt/
│
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── assets/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── layouts/
│   │   ├── services/
│   │   ├── hooks/
│   │   ├── context/
│   │   ├── routes/
│   │   └── utils/
│   │
│   └── package.json
│
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── auth/
│   │   ├── users/
│   │   ├── tenants/
│   │   ├── landlords/
│   │   ├── agents/
│   │   ├── properties/
│   │   ├── recommendations/
│   │   ├── chatbot/
│   │   ├── fraud/
│   │   ├── security/
│   │   ├── notifications/
│   │   └── reports/
│   │
│   ├── database/
│   │   ├── migrations/
│   │   └── seeders/
│   │
│   ├── middleware/
│   ├── routes/
│   ├── config/
│   ├── storage/
│   │   ├── property-images/
│   │   ├── user-documents/
│   │   └── logs/
│   │
│   └── go.mod
│
├── ai-engine/
│   ├── recommendation/
│   ├── price-prediction/
│   ├── scam-detection/
│   ├── face-verification/
│   ├── chatbot/
│   ├── models/
│   └── app.py
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
└── README.md


---

Phase 1: MVP (Minimum Viable Product)

Build this first before touching AI.

Authentication

Features

Register

Login

Logout

Password reset

JWT authentication


Roles

Tenant

Landlord

Agent

Admin



---

Property Module

Landlord

Create listing

Edit listing

Delete listing

Upload photos


Tenant

Search houses

View details

Save favorites



---

Admin

Approve listings

Suspend users

View reports



---

Phase 2: Verification System

This is where your project starts becoming unique.

Landlord Verification

Upload ID
      |
      v
Upload Selfie
      |
      v
Admin Review
      |
      v
Verified Badge

Database table:

verification_requests

Fields:

id
user_id
id_document
selfie
status
reviewed_by
created_at


---

Phase 3: AI Recommendation Engine

User Preferences

Store:

Location
Budget
Bedrooms
Bathrooms
Amenities

AI suggests:

Recommended Houses
Match Score

Example:

House A → 95%
House B → 87%
House C → 75%


---

Phase 4: AI Chatbot

User types:

Need a 2-bedroom house below 15k in Kakamega.

Bot translates this into search filters.

Returns matching listings.


---

Phase 5: Fraud Detection

This is a major selling point.

Checks

Duplicate Images

Image A
     ==
Image B

Flag listing.


---

Suspicious Pricing

Example:

Average rent:

KES 15,000

Listed:

KES 4,000

Flag as suspicious.


---

Complaint Analysis

If many users report a landlord:

Fraud Risk: High


---

Phase 6: Security System

Multi-Factor Authentication

Email OTP

or

SMS OTP


---

Login Monitoring

Store:

IP Address
Device
Country
Time

Send alert:

New login detected.


---

Audit Logs

Table:

audit_logs

Tracks:

Login
Delete Property
Update Listing
Verification Approval


---

Phase 7: Rental Price Prediction

AI model trains on:

Location
Bedrooms
Bathrooms
Amenities

Output:

Expected Rent:
KES 18,000

Useful for both landlords and tenants.


---

Database Structure

Core tables:

users
roles
properties
property_images
favorites
viewing_requests
messages
verification_requests
recommendations
fraud_reports
audit_logs
notifications


---

Recommended Tech Stack

Frontend

React

Tailwind CSS



---

Backend

Go

Gin Framework

JWT



---

Database

PostgreSQL



---

AI

Python

FastAPI

Scikit-Learn



---

Storage

MinIO (local development)

AWS S3 (production)



---

DevOps

Docker

GitHub Actions

Nginx



---

What I'd Build First (Week 1–4)

1. Authentication


2. Property Listings


3. Search System


4. Favorites


5. Admin Dashboard


6. Landlord Verification



Only after these are working would I start:

AI Recommendations

AI Chatbot

Fraud Detection

Price Prediction


This order prevents you from spending months on AI before you have actual property and user data to work with. The AI features become much easier once the core platform is already functioning.