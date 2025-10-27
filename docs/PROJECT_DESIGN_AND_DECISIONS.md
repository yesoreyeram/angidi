# Project Design and Key Decisions

## 1. Architectural Style Decision

**Decision:** Adopt a **Microservices Architecture** over a monolithic approach.

| Rationale                  | Link to Requirement                                                                                                                                                  |
| :------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Scalability**            | Allows individual, high-load services (e.g., Search, Inventory) to scale horizontally and independently from less-used services (e.g., Admin Panel). (**NFR-P-04**). |
| **Resilience**             | Failure in one service (e.g., Reviews) does not bring down core functionality (e.g., Cart/Checkout) (**NFR-R-02**).                                                  |
| **Technology Flexibility** | Enables the selection of the best-fit technology for each task (e.g., Golang for performance, MongoDB for flexibility) (**NFR-M-01**).                               |
| **Deployment**             | Facilitates independent, rapid deployment cycles using Kubernetes (DEP-CD-04).                                                                                       |

---

## 2. Core Technology Stack Rationale

The technology stack was chosen to meet the high-performance and scalability requirements of an enterprise e-commerce platform.

| Component              | Technology               | Decision Rationale                                                                                                                                                |
| :--------------------- | :----------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Backend Language**   | **Golang**               | Selected for its superior concurrency model, low memory footprint, and CPU efficiency, directly addressing the **500ms P95 latency** target (**NFR-P-01**).       |
| **Frontend Framework** | **Next.js (React)**      | Selected to enable **Server-Side Rendering (SSR)**, which is crucial for SEO and achieving a fast initial page load time (**NFR-P-02**).                          |
| **Data Architecture**  | **Polyglot Persistence** | Utilizes **PostgreSQL** for transactional integrity (Orders, Payments) and **MongoDB** for flexible, non-transactional data (Product Attributes, Review Content). |
| **Async Messaging**    | **Apache Kafka**         | Mandatory for decoupling services and handling high-volume, non-critical updates (e.g., product updates, rating recalculations) asynchronously (**NFR-P-05**).    |
| **Deployment**         | **Kubernetes (K8s)**     | Standardized orchestration platform for containerized applications, enabling automatic self-healing and rolling updates (DEP-ENV-01).                             |

---

## 3. Key Data Model and Service Decisions

### 3.1. Data Consistency Strategy

| Decision                 | Rationale                                                                                                                                                                                  |
| :----------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Strong Consistency**   | Used for critical financial and inventory data (PostgreSQL). We cannot allow for eventual consistency in the checkout process.                                                             |
| **Eventual Consistency** | Used for user-facing, non-critical content (Search Index, Product Ratings). Updates are propagated via **Kafka** to prioritize high availability and performance over instant consistency. |

### 3.2. Product Catalog Data Split

| Data Type                 | Data Store     | Rationale                                                                                                                                                                                 |
| :------------------------ | :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Core Product Data**     | **PostgreSQL** | Stores SKU, name, core price, and Category IDs—data that requires transactional guarantees and strong foreign key relationships.                                                          |
| **Attributes & Variants** | **MongoDB**    | Stores flexible data like size charts, material specifications, and product variants. Allows Content Managers to quickly introduce new attribute types without schema migration downtime. |

### 3.3. Security Mitigation (SSRF Example)

| Vulnerability                          | Mitigation Strategy                                                                                                                                                                                                                          | Linked Requirement |
| :------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :----------------- |
| **Server-Side Request Forgery (SSRF)** | Any API endpoint that accepts a user- or admin-supplied URL (e.g., image upload, data import) **MUST** have a dedicated validation layer that performs **strict whitelisting** of allowed domains and blocks all private/internal IP ranges. | **SEC-V-03**       |

---

## 4. Observability and Quality Decision

**Decision:** Full implementation of the **Three Pillars of Observability** is mandatory from Phase 1.

| Pillar      | Implementation Standard                                                                                                                             | Impact on Development                                                                                                       |
| :---------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------------------------------------- |
| **Tracing** | **Distributed Tracing (Jaeger)**. A **Trace ID** must be generated at the Next.js frontend/API Gateway and propagated to every Golang microservice. | Developers must ensure the Trace ID is correctly passed via context in all internal API calls.                              |
| **Metrics** | **Prometheus/Grafana**. Standard operational and key business metrics (e.g., `Coupon_Applied_Success`) must be exposed by every service.            | Teams must define and implement custom metrics alongside feature development.                                               |
| **Quality** | **80% Unit Test Coverage (Golang)**. Enforced via the CI/CD pipeline (DEP-CI-01).                                                                   | Mandates a test-driven or highly tested development approach to ensure code reliability and maintainability (**NFR-M-02**). |
