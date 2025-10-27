# Application Modules and Service Definition

## 1. Core Services (Foundation for Initial Phase)

These services are mandatory for the initial platform launch, covering User, Catalog, Search, and Reviews.

### 1.1. Identity & Access Management (IAM) Service 🔑

| Attribute            | Detail                                                                    |
| :------------------- | :------------------------------------------------------------------------ |
| **Technology**       | Golang Microservice                                                       |
| **Data Store**       | PostgreSQL                                                                |
| **Responsibilities** | - User Registration (FR-U-01) and Authentication (Login).                 |
|                      | - Token generation (JWT) and validation.                                  |
|                      | - User Profile and Address Management (FR-U-03).                          |
|                      | - Centralized **Role-Based Access Control (RBAC)** enforcement (FR-U-04). |
|                      | - Password recovery and reset (FR-U-02).                                  |

### 1.2. Product Catalog Service 📦

| Attribute            | Detail                                                                                      |
| :------------------- | :------------------------------------------------------------------------------------------ |
| **Technology**       | Golang Microservice                                                                         |
| **Data Store**       | PostgreSQL (Core Data), MongoDB (Attributes/Variants)                                       |
| **Responsibilities** | - **Product CRUD** operations for Content Managers (FR-P-01).                               |
|                      | - Managing Product Core Data (Name, Base Price, SKU) in **PostgreSQL**.                     |
|                      | - Managing flexible Product Attributes and Variants (Color, Size) in **MongoDB** (FR-P-03). |
|                      | - Managing **Hierarchical Categories** (FR-P-02).                                           |
|                      | - Publishing product changes to the Search Service (via Kafka) for indexing.                |

### 1.3. Search & Discovery Service 🔍

| Attribute            | Detail                                                                                       |
| :------------------- | :------------------------------------------------------------------------------------------- |
| **Technology**       | Golang Microservice                                                                          |
| **Data Store**       | Elasticsearch (Primary Index), Redis (Caching)                                               |
| **Responsibilities** | - Handling all user search queries (FR-S-01).                                                |
|                      | - Providing **Faceted Search** and filtering capabilities (FR-S-02).                         |
|                      | - Providing **Auto-Completion** suggestions (FR-S-04).                                       |
|                      | - Maintaining and updating the Elasticsearch index based on events from the Catalog Service. |

### 1.4. Product Review & Rating Service ⭐

| Attribute            | Detail                                                                                        |
| :------------------- | :-------------------------------------------------------------------------------------------- |
| **Technology**       | Golang Microservice                                                                           |
| **Data Store**       | MongoDB (Review Content), PostgreSQL (Vote Tracking)                                          |
| **Responsibilities** | - Handling submission and storage of user reviews and star ratings (FR-R-01).                 |
|                      | - Managing Review **Moderation** queues for Moderators (FR-R-03).                             |
|                      | - Tracking user votes on review helpfulness (FR-R-04) using PostgreSQL to enforce uniqueness. |
|                      | - Publishing events (via Kafka) when a rating is submitted/approved.                          |

### 1.5. Frontend Application 🌐

| Attribute            | Detail                                                                                             |
| :------------------- | :------------------------------------------------------------------------------------------------- |
| **Technology**       | Next.js (React)                                                                                    |
| **Responsibilities** | - User Interface and User Experience (UI/UX).                                                      |
|                      | - Server-Side Rendering (SSR) for optimal performance and SEO.                                     |
|                      | - Communication with all backend microservices via the API Gateway.                                |
|                      | - Handling client-side routing, state management, and interaction with Search and Catalog modules. |

---
