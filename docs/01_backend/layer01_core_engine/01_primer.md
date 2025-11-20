# Book 1: The `go-quadgit` Primer
**Layer I: The Core Engine**

**Sub-System 1: Foundational Concepts & Interface**

**Subtitle:** *From Zero to Your First Versioned Knowledge Graph*

**Book Description:** This book is the definitive starting point for anyone new to `go-quadgit`. It is a practical, hands-on guide that establishes the core philosophy of versioning knowledge graphs and walks you through every step of creating, managing, and inspecting a local repository. We will begin with the high-level concepts, move to a hands-on tutorial where you'll make your first commit within minutes, and then dive deeper into the underlying data model and the full command suite. By the end of this book, you will have mastered the single-user workflow and gained a solid conceptual foundation for the entire `go-quadgit` platform.

## **Part I: Getting Started**

*This part is designed to get a user up and running as quickly as possible, providing immediate value and a successful first experience.*

## **Chapter 1: Introduction to Versioned Graphs**
*   **1.1: The Problem: Data's Missing Past**
    *   A detailed exploration of the shortcomings of traditional databases for audit and reproducibility. We'll use concrete examples from finance, science, and compliance to illustrate *why* answering "what was true then?" is a critical, unsolved problem.

*   **1.2: The Solution: The Git Mental Model for Knowledge**
    *   A deep dive into the core analogy: Repository, File (Named Graph), Line of Code (Quad), Commit, Branch, and Merge. This section builds the intuitive bridge between source code management and knowledge graph management.

*   **1.3: Introducing `go-quadgit`: Goals and Design Principles**
    *   A mission statement for the project. We'll formally introduce the guiding principles: Scalability, Performance, Semantic Integrity, and Usability. This sets expectations for what the system is designed to do well.

*   **1.4: The Technology Stack: Why Go and BadgerDB?**
    *   A justification of our core technology choices, aimed at a technical audience. We'll discuss the benefits of Go's concurrency model and static binaries, and why BadgerDB's pure Go nature, performance, and transactional guarantees make it the ideal foundation.

## **Chapter 2: Installation and Initial Configuration**
*   **2.1: Installing the `go-quadgit` CLI**
    *   Step-by-step instructions for different platforms (macOS, Linux, Windows). This will cover installation via pre-compiled binaries from the releases page, using `go install`, and building from source.

*   **2.2: Setting Up Your Identity**
    *   Explains the importance of the author field in commits. It walks the user through setting up their name and email via the `go-quadgit config --global user.name "..."` command, explaining where the `~/.quadgit/config` file is stored and what it does.

*   **2.3: Shell Completion**
    *   Instructions for enabling command-line completion for Bash, Zsh, and PowerShell to improve usability.

## **Chapter 3: Your First Repository: A Hands-On Tutorial**
*   **3.1: Initialization: `init`**
    *   Walks the user through running `go-quadgit init` and explains the structure of the newly created `.quadgit/` directory.

*   **3.2: Preparing Your Data: The `add` Command**
    *   Explains how to format data in an N-Quads file. This section will introduce the simple syntax for marking quads for deletion (e.g., prefixing with `D `) versus addition (no prefix).

*   **3.3: Making History: Your First `commit`**
    *   Guides the user through running `go-quadgit add data.nq` followed by `go-quadgit commit -m "Initial dataset"`. It explains what a good commit message looks like.

*   **3.4: Seeing the Results: `log` and `status`**
    *   Shows the user how to run `go-quadgit log` to see the commit they just created. It explains each part of the log output (hash, author, date, message). It also shows how `go-quadgit status` now reports a clean state.

*   **3.5: Making and Inspecting a Change: `diff` and `show`**
    *   This completes the core loop. The user is guided to make a change to their data file, `add` it, and `commit` it again. Then, they will use `go-quadgit diff HEAD~1 HEAD` to see the exact changes they made, and `go-quadgit show HEAD` to see the full commit details plus the patch.


## **Part II: Core Concepts in Detail**

*This part moves from tutorial to reference, providing a deeper, more technical explanation of the concepts introduced in Part I. It's for the user who now asks, "How does this actually work?"*

## **Chapter 4: The `go-quadgit` Data Model Under the Hood**
*   **4.1: Content-Addressable Storage: Blobs and Trees**
    *   A detailed look at how a named graph's content is hashed to create a `Blob` and how a `Tree` object acts as a manifest, pointing to these blobs. Includes diagrams of the object relationships.

*   **4.2: The Immutable Ledger: Commits and the DAG**
    *   Explains how `Commit` objects link to a `Tree` and their `parents`, forming a Directed Acyclic Graph (DAG). It emphasizes the cryptographic integrity of this chain.

*   **4.3: Mutable Pointers: Branches, Tags, and `HEAD`**
    *   Contrasts the immutable objects with the mutable references. Explains that a branch is just a simple named pointer that moves with each commit.

*   **4.4: On-Disk Representation: The Key-Value Layout**
    *   This is a crucial technical section. It shows the actual key structures used in BadgerDB: `obj:<hash>`, `ref:head:<name>`, `ref:tag:<name>`, and `HEAD`. This makes the abstract data model concrete.

*   **4.5: Anatomy of a Commit: A Detailed Walkthrough**
    *   Revisits the commit process from Chapter 3, but this time from the system's perspective. It shows exactly which keys are read from and written to the database at each step, tying all the data model concepts together.

## **Chapter 5: The CLI Command Reference**
*   **5.1: Repository Management Commands**
    *   Detailed reference page for `init` and `status`.

*   **5.2: State Modification Commands**
    *   Detailed reference page for `add` and `commit`, including all their flags (e.g., `-S` for signing, even if it's implemented later).

*   **5.3: History Inspection Commands**
    *   Detailed reference page for `log`, `diff`, and `show`, covering various argument formats (e.g., `HEAD~3`, branch names, commit hashes) and flags (`--stat`, `--show-signature`).

*   **5.4: Branching and Tagging Commands**
    *   Detailed reference page for `branch`, `checkout`, and `tag`, including flags for deleting, renaming, and listing.



## **Part III: The Path Forward**

*This part concludes the primer and provides a bridge to the more advanced topics covered in the other books.*

## **Chapter 6: Beyond the Local Repository**
*   **6.1: A Glimpse into Merging**
    *   Briefly explains the concept of merging and why it's necessary for collaboration. Points the reader to **Book 2: Engineering for Scale and Performance** for the full implementation details.

*   **6.2: Preparing for Scale**
    *   Touches on the memory limitations of the simple `diff` algorithm implemented in this book and explains that the true scalable algorithms are covered in **Book 2**.

*   **_6.3: The Journey to a Networked Service**
    *   Paints the vision for the REST server and distributed workflows (`push`/`pull`). It explains that the clean API built for the CLI is what makes this possible and directs the reader to **Layer II: The Networked Platform** for the full story.

*   **6.4: Where to Go Next: A Guide to the `go-quadgit` Library**
    *   A summary of the other books in the series, explaining what each one covers and who it's for. This provides a clear roadmap for a user wanting to continue their learning journey.