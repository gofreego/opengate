# OpenGate Config Management - Backend Refactoring Plan

This plan outlines the implementation of Config management features in the OpenGate service with PostgreSQL as the primary database.

## Overview
- **Feature**: Config CRUD operations (Create, GetByID, List with pagination, GetRoutes)
- **Database**: PostgreSQL (primary), other implementations marked as unimplemented
- **API**: gRPC with HTTP gateway (following catalogservice patterns)
- **Constraint**: Config `name` must be unique

---

## Tasks

### Phase 1: Database Schema

- [x] **1.1 Create PostgreSQL migration files** ✅
  - Create `resources/migrations/001_create_configs_table.up.sql`
  - Create `resources/migrations/001_create_configs_table.down.sql`
  - Table: `configs` with columns: id, name (unique), path_prefix, target_url, strip_prefix, authentication (jsonb), middleware (jsonb array), timeout, created_at, updated_at

---

### Phase 2: Proto Files & Code Generation

- [x] **2.1 Create proto directory structure** ✅
  - Create `api/proto/opengate/v1/` directory
  - Create `api/proto/common/` directory for shared types

- [x] **2.2 Create config.proto** ✅
  - Define `Config` message with all fields
  - Define `CreateConfigRequest`, `CreateConfigResponse`
  - Define `GetConfigRequest`, `GetConfigResponse`
  - Define `ListConfigsRequest`, `ListConfigsResponse` (with pagination)
  - Define `GetRoutesRequest`, `GetRoutesResponse`

- [x] **2.3 Create opengate.proto (service definition)** ✅
  - Define `OpenGateService` with gRPC methods
  - Add HTTP annotations for REST endpoints:
    - `POST /opengate/v1/configs` - CreateConfig
    - `GET /opengate/v1/configs/:id` - GetConfig
    - `GET /opengate/v1/configs` - ListConfigs
    - `GET /opengate/v1/routes` - GetRoutes

- [x] **2.4 Create ping.proto (common)** ✅
  - Copy/adapt from catalogservice for consistency

- [x] **2.5 Setup buf configuration** ✅
  - Create `api/buf.yaml`
  - Create `api/buf.gen.yaml`
  - Create `api/protoc.sh` generation script

- [x] **2.6 Generate Go code from proto** ✅
  - Run buf generate to create pb.go files

---

### Phase 3: Repository Layer

- [x] **3.1 Update Repository interface** ✅
  - Update `internal/service/service.go` Repository interface
  - Add methods: `CreateConfig`, `GetConfigByID`, `ListConfigs`, `GetRoutes`

- [x] **3.2 Implement PostgreSQL repository** ✅
  - Create `internal/repository/postgresql/config.go`
  - Create `internal/repository/postgresql/repository.go` with connection setup
  - Implement all Config CRUD methods
  - Handle unique name constraint with proper error handling

- [x] **3.3 Mark other repositories as unimplemented** ✅
  - Update `internal/repository/local/repository.go` - add unimplemented stubs
  - Update `internal/repository/openauth/repository.go` - add unimplemented stubs

- [x] **3.4 Update repository factory** ✅
  - Update `internal/repository/repository.go` to support PostgreSQL
  - Add PostgreSQL config struct and initialization

---

### Phase 4: Service Layer

- [x] **4.1 Create config service methods** ✅
  - Create `internal/service/config.go`
  - Implement `CreateConfig` with validation
  - Implement `GetConfigByID`
  - Implement `ListConfigs` with pagination
  - Implement `GetRoutes` (returns all routes for routing)

- [x] **4.2 Add business logic validation** ✅
  - Validate required fields
  - Validate unique name constraint (at service level)
  - Validate URL format for target_url

---

### Phase 5: gRPC/HTTP Handlers

- [x] **5.1 Create gRPC server implementation** ✅
  - Service struct embeds `UnimplementedOpenGateServiceServer`
  - Service directly implements gRPC methods (Ping, CreateConfig, etc.)
  - Following catalogservice pattern

- [x] **5.2 Update HTTP server** ✅
  - Updated `cmd/http_server/http.go` to combine gin + grpc-gateway
  - `/opengate/v1/*` routes use grpc-gateway mux
  - All other routes use gin router for proxy
  - Single port serves both API and proxy

---

### Phase 6: Configuration & Wiring

- [x] **6.1 Update dev.yaml configuration** ✅
  - Add PostgreSQL connection configuration
  - Add repository type selection (postgresql)

- [x] **6.2 Update main.go** ✅
  - Wire up PostgreSQL repository
  - Initialize config service

- [ ] **6.3 Add database migration runner**
  - Add migration execution on startup or as separate command

---

### Phase 7: Testing & Documentation

- [ ] **7.1 Add unit tests for repository layer**
  - Test PostgreSQL queries
  - Test unique constraint handling

- [ ] **7.2 Add unit tests for service layer**
  - Test business logic validation
  - Test pagination

- [ ] **7.3 Update README.md**
  - Document new config management API
  - Document PostgreSQL setup requirements

---

## File Structure (New/Modified)

```
opengate/
├── api/
│   ├── buf.yaml                          # NEW
│   ├── buf.gen.yaml                      # NEW
│   ├── protoc.sh                         # NEW
│   ├── opengate_v1/                      # NEW (generated)
│   │   ├── config.pb.go
│   │   ├── config.pb.validate.go
│   │   ├── opengate.pb.go
│   │   ├── opengate.pb.gw.go
│   │   ├── opengate_grpc.pb.go
│   │   └── ping.pb.go
│   └── proto/
│       ├── common/
│       │   └── ping.proto                # NEW
│       └── opengate/
│           └── v1/
│               ├── config.proto          # NEW
│               └── opengate.proto        # NEW
├── resources/
│   └── migrations/
│       ├── 001_create_configs_table.up.sql   # NEW
│       └── 001_create_configs_table.down.sql # NEW
├── internal/
│   ├── repository/
│   │   ├── postgresql/
│   │   │   ├── repository.go             # NEW
│   │   │   └── config.go                 # NEW
│   │   ├── local/
│   │   │   └── repository.go             # MODIFIED
│   │   └── repository.go                 # MODIFIED
│   └── service/
│       ├── config.go                     # NEW
│       └── service.go                    # MODIFIED
├── cmd/
│   ├── grpc_server/
│   │   └── grpc.go                       # NEW
│   └── http_server/
│       └── http.go                       # MODIFIED
├── dev.yaml                              # MODIFIED
└── main.go                               # MODIFIED
```

---

## Notes

- Following catalogservice patterns for consistency
- Using `buf` for proto compilation
- Proto validation using `protoc-gen-validate`
- gRPC-gateway for HTTP/REST support
- PostgreSQL with `lib/pq` driver (already in go.mod)
