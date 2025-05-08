# Time Record Service

[![Go Report Card](https://goreportcard.com/badge/github.com/patricksferraz/time-record-service)](https://goreportcard.com/report/github.com/patricksferraz/time-record-service)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://godoc.org/github.com/patricksferraz/time-record-service?status.svg)](https://godoc.org/github.com/patricksferraz/time-record-service)

A modern, scalable time record management service built with Go, featuring gRPC and REST APIs, event-driven architecture, and comprehensive monitoring.

## 🌟 Features

- **Dual API Support**: REST and gRPC endpoints for flexible integration
- **Event-Driven Architecture**: Kafka integration for reliable event processing
- **Database Management**: PostgreSQL with pgAdmin interface
- **Monitoring & Observability**: Elastic APM integration
- **Containerized**: Docker and Kubernetes support
- **Testing**: Comprehensive test suite with coverage reporting
- **Documentation**: Swagger/OpenAPI documentation

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.16 or later
- Make (optional, but recommended)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/patricksferraz/time-record-service.git
cd time-record-service
```

2. Copy the environment file and configure it:
```bash
cp .env.example .env
```

3. Start the services:
```bash
make up
```

The service will be available at:
- REST API: http://localhost:8080
- gRPC: localhost:50051
- pgAdmin: http://localhost:9000
- Kafka Control Center: http://localhost:9021

## 🛠️ Development

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

### Generating gRPC Code

```bash
make gen
```

### Viewing Logs

```bash
make logs
```

## 📚 API Documentation

The REST API documentation is available at `/swagger/index.html` when the service is running.

## 🏗️ Architecture

The service follows a clean architecture pattern with the following components:

- **Domain**: Core business logic and entities
- **Application**: Use cases and business rules
- **Infrastructure**: External services integration
- **Interface**: API handlers and controllers

## 🔄 Event Flow

The service processes the following events:
- `NEW_EMPLOYEE`
- `NEW_COMPANY`
- `NEW_TIME_RECORD`

## 📊 Monitoring

The service integrates with Elastic APM for:
- Performance monitoring
- Error tracking
- Distributed tracing
- Log aggregation

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Patrick Ferraz** - *Initial work*

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Confluent Kafka](https://github.com/confluentinc/confluent-kafka-go)
- [Elastic APM](https://www.elastic.co/apm)
- [GORM](https://gorm.io/)
