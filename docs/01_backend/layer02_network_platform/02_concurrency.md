
### **Part II: Concurrency, Safety, and Long-Running Tasks**

*This part tackles the hardest problems in server-side development: ensuring data is never corrupted and the service remains responsive under load.*

#### **Chapter 3: Transactional Integrity and ACID Compliance**
*   **3.1: A Review of BadgerDB's Guarantees (MVCC & SSI)**
    *   A quick refresher on why the underlying database provides the tools for safety, focusing on snapshot isolation.
*   **3.2: The Datastore's Contract: Atomicity with `db.Update()`**
    *   Reinforces that the `Store` implementation's primary responsibility is to wrap every logical write operation in an atomic transaction. This is the bedrock of concurrency safety.
*   **3.3: Case Study: How Transactions Prevent a "Commit Race"**
    *   A detailed, step-by-step walkthrough of the classic concurrent commit scenario. It explains exactly how BadgerDB's optimistic locking and automatic retry mechanism work at a low level to prevent a lost update and ensure a linear history.
*   **3.4: A Formal Review of ACID Properties for `go-quadgit`**
    *   Systematically evaluates the full application architecture against the definitions of Atomicity, Consistency, Isolation, and Durability, clarifying the responsibilities of each layer.

#### **Chapter 4: Managed Optimistic Locking in the API**
*   **4.1: Why Automatic Retries Are Not Enough for a UI**
    *   Explains the user experience problem: a user needs to be *told* their data is stale; a silent backend retry can be confusing.
*   **4.2: The ETag and `If-Match` Workflow**
    *   A complete guide to implementing this standard REST pattern for optimistic locking.
*   **4.3: Implementing the `If-Match` Check in a Middleware**
    *   Shows the Go code for a middleware that extracts the ETag, reads the current resource version, and returns `HTTP 412 Precondition Failed` if they don't match.
*   **4.4: Identifying All Endpoints That Need Protection**
    *   Provides a checklist of every write-based endpoint (`POST /commits`, `POST /merges`, `DELETE /refs/heads`, etc.) and explains why each requires this protection.

#### **Chapter 5: Handling Long-Running Operations with Asynchronous Job Queues**
*   **5.1: The HTTP Timeout Problem**
    *   Explains why a synchronous `merge` operation in an HTTP handler is a recipe for disaster.
*   **5.2: Designing the Asynchronous Workflow**
    *   Details the `202 Accepted` response pattern, the job status object, and the polling mechanism.
*   **5.3: Implementation: A Background Worker Pool**
    *   Provides a practical Go implementation of a job queue using channels and a pool of long-lived worker goroutines that consume tasks.
*   **5.4: Storing and Updating Job State in `app.db`**
    *   Shows how to use the `app.db` instance to store the status (`pending`, `complete`, `failed`) and result of each job, making the system resilient to server restarts.
