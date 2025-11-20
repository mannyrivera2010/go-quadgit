# go-quadgit
Excellent suggestion. A more detailed README, especially for a complex project, can significantly improve adoption by answering a potential user's questions before they even ask them. It shows maturity and a commitment to user experience.

This version expands on the previous one by adding more technical depth, more concrete examples, and a "Why It's Different" section to directly address how it compares to other tools.

---

# Graph Git (`quadgit`)

<p align="center">
  <a href="https://goreportcard.com/report/github.com/your-user/quadgit"><img src="https://goreportcard.com/badge/github.com/your-user/quadgit" alt="Go Report Card"></a>
  <a href="https://github.com/your-user/quadgit/actions/workflows/ci.yml"><img src="https://github.com/your-user/quadgit/actions/workflows/ci.yml/badge.svg" alt="CI Status"></a>
  <a href="https://github.com/your-user/quadgit/blob/main/LICENSE"><img src="https://img.shields.io/badge/License-Apache--2.0-blue.svg" alt="License"></a>
  <a href="https://pkg.go.dev/github.com/your-user/quadgit/pkg/quadstore"><img src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white" alt="Go Reference"></a>
  <a href="https://discord.gg/your-invite-code"><img src="https://img.shields.io/discord/your-server-id?logo=discord&label=discord" alt="Discord"></a>
</p>

<p align="center">
  <strong>The Auditable Knowledge Graph. Version your RDF data like you version your code.</strong>
</p>

**Graph Git (`quadgit`)** is a high-performance, distributed database for RDF knowledge graphs that brings the power, safety, and familiarity of Git's branching, merging, and history tracking to the world of linked data. It is designed from the ground up to provide a cryptographically secure, fully auditable history of your data's evolution.

Stop treating your knowledge graph as a black box that only knows its current state. With `quadgit`, you can ask not just "what is true now?", but also "what was true last year?", "who changed this fact?", and "what was the exact state of our knowledge when we ran that analysis?".

This repository contains the source for:
*   **`quadgit`**: A powerful CLI for creating, versioning, and inspecting graph repositories.
*   **`quadstore`**: The underlying Go library designed for high performance and embedding in your own applications.
*   **`quadgit-server` (Future)**: A planned REST server for multi-user collaboration, security, and higher-level workflows.

---

## Table of Contents

- [Why Graph Git? The Problem It Solves](#why-graph-git-the-problem-it-solves)
- [How Is It Different?](#how-is-it-different)
- [Core Features at a Glance](#core-features-at-a-glance)
- [Core Concepts: The Git Analogy](#core-concepts-the-git-analogy)
- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Quick Start: A 5-Minute Tour](#quick-start-a-5-minute-tour)
- [Architecture Overview](#architecture-overview)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

## Why Graph Git? The Problem It Solves

Go to problem_solutions

## How Is It Different?

| Tool | Focus | Versioning Model | Our Advantage |
| :--- | :--- | :--- | :--- |
| **`quadgit`** | **Version Control & Provenance** | **Full Git Model:** Decentralized, cryptographic history with branching and merging. | The only tool that treats data history with the same rigor as source code. Built for audit and collaboration. |
| **Standard Graph DBs**<br>(e.g., Neo4j, JanusGraph) | **Real-time Query Performance** | **None (or limited temporal features).** They primarily store the current state of the graph. | While they may be faster for certain complex graph traversals on live data, they cannot answer "who changed what, and when?" |
| **Git for Text Files**<br>(e.g., storing Turtle files in Git) | **Text Versioning** | **Line-based diffs.** It doesn't understand the graph structure. | `quadgit` understands the data. A `diff` shows semantic quad changes, and `merge` can detect logical conflicts, not just textual ones. |
| **Temporal Databases** | **Time-based Queries** | **Bitemporal model** (valid time vs. transaction time). | Temporal DBs are powerful but complex. `quadgit` provides a more intuitive, developer-friendly "commit-based" versioning model that aligns with existing DevOps workflows. |

## Core Features at a Glance

- **Full Git-like Versioning:** `commit`, `branch`, `merge`, `tag`, `log`, `diff`, `blame`.
- **Cryptographic Integrity:** Content-addressing (SHA-1) for all data and history objects.
- **GPG Signing:** Commits can be signed to prove authorship and prevent tampering.
- **Data Provenance:** Use `quadgit blame` to trace any quad back to its origin commit.
- **Scalable by Design:** Built on BadgerDB with stream-based algorithms to handle graphs with billions of quads.
- **High Performance:** A multi-instance database architecture tuned for different read/write patterns.
- **Multi-Tenancy:** Manage isolated knowledge graphs with `namespaces`.
- **Schema-Aware Merging:** Prevents merges that would introduce logical inconsistencies based on your RDFS/OWL ontology.
- **Distributed:** `push`, `pull`, and `fetch` to synchronize repositories.
- **Extensible:** A hook system (`pre-commit`, `post-commit`) allows you to trigger custom validation or notification scripts.

## Core Concepts: The Git Analogy

| Git Concept | `quadgit` Analogy |
| :--- | :--- |
| File | A **Named Graph** |
| Line of Code | An **RDF Quad** |
| Commit | An immutable, content-addressable **snapshot** of the entire database state. |
| Branch | A lightweight pointer to a line of development. |
| Merge | A semantically-aware integration of knowledge from two different branches. |

## Getting Started

### Installation

*From Source (requires Go 1.18+):*
```bash
git clone https://github.com/your-user/quadgit.git
cd quadgit
# Installs the binary to your $GOPATH/bin
go install ./cmd/quadgit
```

*Using a Pre-compiled Binary:*
Download the latest release for your OS and architecture from the [Releases page](https://github.com/your-user/quadgit/releases).

### Quick Start: A 5-Minute Tour

1.  **Initialize a new repository:**
    ```bash
    mkdir financial-data
    cd financial-data
    quadgit init
    # -> Initialized empty Graph Git repository in .quadgit/
    ```

2.  **Create your first data file (`ownership_v1.nq`):**
    ```n-quads
    <ex:company_A> <ex:owns> <ex:company_B> <ex:ownership_data> .
    ```

3.  **Add and commit the initial state:**
    ```bash
    quadgit add ownership_v1.nq
    quadgit commit -m "Feat: Record initial ownership of Company B"
    ```

4.  **Create a branch to model a proposed acquisition:**
    ```bash
    quadgit branch feature/acme_acquisition
    quadgit checkout feature/acme_acquisition
    ```

5.  **Create a new file (`acquisition.nq`) representing the change:**
    *   `D` prefix marks a quad for deletion.
    *   Lines without a prefix are additions.
    ```n-quads
    D <ex:company_A> <ex:owns> <ex:company_B> <ex:ownership_data> .
    <ex:acme_corp> <ex:owns> <ex:company_B> <ex:ownership_data> .
    ```

6.  **Commit the change on your feature branch:**
    ```bash
    quadgit add acquisition.nq
    quadgit commit -m "Feat: Model acquisition of Company B by Acme Corp"
    ```

7.  **See the difference:**
    ```bash
    quadgit diff main
    # Output:
    # - <ex:company_A> <ex:owns> <ex:company_B> <ex:ownership_data> .
    # + <ex:acme_corp> <ex:owns> <ex:company_B> <ex:ownership_data> .
    ```

8.  **Find out who is responsible for the current state:**
    ```bash
    quadgit blame ex:ownership_data
    # Output:
    # (a1b2c3d Alice) <ex:acme_corp> <ex:owns> <ex:company_B> <ex:ownership_data> .
    ```

## Architecture Overview

`quadgit` is not a monolithic application but a layered system designed for clean separation of concerns, testability, and performance. Understanding this architecture is key to understanding the project's capabilities and how to contribute to it.

The architecture flows from a high-performance key-value store at the bottom to multiple user-facing applications at the top. Each layer has a distinct responsibility and communicates through a well-defined API.

```
+-------------------------------------------------+
|               Presentation Layer                |
|                                                 |
|  +-----------------+     +--------------------+ |
|  |  quadgit CLI    |     | quadgit-server     | |
|  | (Cobra App)     |     | (REST API)         | |
|  +-----------------+     +--------------------+ |
+-----------------|-------------------------------|
                  |
                  | Both consume the public API...
                  v
+-------------------------------------------------+
|                 Core API Layer                  |
|                                                 |
|         `quadstore.Store` Go Interface          |
|      (The "Headless" Library / Public API)      |
|                                                 |
+-----------------|-------------------------------|
                  |
                  | Implemented by...
                  v
+-------------------------------------------------+
|                Datastore Layer                  |
|                                                 |
| +-------------+  +------------+  +------------+ |
| | history.db  |  |  index.db  |  |   app.db   | |
| | (BadgerDB)  |  | (BadgerDB) |  | (BadgerDB) | |
| +-------------+  +------------+  +------------+ |
+-------------------------------------------------+
```

### 1. The Datastore Layer: A Tuned Persistence Engine

The foundation of `quadgit` is **BadgerDB**, a high-performance, embeddable, key-value store written in pure Go. We don't just use a single database; we use a **multi-instance architecture**, where different types of data are stored in separate, specially-tuned BadgerDB instances. This is critical for achieving maximum performance.

*   **`history.db` (The Archive):**
    *   **Stores:** Immutable, content-addressable objects (`Commit`, `Tree`, `Blob`).
    *   **Access Pattern:** Write-once, read-rarely (WORR). This is an append-only ledger.
    *   **Tuning:** Optimized for large values and storage efficiency. It uses a **low `ValueThreshold`** to push large blob data into Badger's Value Log and enables **ZSTD compression** to minimize disk footprint.

*   **`index.db` (The Hot Index):**
    *   **Stores:** High-churn data: query indices (`spog:`, `posg:`), references (`ref:`), and temporary data for scalable operations.
    *   **Access Pattern:** Frequent writes, deletes, and fast, sorted prefix scans.
    *   **Tuning:** Optimized for read performance. It uses a **high `ValueThreshold`**, forcing small key-value pairs to live directly in the LSM-tree. This makes index scans incredibly fast as no secondary disk lookup is needed. A large block cache is allocated to keep as much of the index in RAM as possible.

*   **`app.db` (The Application State):**
    *   **Stores:** Mutable, non-versioned data like Merge Requests, user sessions, and job queue state.
    *   **Access Pattern:** Classic Online Transaction Processing (OLTP) with frequent updates to the same keys.
    *   **Tuning:** Similar to `index.db`, it's tuned for low-latency reads and writes on small values to ensure a responsive UI for application features.

This separation allows us to tailor the storage engine's behavior to the data's lifecycle, avoiding the performance compromises of a "one-size-fits-all" approach.

### 2. The Core API (`quadstore`): The Brains of the Operation

This is the heart of the project. It is a "headless" Go library that provides a clean, public API for all versioning and query operations. This layer acts as a strict firewall between the application logic and the low-level database details.

*   **The Contract (`quadstore.Store` interface):** This Go interface defines every high-level action possible in the system (`Commit`, `Merge`, `Diff`, `Blame`, `Log`, etc.). It is the formal contract that the rest of the application will code against.

*   **Concurrency Safety:** This is the layer responsible for ensuring thread safety. **Every write operation is wrapped in a single, atomic `db.Update()` transaction.** This leverages BadgerDB's Serializable Snapshot Isolation (SSI) to prevent race conditions and guarantee data consistency, no matter how many concurrent requests are made.

*   **Scalability Logic:** All the scalable algorithms from Chapter 5 are implemented here. The complex, iterator-based logic for `diff`, `merge`, and `blame` is completely encapsulated within this layer.

*   **Orchestration:** The concrete implementation of the `Store` interface (`Repository`) holds connections to all three database instances (`history.db`, `index.db`, `app.db`) and is responsible for routing reads and writes to the correct one. A `Commit` operation, for example, will write to both `history.db` and `index.db` within a coordinated transaction.

### 3. The Presentation Layer: The "Heads"

This layer is responsible for user interaction. It consumes the `quadstore` Core API to provide a usable interface. Because the core API is decoupled, we can have multiple "heads" for different purposes.

*   **`quadgit` CLI:**
    *   **Role:** A powerful tool for developers, data scientists, and CI/CD automation.
    *   **Behavior:** A **short-lived process.** When you run `quadgit commit`, the application starts, opens the database, calls the `store.Commit()` method, prints the result, and exits. It is designed for human-readable output and scriptability.

*   **`quadgit-server` (REST API):**
    *   **Role:** A multi-user, collaborative platform for web applications.
    *   **Behavior:** A **long-running, stateful service.** It opens the database stores once at startup and serves requests concurrently. It is responsible for:
        *   **Authentication & Authorization:** Verifying user identity and checking permissions.
        *   **HTTP Semantics:** Translating user actions into the correct HTTP verbs, status codes, and headers (like `ETag` for optimistic locking).
        *   **Higher-Level Workflows:** Managing application-specific concepts like Merge Requests, which live entirely within this layer.

## Roadmap

This project is ambitious and under active development. Our roadmap is structured to deliver value incrementally.

-   [ ] **Phase 1: Core CLI:** `init`, `add`, `commit`, `log`, `branch`, `merge`.
-   [ ] **Phase 2: Scalability & Production Readiness:** Stream-based algorithms, multi-DB architecture, GPG signing, `blame`, `backup`.
-   [ ] **Phase 3: The Networked Service:** REST API, authentication/authorization, optimistic locking.
-   [ ] **Phase 4: Collaborative Workflows:** Namespaces, Merge Request system.
-   [ ] **Phase 5: Ecosystem Integration:** SPARQL endpoint, `push`/`pull` replication.

## Contributing

We welcome contributions of all kinds! Whether it's reporting a bug, proposing a new feature, improving documentation, or writing code, your help is appreciated. This project is in its early stages, and it's a great time to get involved and help shape its future.

Please see our [CONTRIBUTING.md](CONTRIBUTING.md) guide for more details on our development process, code of conduct, and how to submit a pull request.

## License

`quadgit` is licensed under the [Apache 2.0 License](LICENSE).