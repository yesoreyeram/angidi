# Non-Functional Requirements Specification (NFRD)

## 1. Performance & Scalability (P&S)

| ID           | Requirement                 | Metric / Constraint                                                                                                                                  |
| :----------- | :-------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------- |
| **NFR-P-01** | **API Latency (Core)**      | Core transactional APIs (Login, Cart Update, Product Details Fetch) **MUST** respond within **500 milliseconds (P95)**.                              |
| **NFR-P-02** | **Page Load Time**          | All major public-facing pages (Homepage, Product Listing, Search Results) **MUST** load within **2.0 seconds** over a standard broadband connection. |
| **NFR-P-03** | **Concurrent Load**         | The system **MUST** be capable of handling **10,000 concurrent users** browsing and transacting without service degradation.                         |
| **NFR-P-04** | **Horizontal Scaling**      | All microservices **MUST** be horizontally scalable to meet peak demand spikes (e.g., during major sales events).                                    |
| **NFR-P-05** | **Asynchronous Processing** | Asynchronous jobs (e.g., rating aggregation via Kafka) **MUST** complete processing within **5 minutes** of the triggering event.                    |

---

## 2. Reliability & Availability (R&A)

| ID           | Requirement           | Metric / Constraint                                                                                                                                                       |
| :----------- | :-------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **NFR-R-01** | **System Uptime**     | The overall platform **MUST** target an annual average availability of **99.9%** (less than 8.76 hours of unplanned downtime per year).                                   |
| **NFR-R-02** | **Fault Tolerance**   | The architecture **MUST** be resilient to the failure of any single microservice instance (no single point of failure), managed via Kubernetes.                           |
| **NFR-R-03** | **Data Consistency**  | Core transactional data (Orders, Inventory, Payments) **MUST** maintain **strong consistency** (ACID).                                                                    |
| **NFR-R-04** | **Disaster Recovery** | A Disaster Recovery Plan **MUST** be established to ensure a **Recovery Time Objective (RTO)** of **< 4 hours** and a **Recovery Point Objective (RPO)** of **< 1 hour**. |

---

## 3. Security (SEC)

| ID           | Requirement                   | Metric / Constraint                                                                                                                                                                          |
| :----------- | :---------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **NFR-S-01** | **Compliance (Payments)**     | The system's payment processing components **MUST** be compliant with the latest **PCI DSS** standards.                                                                                      |
| **NFR-S-02** | **Data Encryption (Rest)**    | All sensitive Personally Identifiable Information (PII) and payment tokens **MUST** be encrypted **at rest** using industry-standard algorithms (e.g., AES-256).                             |
| **NFR-S-03** | **Data Encryption (Transit)** | All external and inter-service communication **MUST** be secured using **TLS 1.2 or higher**.                                                                                                |
| **NFR-S-04** | **Input Validation**          | All user and Content Manager inputs **MUST** be strictly validated and sanitized to prevent **OWASP Top 10** vulnerabilities, including SQL Injection and XSS.                               |
| **NFR-S-05** | **SSRF Mitigation**           | Strict validation and whitelisting **MUST** be applied to any user-provided URLs (e.g., image embeds in reviews, data import sources) to prevent Server-Side Request Forgery (SSRF) attacks. |
| **NFR-S-06** | **Authentication Security**   | Password hashes **MUST** be stored using a strong, salted, adaptive function (e.g., Argon2 or bcrypt).                                                                                       |

---

## 4. Maintainability & Technology (M&T)

| ID           | Requirement                    | Metric / Constraint                                                                                                                                                                      |
| :----------- | :----------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **NFR-M-01** | **Technology Stack Adherence** | All core services **MUST** utilize **Golang** (backend) and **Next.js (React)** (frontend) running on **Kubernetes**.                                                                    |
| **NFR-M-02** | **Code Quality**               | A minimum of **80% unit test coverage** is required for all business logic in the Golang backend services.                                                                               |
| **NFR-M-03** | **API Documentation**          | All public and internal APIs **MUST** be documented using the **OpenAPI (Swagger)** specification.                                                                                       |
| **NFR-M-04** | **Observability**              | All microservices **MUST** implement comprehensive observability through: **Structured Logging** (Loki), **Metrics** (Prometheus/Grafana), and **Distributed Tracing** (Jaeger).         |
| **NFR-M-05** | **Tracing Standard**           | Every incoming request **MUST** generate a unique **Trace ID** at the frontend or API Gateway, which is propagated across all downstream microservices to monitor the full user journey. |

---

## 5. Portability & Localization (P&L)

| ID           | Requirement                     | Metric / Constraint                                                                                                                                 |
| :----------- | :------------------------------ | :-------------------------------------------------------------------------------------------------------------------------------------------------- |
| **NFR-L-01** | **Internationalization (i18n)** | The application architecture **MUST** be designed to support multiple languages and character sets with minimal code changes.                       |
| **NFR-L-02** | **Multi-Currency Display**      | The frontend **MUST** correctly display prices, taxes, and totals in the user's selected or inferred local currency (e.g., Dollars, Euros, Rupees). |
| **NFR-L-03** | **Localization of Data**        | The database design **MUST** accommodate country-specific address formats, tax rules, and regional time zones.                                      |
