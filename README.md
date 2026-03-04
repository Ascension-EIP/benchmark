> **Last updated:** 4th March 2026  
> **Version:** 1.0  
> **Authors:** Nicolas TORO  
> **Status:** Done  
> {.is-success}

---

# Benchmark

This repository holds all the prototype benchmarks produced during the Ascension tech-stack selection phase.
Each sub-project is a minimal, self-contained implementation used to validate a technology choice before committing to it in the main monorepo.

---

## Table of Contents

- [Benchmark](#benchmark)
  - [Table of Contents](#table-of-contents)
  - [Context](#context)
  - [Repository Structure](#repository-structure)
  - [Backend / Database Benchmarks](#backend--database-benchmarks)
    - [Rust + Axum + PostgreSQL](#rust--axum--postgresql)
    - [Go + Gin + MariaDB](#go--gin--mariadb)
    - [Backend Decision](#backend-decision)
  - [Mobile Benchmarks](#mobile-benchmarks)
    - [Flutter](#flutter)
    - [React Native (Expo)](#react-native-expo)
    - [Mobile Decision](#mobile-decision)
  - [What Each Benchmark Tests](#what-each-benchmark-tests)
  - [How to Run the Benchmarks](#how-to-run-the-benchmarks)
    - [Backend — Rust + PostgreSQL](#backend--rust--postgresql)
    - [Backend — Go + MariaDB](#backend--go--mariadb)
    - [Mobile — Flutter](#mobile--flutter)
    - [Mobile — React Native](#mobile--react-native)
  - [Next Steps](#next-steps)

---

## Context

Before writing a single line of production code, the team ran focused benchmarks on the two most debated stack choices:

1. **Backend + database**: which API framework and database combination handles our workload best?
2. **Mobile**: which cross-platform framework gives us the best camera/video experience on Android and iOS?

Each benchmark is intentionally minimal — it implements only the features that matter for the decision criterion (performance, DX, ecosystem fit, video pipeline), not a full application.

---

## Repository Structure

```
benchmark/
├── backend-db/
│   ├── rust-postgresql/   # Axum + SeaORM + PostgreSQL
│   └── go-mariadb/        # Gin + GORM + MariaDB
└── mobile/
    ├── flutter/           # Flutter + Dart
    └── react-native/      # React Native + Expo (TypeScript)
```

---

## Backend / Database Benchmarks

Both backends implement the same minimal feature set so results are directly comparable:

- `POST /register` — user sign-up with bcrypt password hashing
- `POST /login` — authentication returning a JWT token
- `POST /upload` — authenticated video upload endpoint (JWT-protected)

### Rust + Axum + PostgreSQL

| Item             | Detail                                       |
| ---------------- | -------------------------------------------- |
| Framework        | [Axum](https://github.com/tokio-rs/axum) 0.8 |
| ORM              | [SeaORM](https://www.sea-ql.org/SeaORM/) 1.x |
| Database         | PostgreSQL 16 (Docker)                       |
| Async runtime    | Tokio                                        |
| Password hashing | bcrypt                                       |
| Observability    | `tracing` + `tower-http` TraceLayer          |

The Rust implementation is the most bare-bones of the two: the HTTP layer, ORM, and database schema are all wired up in a single `main.rs` with no framework boilerplate. The database table is created at startup if it does not exist.

### Go + Gin + MariaDB

| Item             | Detail                                    |
| ---------------- | ----------------------------------------- |
| Framework        | [Gin](https://gin-gonic.com) 1.11         |
| ORM              | [GORM](https://gorm.io) + MariaDB driver  |
| Database         | MariaDB (via MySQL driver)                |
| Auth             | JWT (`golang-jwt/jwt` v5)                 |
| Password hashing | `golang.org/x/crypto/bcrypt`              |
| Architecture     | Layered: handler → service → repository   |

The Go implementation uses a full layered architecture (transport / service / repository / model) to evaluate how the language handles a realistic project structure, not just raw performance.

### Backend Decision

✅ **Chosen: Rust + Axum + PostgreSQL**

Key reasons:

- **Performance**: Rust's zero-cost abstractions and async Tokio runtime outperform Go/Gin in raw throughput and memory usage at every concurrency level tested.
- **Type safety**: SeaORM's compile-time checked queries eliminate an entire class of runtime errors that showed up during the Go prototype.
- **PostgreSQL over MariaDB**: better JSON support, full ACID compliance, and stronger ecosystem for the analytical queries the AI pipeline requires.
- **Consistency**: the AI worker already uses Python; adding Go would introduce a third language. Rust gives us systems-level control without the complexity of a third runtime.

---

## Mobile Benchmarks

Both mobile prototypes implement the same user-facing flow:

1. **Home screen** — entry point, navigation trigger.
2. **Camera screen** — request camera + microphone permissions, record a video (max 30 s).
3. **Preview screen** — play back the recorded video, display file size and duration, simulate an upload.

### Flutter

| Item           | Detail                     |
| -------------- | -------------------------- |
| Language       | Dart                       |
| Navigation     | `go_router` ^14            |
| Camera         | `camera` ^0.10.5           |
| Video playback | `video_player` ^2.8.2      |
| Permissions    | `permission_handler` ^11.3 |
| HTTP           | `http` ^1.2                |
| Min SDK        | Dart ≥ 3.0                 |

The Flutter prototype targets Android and iOS from a single codebase. The camera initialisation flow, permission handling, and video metadata extraction (path, size, duration) are all implemented from scratch using native Flutter packages.

### React Native (Expo)

| Item           | Detail                            |
| -------------- | --------------------------------- |
| Language       | TypeScript                        |
| Framework      | React Native 0.81 + Expo ~54      |
| Navigation     | React Navigation 7 (native stack) |
| Camera         | `expo-camera` ~17                 |
| Video playback | `expo-av` ~16 / `expo-video` ~3   |
| HTTP           | `axios` ^1.13                     |
| Target         | Android, iOS, Web (Expo Go)       |

The React Native prototype uses the Expo managed workflow for rapid iteration. The same three-screen flow (Home → Camera → Preview) is implemented in TypeScript with typed navigation.

### Mobile Decision

✅ **Chosen: Flutter**

Key reasons:

- **Camera reliability**: the Flutter `camera` package gave consistent behaviour across Android and iOS during testing. The Expo camera stack required more workarounds for permission flows and codec compatibility.
- **Performance**: Flutter's compiled Dart code rendered the camera preview and video playback noticeably more smoothly than the React Native JS bridge on the same physical devices.
- **Single language**: the team already uses Dart for Flutter; React Native would require maintaining TypeScript tooling separately.
- **Video pipeline fit**: Flutter's direct FFI access will be essential when integrating the AI video analysis pipeline. The React Native bridge adds latency that is unacceptable for real-time feedback.

---

## What Each Benchmark Tests

| Benchmark         | Key criteria evaluated                                                        |
| ----------------- | ----------------------------------------------------------------------------- |
| Rust + PostgreSQL | Throughput under concurrent load, type-safe ORM, Docker setup time            |
| Go + MariaDB      | Architecture scalability, JWT middleware, GORM ergonomics                     |
| Flutter           | Camera init reliability, permission UX, video metadata, cross-platform parity |
| React Native      | Expo DX, camera codec compatibility, navigation type safety, bundle size      |

---

## How to Run the Benchmarks

### Backend — Rust + PostgreSQL

```bash
cd benchmark/backend-db/rust-postgresql

# Start PostgreSQL
docker compose up -d

# Run the API (reads DATABASE_URL from .env)
cargo run
```

The server starts on `http://localhost:3000`.

Available routes:

```
GET  /           → health check
POST /register   → { "email": "...", "password": "..." }
POST /login      → { "email": "...", "password": "..." }
```

### Backend — Go + MariaDB

```bash
cd benchmark/backend-db/go-mariadb

# Copy and fill in the environment file
cp .env.example .env

# Run
go run ./cmd/app
```

### Mobile — Flutter

```bash
cd benchmark/mobile/flutter

flutter pub get
flutter run
```

Requires a physical device or emulator with camera support. Android and iOS are both supported.

### Mobile — React Native

```bash
cd benchmark/mobile/react-native

npm install
npx expo start
```

Scan the QR code with Expo Go, or press `a` / `i` to open the Android emulator / iOS simulator.

---

## Next Steps

The benchmarks informed the final stack for the production monorepo:

- **Backend**: Rust + Axum + PostgreSQL → `apps/server/`
- **Mobile**: Flutter + Dart → `apps/mobile/`

These benchmark sub-projects are **archived** — they are kept for reference and are not maintained or updated alongside the production code.
