
### **Jira Project: Graph Git (GG)**
  **GG-2 (Story): Define Core API & Structs**
    *   AC 1: A Go interface named `Store` exists in `pkg/quadstore`.
    *   AC 2: The `Store` interface includes method signatures for at least `Init`, `Commit`, `ReadCommit`, and `ResolveRef`.
    *   AC 3: Public data structs `Commit`, `Tree`, `Blob`, `Quad`, and `Author` are defined in `pkg/quadstore/types.go` with appropriate JSON tags.
    *   AC 4: All public structs and interface methods are documented with GoDoc comments.

*   **GG-3 (Task): Implement `Repository` Struct**
    *   AC 1: A `Repository` struct exists in `internal/datastore` that correctly implements the `quadstore.Store` interface.
    *   AC 2: The `Repository` holds a `*badger.DB` instance and a `namespace` string.
    *   AC 3: All methods that interact with BadgerDB correctly prepend `ns:<namespace>:` to the keys.
    *   AC 4: A public constructor function `quadstore.Open(...)` is implemented, which returns the `Repository` as a `Store` interface.

*   **GG-4 (Story): As a User, I can run `quad-db init`...**
    *   AC 1: **Given** a directory without a `.quad-db` folder, **When** I run `quad-db init`, **Then** a `.quad-db` directory is created.
    *   AC 2: **And** the database contains a single root `Commit` object with no parents.
    *   AC 3: **And** a `ref:head:main` reference exists, pointing to the root commit's hash.
    *   AC 4: **And** a `HEAD` reference exists with the value `"ref:head:main"`.
    *   AC 5: **Given** an existing repository, **When** I run `quad-db init` again, **Then** the command prints an error message and exits with a non-zero status code.

*   **GG-5 (Story): As a User, I can run `quad-db add <file>`...**
    *   AC 1: **Given** an initialized repository, **When** I run `quad-db add data.nq`, **Then** the quads from the file are added to a temporary staging area (e.g., an index file or dedicated DB keys).
    *   AC 2: The command must support a simple syntax for additions (default) and deletions (e.g., lines prefixed with `D` or `-`).

*   **GG-6 (Story): As a User, I can run `quad-db commit -m "..."`...**
    *   AC 1: **Given** quads have been staged using `add`, **When** I run `quad-db commit -m "My first commit"`, **Then** a new `Commit` object is created in the database.
    *   AC 2: **And** the new commit's parent is correctly set to the previous `HEAD` commit.
    *   AC 3: **And** the `ref:head:main` reference is updated to point to the new commit hash.
    *   AC 4: **And** the staging area is cleared after the commit is successful.
    *   AC 5: **Given** an empty staging area, **When** I run `commit`, **Then** the command prints "Nothing to commit" and exits.

*   **GG-7 (Story): As a User, I can run `quad-db log`...**
    *   AC 1: **When** I run `quad-db log`, **Then** the output displays a list of commits from the current branch's `HEAD` backward.
    *   AC 2: Each entry must display at least the commit hash, author, date, and message.
    *   AC 3: The command correctly follows the parent chain back to the root commit.

*   **GG-8 (Task): Create Initial Test Suite**
    *   AC 1: A `repository_test.go` file exists.
    *   AC 2: There is a test case that successfully runs the full `init` -> `add` -> `commit` -> `log` workflow.
    *   AC 3: Test helpers exist to create and clean up temporary BadgerDB instances for isolated test runs.

*   **GG-9 (Story): As a User, I can run `quad-db status`...**
    *   AC 1: **Given** changes have been staged with `add`, **When** I run `quad-db status`, **Then** the output lists the staged additions and deletions.
    *   AC 2: **Given** the staging area is empty, **When** I run `quad-db status`, **Then** the output indicates there are no changes staged for commit.

#### **Epic 2: GG-10 - Branching & History Inspection**

*   **GG-11 (Story): As a User, I can manage branches...**
    *   AC 1: **When** I run `quad-db branch feature/login`, **Then** a new ref `ref:head:feature/login` is created pointing to the current `HEAD` commit.
    *   AC 2: **When** I run `quad-db branch`, **Then** the output lists all existing branches, with the current branch marked with an asterisk.
    *   AC 3: **When** I run `quad-db branch -d feature/login`, **Then** the `ref:head:feature/login` key is deleted.

*   **GG-12 (Story): As a User, I can switch my active context...**
    *   AC 1: **Given** I am on the `main` branch, **When** I run `quad-db checkout feature/login`, **Then** the `HEAD` key's value is updated to `"ref:head:feature/login"`.
    *   AC 2: **And** a subsequent `commit` will add to the `feature/login` branch, not `main`.

*   **GG-13 (Story): As a User, I can create and list tags...**
    *   AC 1: **When** I run `quad-db tag v1.0`, **Then** a new ref `ref:tag:v1.0` is created pointing to the current `HEAD` commit.
    *   AC 2: **When** I run `quad-db tag`, **Then** the output lists all existing tags.

*   **GG-14 (Story): As a User, I can run `quad-db diff`...**
    *   AC 1: **When** I run `quad-db diff <hash1> <hash2>`, **Then** the output shows quads present in `<hash2>` but not `<hash1>` prefixed with `+`.
    *   AC 2: **And** the output shows quads present in `<hash1>` but not `<hash2>` prefixed with `-`.
    *   AC 3: The command works correctly when using branch or tag names as arguments.

*   **GG-15 (Story): As a User, I can run `quad-db show`...**
    *   AC 1: **When** I run `quad-db show <hash>`, **Then** the output first displays the commit metadata (author, date, message).
    *   AC 2: **And** the output then displays the diff between that commit and its primary parent.

#### **Epic 3: GG-16 - Merging & Concurrency Hardening**

*   **GG-17 (Story): Implement three-way merge base detection**
    *   AC 1: Given two commit hashes, the algorithm must correctly identify the nearest common ancestor commit.
    *   AC 2: The algorithm works correctly for both simple linear histories and post-merge histories.

*   **GG-18 (Story): As a User, I can run `quad-db merge`...**
    *   AC 1: **Given** no conflicts, **When** I run `quad-db merge feature`, **Then** a new merge commit is created on the current branch.
    *   AC 2: **And** the new merge commit has two parents: the previous branch `HEAD` and the `feature` branch `HEAD`.
    *   AC 3: **And** the state of the graph after the merge correctly reflects the combination of changes from both branches.

*   **GG-19 (Task): Implement basic conflict detection...**
    *   AC 1: The merge operation must detect when one branch adds a quad and another branch deletes the same quad, and halt with a conflict error.
    *   AC 2: The merge operation must detect when two branches add a quad with the same Subject-Predicate but different Objects, and halt with a conflict error.

*   **GG-20 (Story): Create a "Concurrent Test Harness"...**
    *   AC 1: A new test function `TestConcurrentCommits` exists.
    *   AC 2: The test launches at least 10 goroutines that concurrently call the `Store.Commit()` method on the same branch.
    *   AC 3: The test uses `sync.WaitGroup` to ensure it waits for all goroutines to complete.

*   **GG-21 (Story): Create an "Auditor" function...**
    *   AC 1: An `assertRepositoryIsInvariants` function exists.
    *   AC 2: The function checks for reference integrity (all refs point to existing commits).
    *   AC 3: The function checks for object integrity (all commit/tree/blob links are valid).
    *   AC 4: The function checks for history integrity (all parent links are valid).
    *   AC 5: The `TestConcurrentCommits` test calls this auditor function upon completion.

*   **GG-22 (Bug): The "Commit Race" test fails...**
    *   AC 1: The `TestConcurrentCommits` test, when run with `go test -race`, reports a data race.
    *   AC 2: (After Fix) The `TestConcurrentCommits` test passes with `go test -race`.
    *   AC 3: (After Fix) The `assertLinearHistory` check within the test confirms no commits were lost and the history is a perfect chain.

*   **GG-23 (Task): Set up CI pipeline...**
    *   AC 1: A GitHub Actions (or similar) workflow is created.
    *   AC 2: The workflow runs `go test -race ./...` on every push to `main` and on every pull request.
    *   AC 3: The build fails if any tests fail or if the race detector finds an issue.

#### **Epic 4: GG-24 - Scalability & Performance Engineering**
*Description: Overhaul the core algorithms and data layout to handle massive datasets and improve performance under load.*

*   **GG-25 (Story):** As a Developer, I need to refactor the `diff` algorithm to be fully stream-based using BadgerDB iterators.
*   **GG-26 (Story):** As a Developer, I need to refactor the `merge` algorithm to use the new scalable diff implementation.
*   **GG-27 (Story):** As a User, I can run the `quad-db stats` command to get both high-level and detailed graph statistics.
*   **GG-28 (Task):** Implement HyperLogLog for unique entity counts in the `stats data` command.
*   **GG-29 (Story):** As a DevOps Engineer, I need to refactor the datastore to use a multi-instance BadgerDB architecture (`history.db`, `index.db`).
*   **GG-30 (Task):** Update the `Commit` struct to include pre-computed churn statistics (`added`/`deleted`) to accelerate `stats history`.

#### **Epic 5: GG-31 - Trust, Provenance, and Production Operations**
*Description: Implement features for data integrity, auditability, and day-to-day database administration.*

*   **GG-32 (Story):** As a User, I can sign my commits with my GPG key using `quad-db commit -S`.
*   **GG-33 (Story):** As a User, I can verify the GPG signature of commits using `quad-db log --show-signature`.
*   **GG-34 (Story):** As a User, I can run `quad-db blame <graph>` to see the provenance of each quad in a named graph.
*   **GG-35 (Story):** As an Administrator, I can perform an online backup of the entire repository using `quad-db backup`.
*   **GG-36 (Story):** As an Administrator, I can restore a repository from a backup using `quad-db restore`.
*   **GG-37 (Story):** As a Data Engineer, I can perform a high-speed initial data import using `quad-db bulk-load`.
*   **GG-38 (Story):** As a Data Scientist, I can improve query performance on historical data using `quad-db materialize <tag>`.

#### **Epic 6: GG-39 - The REST API Server**
*Description: Build the server application that exposes the core `quad-db` functionality over a secure, multi-user REST API.*

*   **GG-40 (Task):** Create the new `quad-db-server` binary with an HTTP router and shared `Store` instance.
*   **GG-41 (Story):** As a Developer, I can authenticate with the API using a bearer token (API Key or JWT).
*   **GG-42 (Story):** As an Administrator, I can define user roles (RBAC) on a per-namespace basis.
*   **GG-43 (Task):** Implement authentication and authorization middleware for all endpoints.
*   **GG-44 (Story):** As a Developer, I can create and view commits via `POST /commits` and `GET /commits/:hash`.
*   **GG-45 (Story):** As a Developer, I can manage branches via the `GET /refs/heads` and `POST /refs/heads` endpoints.
*   **GG-46 (Story):** As a Developer, my write requests are rejected with `HTTP 412` if my data is stale, enforced by `ETag`/`If-Match` headers.
*   **GG-47 (Story):** As a Developer, long-running operations like `merge` are handled asynchronously via a job queue.


### **Epic 7: GG-48 - Multi-Tenancy & Collaborative Workflows**

*   **GG-49 (Story): As an Administrator, I can create, list, and delete namespaces using the `quad-db ns` commands.**
    *   **AC 1 (`list`):** **Given** namespaces "default" and "proj-a" exist, **When** I run `quad-db ns list`, **Then** the output contains both "default" and "proj-a".
    *   **AC 2 (`create`):** **When** I run `quad-db ns create proj-b`, **Then** the `sys:namespaces` key in the database is updated to include "proj-b".
    *   **AC 3 (`create`):** **And** the new "proj-b" namespace contains a root commit and a `main` branch, fully prefixed with `ns:proj-b:`.
    *   **AC 4 (`create`):** **Given** a namespace "proj-b" already exists, **When** I run `quad-db ns create proj-b`, **Then** an error message is printed and the command exits.
    *   **AC 5 (`rm`):** **Given** a namespace "proj-b" exists, **When** I run `quad-db ns rm proj-b --force`, **Then** all keys with the prefix `ns:proj-b:` are deleted from the database.
    *   **AC 6 (`rm`):** **When** I run `quad-db ns rm proj-b` without the `--force` flag, **Then** the command prints a confirmation warning and does not delete any data.

*   **GG-50 (Story): As a User, I can specify a namespace for any command using the `-n` flag.**
    *   **AC 1:** **Given** two namespaces "proj-a" and "proj-b" with different histories, **When** I run `quad-db log -n proj-a`, **Then** I see the commit history for "proj-a".
    *   **AC 2:** **And When** I run `quad-db log -n proj-b` immediately after, **Then** I see the different commit history for "proj-b".
    *   **AC 3:** **When** I run `quad-db ns use proj-a`, **Then** a subsequent `quad-db log` command (without the `-n` flag) shows the history for "proj-a".
    *   **AC 4:** All core commands (`commit`, `log`, `diff`, `branch`, `stats`, etc.) must respect the `-n` flag and the `ns use` context.

*   **GG-51 (Story): As a User, I can copy a branch from one namespace to another using `quad-db cp`.**
    *   **AC 1:** **Given** a branch `feature` exists in namespace "dev" but not "staging", **When** I run `quad-db cp dev:feature staging:feature-copy`, **Then** a new branch `feature-copy` is created in the "staging" namespace.
    *   **AC 2:** **And** the new branch points to a commit with the same hash as the original `feature` branch `HEAD`.
    *   **AC 3:** **And** all necessary `Commit`, `Tree`, and `Blob` objects from the "dev" history have been copied to the "staging" namespace key-space.
    *   **AC 4:** The command must efficiently skip copying any objects that already exist in the target namespace (due to content-addressing).

*   **GG-52 (Task): Design and implement the non-versioned `MergeRequest` data model and `mrstore`.**
    *   AC 1: A `MergeRequest` Go struct is defined, including fields for `ID`, `Version`, `Status`, `Title`, source/target branches, and author.
    *   AC 2: A `mrstore.MergeRequestStore` interface is defined with `Create`, `Get`, `Update`, and `List` methods.
    *   AC 3: A concrete implementation of the store is created that uses the `app.db` BadgerDB instance.
    *   AC 4: The implementation uses a dedicated key prefix `app:mr:data:<id>` for MR objects and `app:mr:sequence` for atomically generating unique IDs.
    *   AC 5: The `Update` method must implement the optimistic locking check based on the `Version` field.

*   **GG-53 (Story): As a User, I can create and list Merge Requests via the REST API.**
    *   **AC 1 (`POST`):** **When** I send a `POST` request to `/mergerequests` with a valid title and source/target branches, **Then** the server responds with `HTTP 201 Created`.
    *   **AC 2 (`POST`):** **And** a new MR object is created in the `app.db` with status "open" and a unique ID.
    *   **AC 3 (`GET`):** **When** I send a `GET` request to `/mergerequests`, **Then** I receive a JSON array of `MergeRequest` objects.
    *   AC 4 (`GET`): The `GET /mergerequests` endpoint must support filtering by status (e.g., `?status=open`).

*   **GG-54 (Story): As a Reviewer, I can view the diff and commit history for a Merge Request.**
    *   **AC 1:** **When** I send a `GET` request to `/mergerequests/:id`, **Then** the JSON response contains the full `MergeRequest` object.
    *   AC 2: **And** the response includes a `diff` field containing the result of calling `quadstore.Store.Diff()` between the MR's target and source branches.
    *   AC 3: **And** the response includes a `commits` field containing an array of `Commit` objects from the source branch's history.
    *   AC 4: The response includes an `ETag` header containing the MR's current `Version` number.

*   **GG-55 (Story): As a Maintainer, I can execute a merge via the `POST /mergerequests/:id/merge` endpoint.**
    *   **AC 1:** **Given** an open MR, **When** I send a `POST` request to `/mergerequests/:id/merge` with a valid `If-Match` header (containing the target branch's `HEAD` hash), **Then** the server responds with `HTTP 200 OK`.
    *   **AC 2:** **And** a new merge commit is successfully created on the target branch in the `quadstore`.
    *   **AC 3:** **And** the MR object's status in `app.db` is updated to "merged", with the `merged_by` and `merge_commit_hash` fields populated.
    *   **AC 4:** **Given** the `If-Match` header does not match the current `HEAD` of the target branch, **When** I send the request, **Then** the server responds with `HTTP 412 Precondition Failed` and no merge is performed.
    *   **AC 5:** The endpoint must check for user authorization (e.g., the user must have the `maintainer` role) and return `HTTP 403 Forbidden` if they are not permitted to merge.

#### **Epic 8: GG-56 - Distributed Synchronization**
*Description: Implement the `push` and `pull` mechanisms to allow repositories to be synchronized between a client and a remote server.*

*   **GG-57 (Story):** As a Developer, I need to design the RPC protocol for object synchronization, including `list-refs`, `get-objects`, and `update-ref`.
*   **GG-58 (Task):** Implement the server-side RPC endpoints (`/rpc/...`) in the `quad-db-server`.
*   **GG-59 (Story):** As a User, I can add a named remote repository using `quad-db remote add <name> <url>`.
*   **GG-60 (Story):** As a User, I can run `quad-db fetch <remote>` to download objects and update my remote-tracking branches without merging.
*   **GG-61 (Task):** Implement the client-side "packfile" generation logic for `fetch`.
*   **GG-62 (Story):** As a User, I can run `quad-db push <remote> <branch>` to send my local commits to the remote server.
*   **GG-63 (Task):** Implement the client-side "packfile" creation and transfer logic for `push`.
*   **GG-64 (Story):** As a User, `push` operations on a protected branch are rejected by the server unless I have the correct permissions.
*   **GG-65 (Story):** As a User, I can run `quad-db pull <remote> <branch>` as a convenience command that performs a `fetch` followed by a `merge`.

#### **Epic 9: GG-66 - SPARQL 1.1 Protocol Compliance**
*Description: Implement a standard SPARQL endpoint to allow third-party RDF tools to query the database.*

*   **GG-67 (Task):** Choose and integrate a Go library for parsing SPARQL query strings into an Abstract Syntax Tree (AST).
*   **GG-68 (Story):** As a Developer, I need to implement a query planner that resolves single Basic Graph Patterns (BGPs) using the optimal `spog/posg` indices.
*   **GG-69 (Story):** As a Developer, the query engine must support joining multiple triple patterns using a scalable merge-join or hash-join algorithm.
*   **GG-70 (Story):** As a Developer, I need to implement a query optimizer that reorders triple patterns based on cardinality estimates to improve performance.
*   **GG-71 (Story):** As a SPARQL User, my queries can include `FILTER`, `OPTIONAL`, and `UNION` clauses.
*   **GG-72 (Story):** As a Client Application, I can send a query to the `/sparql` endpoint via `GET` or `POST`.
*   **GG-73 (Story):** As a Client Application, I can request query results in standard formats (`application/sparql-results+json`, `application/sparql-results+xml`) using content negotiation.
*   **GG-74 (Story):** As a SPARQL User, my queries can include the `SERVICE` keyword to execute a federated sub-query against an external SPARQL endpoint.

#### **Epic 10: GG-75 - Advanced Testing & Resilience**
*Description: Implement advanced testing methodologies to prove the system's logical correctness under complex scenarios and its resilience to failure.*

*   **GG-76 (Story):** As a Developer, I need to implement a property-based test for "revert invariance" to find edge cases in the diff/commit logic.
*   **GG-77 (Story):** As a Developer, I need to implement a property-based test for "merge idempotence" to validate the merge algorithm.
*   **GG-78 (Story):** As a Developer, I need to implement a property-based test for "backup/restore integrity" to ensure perfect data serialization.
*   **GG-79 (Story):** As a DevOps Engineer, I need to create a "Crash Test" harness that uses `kill -9` during concurrent writes and then runs the Auditor to verify durability.
*   **GG-80 (Task):** Integrate `toxiproxy` into the CI pipeline to test the resilience of `push`/`pull` operations against network latency and connection drops.
*   **GG-81 (Story):** As a Developer, I need to build a simple, in-memory "golden model" implementation of the `Store` interface.
*   **GG-82 (Task):** Create a "Model-Based Test" suite that runs complex operation sequences against both the real implementation and the golden model, asserting their final states are identical.
*   **GG-83 (Story):** As a QA Engineer, I need to set up a long-running "Soak Test" environment with `pprof` monitoring to detect memory leaks or performance degradation over time.
