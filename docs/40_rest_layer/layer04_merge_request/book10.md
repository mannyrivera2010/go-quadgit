## **Book 10: The Merge Request Workflow**

**Layer III: The Application & Ecosystem**

**Sub-System 5: Collaborative Application Features*

**Subtitle:** *Building a Collaborative Review and Integration System for Knowledge Graphs*

**Book Description:** A version control system provides the mechanism for change, but a true collaborative platform also provides a framework for the *process* of change. The most successful and widely adopted process for team-based development is the **Merge Request** (or Pull Request). It is a formal proposal to integrate changes, creating a space for review, discussion, and automated validation before the final merge is executed.

This book is the definitive guide to building a complete Merge Request system as a high-level application layer on top of the `go-quadgit` platform. You will learn the critical architectural principle of separating mutable application state from immutable versioned history. We will design and build the data model, the REST API, and the core business logic for managing the entire lifecycle of a Merge Request. By the end of this book, you will have transformed `go-quadgit` from a raw versioning engine into a sophisticated platform for collaborative knowledge engineering.

**Prerequisites:** Readers must have a complete understanding of the `go-quadgit` architecture, particularly the Core API (`quadstore.Store`) and the REST server with its multi-instance datastore. Deep familiarity with RESTful API design, Go, and state management in web applications is required.

### **Part I: The Architectural Foundation**

*This part establishes the "why" and "how" of building a stateful application layer on top of our versioned datastore.*

#### **Chapter 1: Architectural Separation: The Core Principle**
*   **1.1: What is a Merge Request? A Proposal, Not a Commit**
    *   Formally defines a Merge Request as a long-lived, stateful object representing a *request to merge*, distinct from the `merge` operation itself. Its state includes `open`, `merged`, `closed`, reviewers, comments, and CI status.
*   **1.2: Why Merge Requests Must Be Non-Versioned**
    *   This is the most critical design decision. The chapter provides a detailed argument against storing MR metadata within the Git-like object model.
    *   **Argument 1: The Mutable State Problem.** The Git model is optimized for immutable snapshots. An MR's state changes constantly. Using commits to track these changes would be inefficient and conceptually wrong.
    *   **Argument 2: The History Pollution Problem.** Committing every comment or status change would flood the `go-quadgit log` with thousands of irrelevant entries, making it impossible to see the actual history of the knowledge graph.
*   **1.3: Our Strategy: A Separate Application Data Store**
    *   Introduces the solution: MR data will be treated as simple, mutable Key-Value data. It will be stored in the dedicated `app.db` instance, completely separate from the `history.db` and `index.db`. This keeps the core versioning engine pure.

#### **Chapter 2: Designing the Merge Request Data Model and Store**
*   **2.1: The `MergeRequest` Go Struct**
    *   Defines the complete data structure, including all necessary fields: `ID`, `Version` (for optimistic locking), `Title`, `Description`, `Status`, `AuthorID`, `SourceBranch`, `TargetBranch`, `CreatedAt`, `MergedAt`, `MergeCommitHash`, etc.
*   **2.2: The Key Layout in `app.db`**
    *   Details the simple and efficient key structure:
        *   `app:mr:data:<id>`: Stores the serialized JSON of the `MergeRequest` object.
        *   `app:mr:sequence`: The atomic counter used by BadgerDB's `Sequence.Next()` to generate unique MR IDs.
*   **2.3: The `mrstore`: A New CRUD Service Interface**
    *   Defines a new `mrstore.MergeRequestStore` interface with methods like `Create`, `GetByID`, `Update`, and `List`. This establishes a clean architectural boundary for all MR-related database operations.
*   **2.4: Implementing the `mrstore`**
    *   Provides the concrete Go implementation of the store, showing how to use the `app.db` BadgerDB instance to perform the CRUD operations.

### **Part II: The REST API and User Workflow**

*This part focuses on building the user-facing API and orchestrating the interactions between the new `mrstore` and the existing `quadstore`.*

#### **Chapter 3: The Merge Request API Endpoints**
*   **3.1: Creating and Listing MRs**
    *   A full implementation of the `POST /mergerequests` handler for creating a new MR.
    *   A full implementation of the `GET /mergerequests` handler, including logic for parsing query parameters to filter by `status`, `author`, etc.
*   **3.2: The Compositional `GET /mergerequests/:id` Handler**
    *   This is a masterclass in composing services. The handler's implementation is detailed step-by-step:
        1.  Call `mrstore.GetByID()` to fetch the core MR object.
        2.  Call `quadstore.Store.Log()` to fetch the list of commits for the "Commits" tab.
        3.  Call `quadstore.Store.Diff()` to generate the list of quad changes for the "Changes" tab.
        4.  Combine these three pieces of data into a single, rich JSON response for the frontend.
*   **3.3: Updating an MR**
    *   Implements the `PUT /mergerequests/:id` endpoint for editing the title or description. This chapter will heavily feature the optimistic locking mechanism.

#### **Chapter 4: The Merge Action: A Secure, Multi-Step Orchestration**
*   **4.1: The `POST /mergerequests/:id/merge` Endpoint**
    *   This is the capstone workflow. The implementation of this handler ties together almost every concept in the entire system.
*   **4.2: Step 1: Pre-Flight Checks and Authorization**
    *   The handler first fetches the MR object and the user's identity. It checks: Is the MR status `open`? Does the user's role (e.g., `maintainer`) grant them `repo:merge` permission on the target branch?
*   **4.3: Step 2: Enforcing Branch Protections**
    *   The handler then loads any `BranchProtection` rules for the target branch. It checks: Does this MR meet the `required_approvals` count? Have all required CI/CD status checks passed?
*   **4.4: Step 3: The Optimistic Lock**
    *   The handler validates the crucial `If-Match` header from the client against the current `HEAD` of the target branch, preventing merges into a stale state.
*   **4.5: Step 4: Executing the Core Merge**
    *   If all checks pass, it finally calls the low-level `quadstore.Store.Merge()` method to perform the actual database operation.
*   **4.6: Step 5: Updating Application State**
    *   Upon a successful merge, the handler calls `mrstore.Update()` to atomically change the MR's status to `merged`, recording the new merge commit hash and the user who performed the merge.

### **Part III: Advanced Features and User Experience**

*This part covers features that make the Merge Request experience truly powerful and user-friendly.*

#### **Chapter 5: Implementing a Commenting System**
*   **5.1: The `Comment` Data Model**
    *   Defines a new non-versioned data model for comments, linked to an MR ID.
*   **5.2: The `POST /mergerequests/:id/comments` Endpoint**
    *   Implements the API for adding a new comment to an MR.
*   **5.3: Inline Comments: Linking Comments to Diffs**
    *   An advanced implementation showing how a comment can be optionally linked to a specific quad (or its hash) within the MR's diff, allowing for line-by-line code review.

#### **Chapter 6: Real-Time Updates with WebSockets**
*   **6.1: The Real-Time Experience**
    *   Explains why polling is inefficient for a collaborative tool.
*   **6.2: Broadcasting MR Events**
    *   The server is modified to broadcast WebSocket events after key actions:
        *   When a new comment is posted.
        *   When a user approves the MR.
        *   When the source branch is updated with new commits.
        *   When the MR is merged or closed.
*   **6.3: The Frontend Implementation**
    *   Shows how the React UI can listen for these events and update the MR page in real-time without requiring the user to refresh the page.

By the end of this book, the reader will have built a complete, production-ready Merge Request system. They will have mastered the art of layering a stateful, mutable application on top of an immutable, version-controlled core, creating a platform that is safe, powerful, and deeply collaborative.
