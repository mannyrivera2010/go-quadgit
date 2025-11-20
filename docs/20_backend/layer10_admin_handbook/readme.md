
## **Part II: Production Operations and Extensibility**

*This part is the administrator's handbook, covering security, data management, and custom integrations.*

## **Chapter 4: Securing Data at Rest with Encryption**
*   **4.1: The Key Hierarchy: KEKs, DEKs, and Password Derivation**
    *   A detailed cryptographic design chapter explaining the roles of the master Key Encryption Key (from `QUADGIT_MASTER_KEY`), the per-namespace Data Encryption Key, and why a KDF like Argon2id is necessary to turn user passwords into secure keys.
*   **4.2: Implementing Per-Namespace, Password-Based Encryption**
    *   The full Go implementation. This includes modifying the `Repository` to manage a cache of unlocked ciphers, creating the `ns:<ns>:sys:crypto` object to store salt and the encrypted DEK, and refactoring `writeObject`/`readObject` to use the correct per-namespace cipher.
*   **4.3: The "Unlock" Workflow and Session Management**
    *   Implements the `POST /unlock` REST endpoint and details how a successful unlock can generate a short-lived JWT that allows subsequent requests to be processed without re-submitting the password.
*   **4.4: The Operational Lifecycle: Key Rotation and Data Migration**
    *   Provides administrator runbooks for critical security operations, such as changing a namespace password and the more complex process of migrating an unencrypted repository to an encrypted one.

## **Chapter 5: The Administrator's Handbook**
*   **5.1: The Production Operations Manual: Backup and Restore**
    *   A practical guide for administrators. It includes strategies for scheduling nightly full backups and hourly incremental backups, and how to stream them directly to cloud storage. Includes a step-by-step disaster recovery plan.
*   **5.2: Large-Scale Data Management: `bulk-load` and `materialize`**
    *   A deep dive into the use cases for these specialized commands. Provides performance benchmarks comparing `bulk-load` to standard commits and `materialize`'s effect on historical query times.
*   **5.3: Monitoring and Insights: A Guide for SREs**
    *   Details the key metrics the `go-quadgit-server` should expose via a `/metrics` endpoint for Prometheus. This includes database-level metrics (cache hit ratios, LSM-tree size), application-level metrics (request latency, error rates), and business-level metrics (commits per hour, number of active users).

## **Chapter 6: The Platform Extensibility Guide**
*   **6.1: The Plugin and Hook System**
    *   The formal design of the hook system, detailing all available hook points (`pre-commit`, `post-commit`, `pre-merge`, `post-receive`, etc.) and the data passed to each script via `stdin`.
*   **6.2: Use Case Deep Dive: Pre-Commit Validation with SHACL**
    *   A complete, end-to-end tutorial. This includes writing the `shapes.ttl` file, creating the `pre-commit` shell script that invokes a SHACL validator, and showing how the error messages are propagated back to the user's CLI or the REST API response.
*   **6.3: Use Case Deep Dive: Post-Commit CI/CD Integration**
    *   Shows how to write a `post-receive` hook on the server that, upon a push to the `main` branch, triggers a Jenkins or GitLab CI pipeline via a webhook `curl` request.
*   **6.4: Building Applications on the `go-quadgit` Platform**
    *   A final, forward-looking section that revisits the application ideas (Governance Studio, ML Lineage) and provides architectural sketches showing how they would use the platform's features, from the REST API to hooks, to build their product.

