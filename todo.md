# UserAuth — Architecture Analysis & TODO

A Go user-authentication service built with **Clean Architecture** (layered / hexagonal style).
Stack: Gin (HTTP), GORM + Postgres, zap (logging), JWT (HS256), bcrypt, golang-migrate, Docker Compose.

## High-level layering

```
                          HTTP (Gin)
                              │
   cmd/main.go ──────────────────────────────────► wires everything up
        │
        ├── bootstrap/   (config, DB connection, migrations)   ← infrastructure setup
        ├── container/   (dependency injection root)
        │
        ▼
┌─────────────────────────────────────────────────────────────┐
│  delivery (HTTP handlers)   internal/api/delivery, route     │  ← outer
│        depends on ▼                                          │
│  usecase (business logic)   internal/usecase                 │
│        depends on ▼                                          │
│  domain (entities + repo interfaces)   domain/               │  ← inner / core
│        ▲ implemented by                                      │
│  repository (GORM/Postgres)  internal/repository             │  ← outer
└─────────────────────────────────────────────────────────────┘
        pkg/  (jwt, bcrypt, logger)  — reusable, framework-agnostic helpers
```

The **dependency rule** is respected: dependencies point inward. `domain` knows nothing about
Gin, GORM, or Postgres — it only defines entities and repository *interfaces*. Outer layers
depend on those interfaces.

## Layer by layer

- **`domain/`** — the core. Pure Go entities (`User`, `Role`, `Permission`,
  `PersonalAccessToken`) plus repository interfaces (`UserRepository`, `RoleRepository`, etc.)
  and domain errors. A shared `Base` struct carries `ID uuid.UUID`, `CreatedAt`, `UpdatedAt`.
- **`internal/usecase/`** — application/business logic. `UserUseCase` orchestrates signup
  (hash password → create user → ensure "user" role exists → assign role) and login
  (lookup → bcrypt compare → issue JWT). Depends only on domain interfaces + `pkg/bcrypt` +
  `pkg/jwt`. Each method wraps a `context.WithTimeout`. Inputs arrive as `*Input` structs.
- **`internal/repository/`** — persistence adapters implementing the domain interfaces with
  GORM/Postgres. Separate **persistence models** (`*_model.go`, GORM tags + `TableName()`) are
  kept distinct from domain entities, with `ToDomain()` / `FromDomain()` mappers. Many-to-many
  relations (`user_has_roles`, `user_has_permissions`, `role_has_permission`) and cursor-based
  pagination (`Helper.go`) live here.
- **`internal/api/delivery/`** — HTTP controllers (Gin). Bind DTOs → call usecase → map to HTTP
  responses. `route/route.go` is currently commented out; routes are registered in `main.go`.
- **`internal/dtos/`** — request/response shapes with Gin validation tags, decoupled from both
  domain and persistence models.
- **`internal/bootstrap/`** — infrastructure: `LoadConfig` (env via godotenv with fallbacks),
  `Connect` (GORM Postgres + pool tuning), migrations.
- **`internal/container/`** — composition root / manual DI. `NewContainer` instantiates
  repos → usecases → deliveries.
- **`pkg/`** — framework-agnostic utilities: `jwt` (HS256 access/refresh), `bcrypt`, `logger` (zap).

## Data / infra

- **Postgres** with versioned SQL migrations in `migration/` (users, PATs, roles, permissions,
  and three join tables) — an RBAC model.
- **Docker Compose**: `app`, `db` (postgres:16, healthchecked), `pgadmin`, one-shot `migrate`.

## Bugs / TODO

- [ ] **`main.go` bypasses the route layer** — `route.go` is fully commented out; signup/login
  are wired inline in `main.go`. Only `UserUseCase` is wired in the container; role/permission/PAT
  usecases exist but aren't exposed.
- [ ] **Config bug** — in `bootstrap/config.go` `LoadConfig`, the `.env` loop sets `loadedFrom`
  only when `godotenv.Load` *errors* (`if err != nil` is inverted), so the log message is backwards.
- [ ] **Two `BeforeCreate` hooks** — both `domain.Base` and `repository.User` define a GORM
  `BeforeCreate`. GORM operates on the persistence model, so the domain one is effectively dead code.
- [ ] **Login error mapping** — invalid credentials / user-not-found both return `500` instead of
  `401`; `Login` also returns `http.StatusBadGateway` (502) for a bad request body.
- [ ] **Review `PRODUCTION_REVIEW.md`** — already in the repo; likely overlaps with some of the above.

## Clean Architecture review (dependency-rule violations & inconsistencies)

- [ ] **`roleUseCase` depends on the concrete repository** — `usecase/role.go` imports
  `internal/repository` and holds a `*repository.RoleRepository`, pointing an inner layer at an
  outer one. The `domain.RoleRepository` interface already exists (and `UserUseCase` uses the
  interface form correctly) — switch `roleUseCase` to depend on `domain.RoleRepository`.
- [ ] **`domain` depends on GORM** — `domain/user.go` imports `gorm.io/gorm` and defines a
  `BeforeCreate(tx *gorm.DB)` hook on the entity. The innermost layer must be framework-free;
  remove the hook + import (it's dead code anyway — see duplicate `BeforeCreate` above).
- [ ] **Use cases hard-coupled to concrete crypto/JWT** — `usecase/user.go` calls `bcrypt` and
  `jwt` package functions directly. Abstract them behind ports (interfaces) so they can be
  mocked/swapped and the use cases become unit-testable in isolation.
- [ ] **`UserUseCase` has no interface; `RoleUseCase` does** — `RoleDelivery` depends on the
  `usecase.RoleUseCase` interface, but `UserDelivary` depends on concrete `*usecase.UserUseCase`.
  Pick one convention — add a `UserUseCase` interface so delivery depends on abstractions everywhere.
- [ ] **Repos return `(nil, nil)` on not-found** — `repository/role.go` `FindByID`/`FindByName`
  return `nil, nil`, so callers (`usecase/user.go` signup default-role branch, `usecase/role.go`
  `FindByID`) dereference a nil pointer. Return `domain.Err*NotFound` like the user repo does.
- [ ] **`FindByID` swallows the error** — `usecase/user.go` `FindByID` returns `nil` instead of
  `err` on repository failure.
- [ ] **Preload casing** — `repository/user.go` `FindByID` uses `.Preload("roles")` (lowercase)
  but the association is `Roles`; roles also aren't mapped in `ToDomain`, so they never load.

