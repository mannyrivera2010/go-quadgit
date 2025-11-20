
## 
**Layer IV: Operations & Reliability**
**Sub-System 7: System Reliability & Testing**

# **Book 15: Engineering for Resilience**

**Subtitle:** *A Practical Guide to Fault Injection, Chaos Engineering, and Durability Testing*

**Book Description:** A production system is defined not by how it behaves when things go right, but by how it behaves when things go wrong. A database that works perfectly under ideal conditions but corrupts data during a power failure is not just useless; it is dangerous. This book is the definitive guide to engineering and verifying the resilience of `go-quadgit`.

We will move beyond testing for correct outputs and begin testing for robust failure modes. This is the practice of **Chaos Engineering**: intentionally and systematically injecting failures into a system to build confidence in its ability to withstand turbulent conditions. You will learn how to build a "Crash Test" harness to validate your database's ACID durability guarantees. We will simulate disk and network failures to ensure every part of the system handles I/O errors gracefully. Finally, you will learn how to use advanced tooling like `toxiproxy` to test the resilience of the distributed `push`/`pull` protocol against real-world network chaos. This book is for the engineer who wants to sleep well at night, knowing their system is prepared for failure.

**Prerequisites:** This is one of the most advanced books in the series. Readers must have a complete understanding of the `go-quadgit` architecture, including its concurrent and distributed features. Expertise in Go, systems programming, shell scripting, and tools like Docker is essential.



## **Part I: The Philosophy of Failure**

*This part establishes the "why" behind intentionally breaking our own system.*

## **Chapter 1: Principles of Fault Injection and Chaos Engineering**
*   **1.1: Why Test for Failure?**
    *   Introduces the core concept: complex systems have emergent behaviors and unknown-unknowns. We can't predict every failure mode, so we must build a system that is resilient to entire classes of failure.
*   **1.2: The Difference Between Correctness and Resilience**
    *   **Correctness:** Does the system produce the right output given a valid input? (Tested in Book 14).
    *   **Resilience:** Does the system enter a safe state when given an *invalid environment* (e.g., no disk space, dead network)? This is the focus of this book.
*   **1.3: Our Toolbox: `kill -9`, Faulty I/O Writers, and Network Proxies**
    *   An overview of the practical tools we will use to inject faults into different layers of the `go-quadgit` stack.



## **Part II: Testing the Core Datastore's Durability**

*This part focuses on a single, critical question: If the server crashes, is my data safe?*

## **Chapter 2: The Crash Test: Verifying ACID Durability**
*   **2.1: Understanding the Role of the Write-Ahead Log (WAL)**
    *   A deep dive into how BadgerDB's WAL provides durability. It explains that a `txn.Commit()` call only returns after the data is safely recorded in the log file on disk. During startup, BadgerDB replays this log to recover from a crash. Our test will validate this exact mechanism.
*   **2.2: Designing the Crash Test Harness**
    *   A step-by-step guide to building the test script. This is not a standard Go test. It's a shell script (`test-crash.sh`) that acts as a test runner.
    *   **Step 1:** The script deletes any old database directory.
    *   **Step 2:** It launches the `TestConcurrentChaos` Go test function (from Book 14) as a background process.
    *   **Step 3:** It waits for a random, short interval (e.g., 5-10 seconds).
    *   **Step 4:** It forcefully terminates the test process with `kill -9 <pid>`. This is crucial as it prevents any graceful shutdown logic from running.
*   **2.3: The Recovery and Verification Phase**
    *   **Step 5:** The script then launches a *new* Go test process, but with a special flag (e.g., `-run-auditor-only`).
    *   **Step 6:** This new process opens the *same database directory* left behind by the crashed process. BadgerDB's `Open()` function will automatically trigger its recovery-from-WAL mechanism.
    *   **Step 7:** The test immediately calls the `assertRepositoryIsInvariants` Auditor function.
*   **2.4: Interpreting the Results**
    *   If the Auditor passes, it is a powerful testament to the system's durability. It proves that even after a catastrophic power failure during a high-concurrency write storm, the database recovered itself to a perfectly consistent and logically correct state, with no corrupted data or broken links. The test script would run this entire cycle in a loop to test many different crash timings.

## **Chapter 3: Simulating I/O Failures**
*   **3.1: The "Faulty Writer": Testing Graceful Error Handling**
    *   Explains that not all failures are crashes. Sometimes, the disk just reports an error. The application must handle this without panicking.
    *   **Implementation:** We will create a custom Go struct that implements the `io.Writer` interface. This `FaultyWriter` will pass through writes normally for a set number of bytes and then begin returning a persistent error (e.g., `io.ErrShortWrite` or a custom "disk full" error).
*   **3.2: Testing the `Backup` Command**
    *   A new integration test is created for `Store.Backup()`.
    *   It passes an instance of the `FaultyWriter` as the destination for the backup stream.
    *   **ACs for the test:**
        1.  The `Backup()` method must not panic.
        2.  It must return a descriptive error that wraps the underlying I/O error.
        3.  It must clean up any partially written, corrupted backup file it may have created.
*   **3.3: Extending the Concept to Other I/O**
    *   Discusses how similar "faulty reader" patterns could be used to test the resilience of the `Restore` and `bulk-load` commands.



## **Part III: Testing the Distributed System's Resilience**

*This part moves beyond a single machine and tests how our networked services behave in the chaotic, unreliable environment of the real world.*

## **Chapter 4: Introducing `toxiproxy`: A Chaos Monkey for Your Network**
*   **4.1: Why `localhost` is a Lie**
    *   Explains that testing network applications on `localhost` gives a false sense of security. The network is instant and never fails. Real networks are slow, laggy, and unreliable.
*   **4.2: Setting up `toxiproxy`**
    *   A practical tutorial on setting up `toxiproxy`. The `go-quadgit-server` listens on port `8081` for Git sync traffic. We configure `toxiproxy` to listen on port `28081` and proxy all traffic to `localhost:8081`. The `go-quadgit` client is then configured to use the proxy port (`28081`) as its remote address.
*   **4.3: Toxics: The Building Blocks of Chaos**
    *   Explains `toxiproxy`'s "toxics"—the specific types of network failure we can inject:
        *   `latency`: Adds a delay to all packets.
        *   `slow_close`: Delays the closing of a TCP connection.
        *   `timeout`: Stops sending data and closes the connection after a period.
        *   `slicer`: Chops data packets into tiny pieces.

## **Chapter 5: Resilience Testing the `push` and `pull` Protocol**
*   **5.1: The Test Harness**
    *   We create a new end-to-end integration test suite that programmatically configures and controls `toxiproxy` via its HTTP API.
*   **5.2: Test Case 1: The High-Latency Push**
    *   **Scenario:** Configure `toxiproxy` to add 500ms of latency to every packet. Run a `go-quadgit push` of a large commit.
    *   **Verification:** The push should eventually succeed. The test measures the total time and verifies it's within an expected (slower) range. This tests the system's tolerance for slow networks.
*   **5.3: Test Case 2: The Mid-Stream Connection Cut**
    *   **Scenario:** This is the most critical test. Begin a `go-quadgit push` of a large packfile. While the client is streaming the `ReceivePack` request body, use the `toxiproxy` API to inject a `timeout` toxic, severing the connection halfway through.
    *   **Verification:**
        1.  The `go-quadgit push` client command must not hang. It must fail with a clear network error message (e.g., "unexpected EOF" or "connection reset").
        2.  The server must not crash. It should gracefully handle the broken stream.
        3.  **Most importantly:** The server's repository state must be unchanged. The Auditor must be run on the server's repository to prove that no partial, corrupted objects from the aborted packfile were written. This validates the atomicity of the `ReceivePack` operation.
*   **5.4: Test Case 3: The Corrupted Packfile**
    *   **Scenario:** Use a custom toxic to randomly flip bits in the data stream.
    *   **Verification:** The server's `PackfileParser` must detect the corruption (e.g., via checksum failures), reject the push with an appropriate error code, and leave the repository unmodified.
