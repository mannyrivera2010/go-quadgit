
### **The Go Project Standard Layout (Adapted for `go-quadgit`)**

This structure is a widely adopted convention in the Go community. It is not enforced by the compiler, but it provides a predictable and logical place for everything.

```
/go-quadgit/
├── /cmd/
│   ├── /go-quadgit/         # CLI Application
│   │   └── main.go
│   └── /go-quadgit-server/  # REST Server Application
│       └── main.go
│
├── /pkg/
│   └── /quadstore/          # The Public Core API
│       ├── api.go           # The `Store` interface
│       └── types.go         # Public structs: Commit, Quad, Author, etc.
│
├── /internal/
│   ├── /datastore/          # Concrete implementation of the `Store` interface
│   │   ├── repository.go    # The `Repository` struct with its methods
│   │   └── badger_helpers.go # Low-level BadgerDB utility functions
│   │
│   ├── /server/             # REST server-specific logic
│   │   ├── middleware/      # AuthN, AuthZ, logging middleware
│   │   ├── handlers/        # HTTP handlers for each endpoint
│   │   └── router.go        # Gin router setup
│   │
│   ├── /rpc/                # Logic for the Git-style push/pull protocol
│   │   ├── handlers.go
│   │   └── packfile.go
│   │
│   └── /config/             # Configuration loading and parsing (Viper)
│       └── config.go
│
├── /docs/                   # All project documentation (the "Books")
│   └── book1_primer/
│       ├── chapter1.md
│       └── ...
│
├── /scripts/                # Helper scripts for build, test, etc.
│   ├── test-crash.sh
│   └── run-dev-server.sh
│
├── /web/                    # The React frontend application
│   ├── /src/
│   └── package.json
│
├── go.mod
├── go.sum
├── README.md
├── LICENSE
└── CONTRIBUTING.md
```

---

### **Detailed Explanation of Each Directory**

#### **`/cmd` - The Entrypoints**

*   **Purpose:** This directory contains the `main` packages for the executables we want to build. It is the only place where you will find `package main`.
*   **Scalability:** This structure makes it trivial to add new binaries in the future (e.g., `go-quadgit-admin-tool`, `go-quadgit-worker`) without cluttering the root directory. Each folder inside `/cmd` corresponds to one compiled application.
*   **Logic:** The code in here should be "thin." Its only job is to parse command-line flags, load configuration, initialize the core application objects (like the `Store`), and start the application (either by executing a CLI command or by starting a server).

#### **`/pkg` - The Public Library (`quadstore`)**

*   **Purpose:** This is the **public, shareable library**. Any external Go project that wants to embed `go-quadgit`'s core logic would import this package.
*   **Key Content:**
    *   `api.go`: The `quadstore.Store` interface. This is the **formal contract** of the core engine.
    *   `types.go`: The public data structures (`Commit`, `Quad`, `Author`, `Change`, etc.) that are used as arguments and return values for the `Store` interface.
*   **Scalability:** By defining a clean, stable public API, we can evolve the internal implementation (`/internal`) without breaking external consumers. This is crucial for building an ecosystem.

#### **`/internal` - The Private Implementation**

*   **Purpose:** This is where the vast majority of your Go code will live. The `internal` directory is a special feature of the Go toolchain: code inside it can only be imported by code within the same repository (i.e., within the `/cmd` directory). This prevents external projects from depending on your private, unstable implementation details.
*   **Scalability:** This gives you the freedom to refactor and change your internal logic aggressively without worrying about breaking downstream users. It enforces the clean API boundary defined in `/pkg`.
*   **Subdirectories:**
    *   `/datastore`: The heart of the engine. The concrete `Repository` struct that implements the `quadstore.Store` interface lives here. All the BadgerDB transactions, key prefixing, and scalable algorithms are implemented here.
    *   `/server`: Contains all the logic specific to the `go-quadgit-server`. This includes HTTP handlers, routing setup, and middleware for authentication and authorization. This code is not needed by the CLI.
    *   `/rpc`: Contains the handlers for the Git-style synchronization protocol (`push`/`pull`).
    *   `/config`: A utility package for loading and parsing the `config.yml` file using a library like Viper.

#### **`/web` - The Frontend Application**

*   **Purpose:** A dedicated, top-level directory for the entire React single-page application.
*   **Scalability:** This completely decouples the frontend from the backend. The frontend team can work in this directory using their own tools (`npm`, `vite`) and build process. The only connection to the backend is through the REST API contract. The Go server can even be configured to serve the static build artifacts from this directory.

#### **`/docs` - The Knowledge Base**

*   **Purpose:** Your "Books" live here as Markdown files. This is the single source of truth for all project documentation.
*   **Scalability:** As the project grows, you can structure this directory just like the code: `/docs/layer1_core/book1_primer/chapter1.md`. This makes the documentation discoverable and easy to maintain.

#### **`/scripts` - Automation & Tooling**

*   **Purpose:** Contains helper scripts for development and CI/CD that aren't part of the main Go application.
*   **Examples:** A shell script to run the server with the correct environment variables for development, a script to run the crash tests, or a script to build all binaries and Docker images for a release.

### **How This Structure Promotes Scalability**

1.  **Clear Ownership:** It's immediately obvious where to put new code. Working on a CLI command? Go to `/cmd/go-quadgit`. Fixing a database bug? Go to `/internal/datastore`. Building a new REST endpoint? Go to `/internal/server/handlers`. This is essential for a growing team.
2.  **Enforced API Boundary:** The `internal` directory prevents "cheating." Developers are forced to use the public `quadstore.Store` interface, ensuring the application layers remain decoupled from the core logic.
3.  **Independent Build Processes:** The `go-quadgit` CLI can be built without pulling in any of the REST server's dependencies (like Gin). The `web` frontend has its own completely separate `npm`-based build process.
4.  **Testability:** The separation makes testing much cleaner. You can write unit/integration tests for the `datastore` package in complete isolation, without needing to spin up a CLI or an HTTP server.
