# Tools and Technology Stack Specification

## 1. Programming Languages & Core Frameworks 💻

| Category               | Technology                           | Version / Rationale                                                                                              |
| :--------------------- | :----------------------------------- | :--------------------------------------------------------------------------------------------------------------- |
| **Backend Core**       | **Go (Golang)**                      | Latest stable version. Mandatory for performance and microservices (Constitution.md).                            |
| **Frontend Core**      | **React**                            | Latest stable version. Core library for UI development.                                                          |
| **Frontend Framework** | **Next.js**                          | Mandatory for Server-Side Rendering (SSR) and routing optimization.                                              |
| **Styling**            | **Tailwind CSS / Styled Components** | TBD based on team preference, but one standard must be adopted globally.                                         |
| **Asynchronous**       | **Node.js**                          | Allowed for peripheral tooling, scripting, and build processes, but **NOT** for core application business logic. |

---

## 2. Data Stores & Databases 💾

| Category                     | Technology        | Rationale / Usage                                                             |
| :--------------------------- | :---------------- | :---------------------------------------------------------------------------- |
| **Primary Transactional DB** | **PostgreSQL**    | Used for all core transactional data (Users, Orders, Prices, Coupons, Votes). |
| **Flexible Data Store**      | **MongoDB**       | Used for flexible schema data (Product Attributes, Review Content).           |
| **Full-Text Search Engine**  | **Elasticsearch** | Mandatory for the Search & Discovery Service indexing and querying.           |
| **In-Memory Cache / Broker** | **Redis**         | High-speed caching, session management, and rate limiting.                    |

---

## 3. Communication & Messaging 🔗

| Category                 | Technology              | Rationale / Usage                                                                                                    |
| :----------------------- | :---------------------- | :------------------------------------------------------------------------------------------------------------------- |
| **API Protocol**         | **RESTful APIs (JSON)** | Standard for external client-server communication and most internal sync communication.                              |
| **Async Message Broker** | **Apache Kafka**        | Mandatory for reliable, high-throughput asynchronous communication (e.g., rating updates, search index propagation). |
| **API Specification**    | **OpenAPI (Swagger)**   | Required for documenting all microservice endpoints (Internal and External).                                         |

---

## 4. DevOps & Cloud Infrastructure ☁️

| Category                         | Technology                                     | Rationale / Usage                                                                   |
| :------------------------------- | :--------------------------------------------- | :---------------------------------------------------------------------------------- |
| **Containerization**             | **Docker**                                     | Mandatory for packaging all microservices and ensuring environment parity.          |
| **Orchestration**                | **Kubernetes (K8s)**                           | Mandatory for deployment, automated scaling, and cluster management.                |
| **CI/CD Pipeline**               | **GitHub Actions / GitLab CI / Jenkins (TBD)** | Automated testing, security scanning (SAST/SCA), and rolling deployments.           |
| **Infrastructure-as-Code (IaC)** | **Terraform (Recommended)**                    | Used to manage and provision cloud resources and Kubernetes cluster infrastructure. |

---

## 5. Observability & Monitoring 📈

| Category                     | Technology               | Usage (Refer to OBSERVABILITY_REQUIREMENTS.md)                        |
| :--------------------------- | :----------------------- | :-------------------------------------------------------------------- |
| **Metrics Collector**        | **Prometheus**           | Used to scrape and store time-series metrics from all services.       |
| **Visualization & Alerting** | **Grafana**              | Used for visualizing metrics and configuring operational alerts.      |
| **Distributed Tracing**      | **Jaeger (Recommended)** | Used for implementing the distributed tracing standard via Trace IDs. |
| **Log Aggregation**          | **Loki (Recommended)**   | Cost-effective, scalable storage and querying for structured logs.    |

---

## 6. Security & Quality Tooling 🛡️

| Category                     | Technology                  | Usage (Refer to SECURITY_REQUIREMENTS.md)                                |
| :--------------------------- | :-------------------------- | :----------------------------------------------------------------------- |
| **SAST (Static Analysis)**   | **GoSec (for Golang)**      | Automated analysis of code for security vulnerabilities.                 |
| **SCA (Component Analysis)** | **Dependabot / Snyk (TBD)** | Detects known vulnerabilities in third-party libraries and dependencies. |
| **Code Formatting**          | **Gofmt / Prettier**        | Enforces consistent code style across the entire codebase.               |
