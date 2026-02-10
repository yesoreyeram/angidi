# Angidi

Angidi is a high-performance, scalable **Enterprise E-commerce Platform** built with a modern **Microservices Architecture**. It is designed to handle high transaction volumes, support multi-regional operations (multi-currency and localization), and ensure best-in-class security compliance (PCI DSS, GDPR).

## Key Features

- **Identity & Access Management** — User registration, authentication (JWT), profile management, and role-based access control
- **Product Catalog** — Product CRUD, hierarchical categories, multi-currency pricing, and flexible product variants
- **Search & Discovery** — Full-text search, faceted filtering, sorting, and auto-completion powered by Elasticsearch
- **Product Reviews & Ratings** — Verified review submission, rating aggregation, moderation, and helpfulness voting

## Tech Stack

| Category | Technology |
| :--- | :--- |
| **Backend** | Go (Golang), Apache Kafka |
| **Frontend** | Next.js (React) |
| **Data Stores** | PostgreSQL, MongoDB, Elasticsearch, Redis |
| **Infrastructure** | Kubernetes, Docker |
| **Observability** | Prometheus, Grafana, Jaeger, Loki |

## Documentation

Detailed project documentation is available in the [`docs/`](docs/) directory:

- [Project Description](docs/PROJECT_DESCRIPTION.md)
- [Functional Requirements](docs/FUNCTIONAL_REQUIREMENTS.md)
- [Non-Functional Requirements](docs/NON_FUNCTIONAL_REQUIREMENTS.md)
- [Security Requirements](docs/SECURITY_REQUIREMENTS.md)
- [Observability Requirements](docs/OBSERVABILITY_REQUIREMENTS.md)
- [Modules & Services](docs/MODULES.md)
- [Tools & Tech Stack](docs/TOOLS_AND_TECH_STACK.md)
- [Delivery Plan](docs/DELIVERY_PLAN.md)
- [Design Decisions](docs/PROJECT_DESIGN_AND_DECISIONS.md)
- [Constitution](docs/CONSTITUION.md)

## License

This project is licensed under the [Apache License 2.0](LICENSE).
