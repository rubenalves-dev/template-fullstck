# Beheer API Documentation

This document describes the available endpoints, their inputs, and outputs.

## Base URL
`http://localhost:8080`

## Authentication
Most endpoints require a JWT token in the `Authorization` header:
`Authorization: Bearer <your-token>`

---

## Public Endpoints

### Health Check
Check if the service is up.
- **URL:** `/health`
- **Method:** `GET`
- **Response:** `200 OK`
  ```json
  {
    "data": {
      "status": "OK"
    }
  }
  ```

### Login
Authenticate and receive a JWT token.
- **URL:** `/login`
- **Method:** `POST`
- **Body:**
  ```json
  {
    "email": "user@example.com",
    "password": "yourpassword"
  }
  ```
- **Response:** `200 OK`
  ```json
  {
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5..."
    }
  }
  ```

### Register
Create a new organization and an admin user.
- **URL:** `/register`
- **Method:** `POST`
- **Body:**
  ```json
  {
    "email": "admin@example.com",
    "password": "yourpassword",
    "full_name": "Admin User",
    "organization_name": "My Sports Club"
  }
  ```
- **Response:** `201 Created`
  ```json
  {
    "data": {
      "message": "User registered successfully"
    }
  }
  ```

---

## Protected Endpoints

### Members

#### Create Member
- **URL:** `/api/v1/members`
- **Method:** `POST`
- **Body:**
  ```json
  {
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "+351912345678",
    "address": "Rua Central 123",
    "zip_code": "1234-567",
    "city": "Porto",
    "country": "Portugal",
    "nif": "123456789",
    "document_type": "CC",
    "document_value": "12345678",
    "birthday": "1990-01-01",
    "gender": "MALE"
  }
  ```
- **Response:** `201 Created`

#### List Members
- **URL:** `/api/v1/members`
- **Method:** `GET`
- **Response:** `200 OK`

#### Get Member
- **URL:** `/api/v1/members/{id}`
- **Method:** `GET`
- **Response:** `200 OK`

#### Update Member
- **URL:** `/api/v1/members/{id}`
- **Method:** `PUT`
- **Body:** Same as Create Member.
- **Response:** `200 OK`

#### Archive Member
- **URL:** `/api/v1/members/{id}`
- **Method:** `DELETE`
- **Response:** `204 No Content`

#### Get Full Profile
- **URL:** `/api/v1/members/{id}/profile`
- **Method:** `GET`
- **Response:** `200 OK`

---

### Modalities & Seasons

#### Create Modality
- **URL:** `/api/v1/modalities`
- **Method:** `POST`
- **Body:**
  ```json
  {
    "name": "Judo",
    "notes": "Notes about the modality"
  }
  ```
- **Response:** `201 Created`

#### List Modalities
- **URL:** `/api/v1/modalities`
- **Method:** `GET`
- **Response:** `200 OK`

#### List Weight Categories
- **URL:** `/api/v1/modalities/{modalityID}/weight-categories`
- **Method:** `GET`
- **Response:** `200 OK`

#### Create Weight Category
- **URL:** `/api/v1/modalities/weight-categories`
- **Method:** `POST`
- **Body:**
  ```json
  {
    "modality_id": "uuid",
    "name": "-60kg",
    "min_weight": 0,
    "max_weight": 60,
    "gender": "MALE"
  }
  ```
- **Response:** `201 Created`

#### Create Season
- **URL:** `/api/v1/seasons`
- **Method:** `POST`
- **Body:**
  ```json
  {
    "name": "2025/2026",
    "start_at": "2025-09-01T00:00:00Z",
    "end_at": "2026-08-31T23:59:59Z"
  }
  ```
- **Response:** `201 Created`

#### List Seasons
- **URL:** `/api/v1/seasons`
- **Method:** `GET`
- **Response:** `200 OK`

#### Create Modality Price
- **URL:** `/api/v1/seasons/prices`
- **Method:** `POST`
- **Body:**
  ```json
  {
    "season_id": "uuid",
    "modality_id": "uuid",
    "enrollment_fee": "20.00",
    "annuity_fee": "15.00",
    "monthly_fee": "35.00"
  }
  ```
- **Response:** `201 Created`

#### List Prices by Season
- **URL:** `/api/v1/seasons/{seasonID}/prices`
- **Method:** `GET`
- **Response:** `200 OK`

---

### Billing

#### Create Invoice
- **URL:** `/api/v1/billing/invoices`
- **Method:** `POST`
- **Body:**
  ```json
  {
    "invoice": {
      "member_id": "uuid",
      "total_amount": "35.00",
      "currency": "EUR",
      "due_date": "2026-02-01T00:00:00Z"
    },
    "items": [
      {
        "description": "Monthly Fee - Feb 2026",
        "amount": "35.00",
        "quantity": 1,
        "type": "SUBSCRIPTION"
      }
    ]
  }
  ```
- **Response:** `201 Created`

#### List Invoices
- **URL:** `/api/v1/billing/invoices`
- **Method:** `GET`
- **Query Params:** `member_id` (optional)
- **Response:** `200 OK`

#### Get Invoice
- **URL:** `/api/v1/billing/invoices/{id}`
- **Method:** `GET`
- **Response:** `200 OK`

#### Mark Invoice as Paid
- **URL:** `/api/v1/billing/invoices/{id}/pay`
- **Method:** `POST`
- **Response:** `200 OK`

#### Create Subscription
- **URL:** `/api/v1/billing/subscriptions`
- **Method:** `POST`
- **Body:**
  ```json
  {
    "member_id": "uuid",
    "modality_price_id": "uuid",
    "period": "MONTHLY",
    "start_date": "2026-01-01T00:00:00Z",
    "end_date": "2026-12-31T23:59:59Z"
  }
  ```
- **Response:** `201 Created`

#### List Subscriptions by Member
- **URL:** `/api/v1/billing/subscriptions/member/{memberID}`
- **Method:** `GET`
- **Response:** `200 OK`

---

### Enrollment

#### Enroll Member
- **URL:** `/api/v1/enrollments`
- **Method:** `POST`
- **Body:**
  ```json
  {
    "member_id": "uuid",
    "modality_id": "uuid",
    "season_id": "uuid"
  }
  ```
- **Response:** `201 Created`

#### List Enrollments by Member
- **URL:** `/api/v1/enrollments/member/{memberID}`
- **Method:** `GET`
- **Response:** `200 OK`

#### Approve Enrollment
- **URL:** `/api/v1/enrollments/{id}/approve`
- **Method:** `POST`
- **Response:** `200 OK`
