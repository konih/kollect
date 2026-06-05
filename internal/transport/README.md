# transport

Lean publish/subscribe boundary for inventory change notifications inside the operator.

## Backends

| Type | Implementation | Integration test |
| --- | --- | --- |
| `inprocess` | `InProcessBus` (default) | unit |
| `redis` | Redis Streams | testcontainers `modules/redis` |
| `nats` | NATS JetStream | testcontainers `modules/nats` |
| `kafka` | Kafka/Redpanda | testcontainers `modules/redpanda` |

Configure via `transport.Config` or `ConfigFromEnv()` (`KOLLECT_TRANSPORT_TYPE`, backend URLs).
Optional TLS for Redis/NATS: `KOLLECT_TRANSPORT_TLS_CA_FILE`, `KOLLECT_TRANSPORT_TLS_CLIENT_CERT_FILE`,
`KOLLECT_TRANSPORT_TLS_CLIENT_KEY_FILE`, `KOLLECT_TRANSPORT_TLS_INSECURE_SKIP_VERIFY` (ADR-0028).

Use cases: spoke → hub inventory reports (ADR-0022), debounced export triggers, and optional
decoupling of collection from export workers.
