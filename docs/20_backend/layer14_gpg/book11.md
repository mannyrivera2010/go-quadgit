# **Book 11: Establishing Cryptographic Trust**

By the end of this book, the reader will have implemented a complete chain of trust. They will be able to create cryptographically signed commits, automatically verify the integrity of the repository's history, and trace any piece of data back to its verified origin, providing the highest level of assurance and auditability for their knowledge graph.

**Layer III: The Application & Ecosystem**
**Sub-System 6: Ecosystem Integration & Trust**

**Subtitle:** *A Practical Guide to GPG Signing, Signature Verification, and Data Provenance*

**Book Description:** An auditable history is a powerful tool, but in a distributed or high-stakes environment, its value is directly proportional to its trustworthiness. How can we be certain that a commit was created by its stated author? How can we prove that the history of our knowledge has not been tampered with? The answer lies in the application of public-key cryptography.

This book is the definitive guide to building a layer of cryptographic trust into `go-quadgit`. We will implement a complete GPG (GNU Privacy Guard) signing and verification workflow, providing non-repudiable proof of authorship for every commit. You will learn how to integrate these checks into core workflows like `log` and `merge` to create a system that actively defends itself against forgery and tampering. Finally, we will build the ultimate data provenance tool—the `blame` command—allowing you to trace any single fact in your knowledge graph back to its specific, cryptographically verified point of origin.

**Prerequisites:** Readers should be familiar with the core `go-quadgit` architecture, particularly the `Commit` data model and the basics of the CLI. A conceptual understanding of public-key cryptography (public/private keys, digital signatures) is essential.

## **Part I: The GPG Signing and Verification Workflow**

*This part focuses on building the end-to-end mechanics of creating and verifying digital signatures.*

## **Chapter 1: The Foundation of Trust: GPG Signing**
*   **1.1: The Threat Model: Forged Commits and Data Tampering**
    *   A detailed exploration of the security risks in a standard version control system. We'll show exactly how a malicious actor could create a commit with a forged author field and why this is unacceptable for a system of record.
*   **1.2: A Practical Introduction to GPG and Digital Signatures**
    *   Explains the core concepts for a developer audience: what a private key is, what a public key is, and how a digital signature provides the two critical guarantees of **Authenticity** (proof of identity) and **Integrity** (proof of unchanged content).
*   **1.3: What We Sign: The Commit as the Unit of Trust**
    *   A critical architectural discussion. It explains why we sign the `Commit` object itself (containing the tree hash, parent, message, etc.), as this creates an unbreakable link between the author's identity and the entire state of the repository at that moment.
*   **1.4: Extending the `Commit` Data Model**
    *   Revisits the `Commit` struct and details the addition of the `Signature string `json:"signature,omitempty"` field. Explains the importance of `omitempty` for ensuring unsigned commits have a consistent hash.

## **Chapter 2: Implementing the Signing Process**
*   **2.1: The Security Principle: Never Handle Private Keys**
    *   Explains why the `go-quadgit` application must *never* load a user's private GPG key into its own memory. We delegate this sensitive operation to the user's trusted, local `gpg` agent.
*   **2.2: The `sign` Callback: A Clean API Abstraction**
    *   Details the design of the `store.Commit` method, which is modified to accept a `sign func(data []byte) (string, error)` callback. This decouples the core commit logic from the specific implementation of signing.
*   **2.3: Orchestrating `gpg` with `os/exec`**
    *   A full, line-by-line Go implementation of the `createGpgSigner` function. It shows how to use `os/exec` to:
        1.  Start the `gpg --detach-sign --armor` process.
        2.  Pipe the canonical commit data (marshalled JSON) to the process's `stdin`.
        3.  Read the resulting ASCII-armored signature block from `stdout`.
        4.  Handle errors by capturing `stderr`.
*   **2.4: Integrating into the `commit -S` Command**
    *   Shows how the CLI's `commit` command handler checks for the `-S` flag and, if present, creates and passes the `gpg` signing function to the `store.Commit()` method.

## **Chapter 3: Implementing the Verification Process**
*   **3.1: The Verification Challenge: Recreating the Original Data**
    *   Explains the core task of verification: you must re-create the *exact* byte stream that was originally signed.
*   **3.2: The `verify` Helper Function**
    *   Provides the full Go implementation for a function that takes a `Commit` object as input. It will:
        1.  Extract the `Signature` and save it to a temporary file (`signature.asc`).
        2.  Create a temporary copy of the commit, set its `Signature` field to be empty, and marshal it to a second temporary file (`data.json`).
        3.  Execute `gpg --verify signature.asc data.json` using `os/exec`.
*   **3.3: Parsing GPG's Output for a User-Friendly Result**
    *   Shows how to parse the complex `stderr` output from the `gpg` command to reliably determine the signature's status: "Good," "Bad," or "Untrusted Key."

## **Part II: Building Trust into the User Workflow**

*This part focuses on integrating the signing and verification mechanics into the everyday commands users will run, making security a visible and active feature.*

## **Chapter 4: Verifying History with `log` and `merge`**
*   **4.1: At-a-Glance Trust: Implementing `log --show-signature`**
    *   Modifies the `log` command to call the `verify` helper for each commit when the `--show-signature` flag is present.
    *   Provides examples of the final formatted output, showing the GPG status block directly under the commit hash, making the trust level of the history immediately apparent.
*   **4.2: The Security Gate: Automatic Verification in `merge`**
    *   Refactors the `store.Merge()` method to add a new, mandatory first step: verifying the signature of the source commit.
    *   Implements the logic to `panic` and abort the merge on a "Bad signature" and to print a clear warning on an "Untrusted key."
*   **4.3: Securing the Network: Signature Verification on `push`**
    *   Extends the server-side `ReceivePack` handler for `go-quadgit push`. As new commits are unpacked, the server can be configured to verify their signatures before accepting the push, preventing a malicious user from pushing forged commits to a remote server.

## **Chapter 5: The `blame` Command: Ultimate Data Provenance**
*   **5.1: The Goal: Answering "Who Wrote This Fact?"**
    *   Defines the user story and the expected output for `go-quadgit blame <graph-iri>`, showing the commit metadata annotated on every single quad.
*   **5.2: The Scalable Reverse History Walk Algorithm**
    *   A detailed, step-by-step implementation guide for the `blame` algorithm. This is one of the most complex algorithms in the system and is explained with clear diagrams.
        1.  **Get Target Set:** Fetch the quads for the target graph at `HEAD`.
        2.  **Initialize `blame` Map:** Create the data structure to hold the results.
        3.  **The Backward Walk:** Use a breadth-first search to traverse the commit DAG backward.
        4.  **Process Each Commit:** In the loop, use the scalable `StreamDiff` function to find the quads added in the current commit.
        5.  **Annotate and Prune:** If an added quad is in our target set, record the blame and remove it from the set of quads we're still looking for.
        6.  **Termination:** The process ends when the blame map is full.
*   **5.3: Building the API for `blame`**
    *   Implements the `--format=json` flag and the `GET /blame/...` REST endpoint. This shows how to stream the `BlameResult` objects as a JSON array, turning the command into a powerful backend for a visual UI.
