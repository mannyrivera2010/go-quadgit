
## **Part II: The Collaborative Layer: Merge Requests**

*This part details the design and implementation of the Merge Request feature, the central hub for team collaboration.*

## **Chapter 4: Designing the Merge Request System**
*   **4.1: The Architectural Separation Principle**
    *   A critical discussion on why Merge Request data is **mutable, non-versioned application state** and must be kept separate from the core, immutable Git object model. This justifies the creation of a new `mrstore`.
*   **4.2: The `MergeRequest` Data Model**
    *   Defines the `MergeRequest` Go struct (`ID`, `Version`, `Status`, `Title`, etc.) and its corresponding key layout in the `app.db` (`app:mr:data:<id>`).
*   **4.3: The `mrstore`: A New CRUD Service**
    *   Implements a new `MergeRequestStore` interface and its concrete implementation, responsible for creating, reading, updating, and listing MR objects.
*   **4.4: Optimistic Locking for Mutable State**
    *   A full implementation guide for using a `Version` counter field on the `MergeRequest` struct, combined with `ETag` and `If-Match` headers in the REST API, to prevent concurrent write conflicts.

## **Chapter 5: The Merge Request API and Workflow**
*   **5.1: Implementing the MR REST Endpoints**
    *   Provides the handler logic for `POST /mergerequests` (create) and `GET /mergerequests` (list).
*   **5.2: The Compositional `GET /mergerequests/:id` Handler**
    *   A detailed walkthrough of how this handler composes a rich response by first fetching the MR object from the `mrstore`, and then calling the `quadstore` to get the associated diff and commit history.
*   **5.3: The `merge` Action: Bringing It All Together**
    *   Implements the `POST /mergerequests/:id/merge` endpoint. This handler demonstrates the full power of the layered system:
        1.  It checks for user authorization (e.g., must be a `maintainer`).
        2.  It performs an `If-Match` check on the *target branch's HEAD*.
        3.  It calls the core `store.Merge()` method to execute the merge.
        4.  If successful, it calls the `mrstore.Update()` method to change the MR's status to "merged".
        5.  All of this is orchestrated as a single logical operation for the user.

