# URL Shortener Service

A highly scalable, low-latency URL Shortener designed for **read-heavy traffic**, supporting **custom aliases**, **expiry**, and **1B+ URLs** with **99.9% availability**.

---

## 🚀 Features

- Generate short URLs from long URLs
- Optional **custom alias** support
- Configurable **expiration time**
- Ultra-fast redirection (cache + CDN optimized)
- Horizontally scalable read & write services
- Collision-free short URL generation
- Designed for **read-heavy workloads**

---

## 📌 Functional Requirements

- Users can shorten a long URL
- Users can optionally:
    - Provide a custom alias
    - Set an expiration date (defaults if not provided)
- Users can retrieve, update, or delete a short URL
- Expired URLs should not redirect

---

## 📈 Non-Functional Requirements

- **Scalability**: Support up to **1B shortened URLs**
- **Availability**: 99.9%
- **Latency**: Minimal redirect latency
- **Reliability**: No two URLs map to the same short code
- **Performance**: Optimized for extremely high read-to-write ratio

---

## 🧠 Core Entities

| Entity        | Description |
|--------------|------------|
| Original URL | The long URL provided by the user |
| Short URL    | Generated or custom short code |
| User         | Creator of the short URL |

---

## 🔌 API Design

### Create Short URL
**POST** `/url`

```json
{
  "url": "https://google.com",
  "expiration_date": "date_in_UTC",
  "alias": "customAlias"
}
```

### Get Short URL
**GET** `/url/{short_url}`

## 🏗️ High Level Architecture

### Components

- **API Gateway**

- **Write Service**
    - Generates short codes
    - Persists URL mappings

- **Read Service**
    - Handles redirects
    - Reads from cache / DB

- **Redis**
    - Global counter
    - Read cache

- **Database**
    - Persistent storage

- **CDN**
    - Edge caching for hot URLs

---

## 🔁 URL Redirection Flow

1. Client accesses the short URL
2. **Read Service**:
    - Checks Redis cache for the short code
    - On cache miss → queries the database
    - Validates expiration time
3. Cache is updated on a cache miss
4. Returns **HTTP 302** with the `Location` header set to the original URL
