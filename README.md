## Folder Structure

```
backend/
├── cmd/
│ ├── synthema-api/
│ │ └── main.go
│ ├── synthema-capture/
│ │ └── main.go
│ └── synthema-worker/
│ └── main.go
│
├── internal/
│ ├── app/
│ │ ├── api/
│ │ │ ├── http/
│ │ │ │ ├── handlers/
│ │ │ │ ├── middlewares/
│ │ │ │ └── routes.go
│ │ │ └── server.go
│ │ │
│ │ ├── capture/
│ │ │ ├── interceptor.go
│ │ │ ├── encoder.go
│ │ │ └── capture_service.go
│ │ │
│ │ ├── replay/
│ │ │ ├── orchestrator.go
│ │ │ ├── scheduler.go
│ │ │ ├── session_graph.go
│ │ │ └── replay_service.go
│ │ │
│ │ ├── diff/
│ │ │ ├── comparator.go
│ │ │ ├── json_diff.go
│ │ │ └── diff_service.go
│ │ │
│ │ ├── transform/
│ │ │ ├── engine.go
│ │ │ ├── ruleset.go
│ │ │ └── sanitizer.go
│ │ │
│ │ └── health/
│ │ └── handler.go
│ │
│ ├── domain/
│ │ ├── traffic/
│ │ │ ├── request.go
│ │ │ ├── response.go
│ │ │ ├── session.go
│ │ │ └── metadata.go
│ │ │
│ │ ├── replay/
│ │ │ ├── replay_job.go
│ │ │ └── replay_result.go
│ │ │
│ │ ├── diff/
│ │ │ └── diff_result.go
│ │ │
│ │ └── common/
│ │ ├── id.go
│ │ ├── time.go
│ │ └── errors.go
│ │
│ ├── ports/
│ │ ├── repository/
│ │ │ ├── session_repository.go
│ │ │ ├── replay_repository.go
│ │ │ └── diff_repository.go
│ │ │
│ │ ├── queue/
│ │ │ └── traffic_queue.go
│ │ │
│ │ └── clock/
│ │ └── clock.go
│ │
│ ├── adapters/
│ │ ├── postgres/
│ │ │ ├── session_repo.go
│ │ │ ├── replay_repo.go
│ │ │ └── diff_repo.go
│ │ │
│ │ ├── redis/
│ │ │ ├── traffic_stream.go
│ │ │ └── ring_buffer.go
│ │ │
│ │ └── httpclient/
│ │ └── shadow_client.go
│ │
│ ├── config/
│ │ ├── config.go
│ │ └── loader.go
│ │
│ ├── bootstrap/
│ │ ├── dependencies.go
│ │ └── wiring.go
│ │
│ └── observability/
│ ├── logger.go
│ ├── metrics.go
│ └── tracing.go
│
├── migrations/
│ └── postgres/
│
├── pkg/
│ ├── backoff/
│ ├── ringbuffer/
│ └── jsonutil/
│
├── scripts/
│
├── go.mod
└── go.sum
```
