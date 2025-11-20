# **Book 3: Building the Multi-User Server**

This book provides all the necessary theory and practical code to take the single-user `go-quadgit` engine and forge it into a hardened, secure, and highly available networked service, ready for collaborative use.

**Layer II: The Networked Platform**

**Sub-System 3: Service Architecture & Concurrency**

**Subtitle:** *From a Local CLI to a Concurrent, Multi-Protocol, Networked Service*

**Book Description:** In the first two books, we forged a powerful, scalable, single-user engine for versioning knowledge graphs. Now, it's time to unleash that power and make it available to the world. This book is the definitive guide to building the `go-quadgit-server`, a long-running, multi-protocol daemon that transforms our core library into a true multi-user platform.

We will start by architecting a server that can listen for different kinds of connections simultaneously. You will learn to build a clean, modern REST API for web applications, a highly efficient RPC protocol for Git-style synchronization, and a remote execution endpoint for "thin client" operations. Along the way, we will tackle the most critical challenge of any server application: handling concurrency safely and reliably. This book provides the blueprint for building a robust, secure, and highly available service on top of the `go-quadgit` core engine.

**Prerequisites:** Readers must have mastered the concepts from **Book 1** and **Book 2**, especially the Core API (`quadstore.Store`) and the transactional guarantees of the datastore. A strong understanding of Go, web services (REST, HTTP), and concurrency (goroutines, mutexes) is essential.

### **Part I: Server Architecture and Design**

*This part lays the architectural groundwork for our networked service.*

#### **Chapter 1: The Multi-Protocol Daemon**
*   **1.1: The Mindset Shift: From Stateless CLI to Stateful Service**
    *   A deep dive into the fundamental differences between a short-lived command-line process and a long-running server. We'll cover state management (opening the database once), the need for graceful shutdowns, and the responsibility of handling untrusted input from the network.
*   **1.2: The Three Server Modes: A Unified Architecture**
    *   Introduces the concept of the unified `go-quadgit-server` binary that listens on multiple ports. We'll detail the distinct purpose of each mode:
        *   **REST API:** For web UIs and application integration.
        *   **Git Sync (RPC):** For `push`/`pull` data synchronization.
        *   **Remote Execution (gRPC):** For offloading heavy client-side commands to the server.
*   **1.3: Configuration Management with Viper and YAML**
    *   A practical guide to building a robust configuration system using a `config.yml` file. We'll define the structure for server addresses, database paths, and security settings.
*   **1.4: The Server Entrypoint: `main.go`**
    *   Presents the complete `main.go` for the server, showing how to initialize the multi-instance `Store`, parse the configuration, and launch separate goroutines to listen on each configured port. Includes signal handling for graceful shutdown.
