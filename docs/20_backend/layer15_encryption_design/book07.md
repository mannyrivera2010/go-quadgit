# **Book 7: Authentication and Authorization**

**Layer II: The Networked Platform**
**Sub-System 4: Security & Access Control**


**Subtitle:** *Building a Secure, Multi-Tenant Access Control Model for the `go-quadgit` Server*

**Book Description:** A multi-user platform without a robust security model is not a platform; it is a liability. Once `go-quadgit` is exposed to the network, we must have an unshakeable system for verifying who is making a request and what they are permitted to do. This book is the definitive guide to designing and implementing a comprehensive security layer for the `go-quadgit` server.

We will begin by establishing a clear and formal distinction between Authentication (AuthN) and Authorization (AuthZ). You will learn to implement industry-standard authentication strategies, from simple API keys for machine-to-machine communication to a full JWT-based system for user-facing applications. We will then build a powerful and flexible Role-Based Access Control (RBAC) model on top, allowing administrators to define granular permissions on a per-namespace basis. Finally, we will put it all together by implementing this logic as a secure middleware chain in our Go web server and extending it with special-case rules for protecting critical branches.

**Prerequisites:** Readers should be familiar with the `go-quadgit` REST server architecture (Book 3). A solid understanding of web security fundamentals (HTTP headers, tokens) and Go middleware patterns is essential.



## **Part I: The Foundations of API Security**

*This part covers the core concepts and initial setup required for securing our server.*

## **Chapter 1: The Security Model: Authentication vs. Authorization**
*   **1.1: The Two Questions of Access Control**
    *   A deep dive into the formal definitions. **Authentication (AuthN): "Who are you?"** - The process of verifying a claimed identity. **Authorization (AuthZ): "What are you allowed to do?"** - The process of checking the permissions associated with that verified identity.
*   **1.2: The Request Lifecycle in a Secure System**
    *   Presents a clear data flow diagram showing how an incoming API request must first pass through an AuthN middleware gate, and only then through an AuthZ middleware gate, before it is allowed to reach the main application handler. This illustrates the layered defense principle.
*   **1.3: Storing User and Policy Data**
    *   Discusses where security-related data will live. User identity information might be external, but our application-specific data—API key hashes and RBAC policies—will be stored in the `app.db` instance under a dedicated `app:authz:` key prefix.

## **Chapter 2: Implementing Authentication Strategies**
*   **2.1: The Authentication Middleware**
    *   Designs a generic Go middleware that inspects the `Authorization` header of an incoming request. It then delegates to one or more specific authentication strategies. If any strategy succeeds, the user's identity is attached to the request context for later use by the authorization layer. If all fail, it returns `HTTP 401 Unauthorized`.
*   **2.2: Strategy 1: Static API Keys**
    *   **Use Case:** Ideal for CI/CD pipelines, scripts, and other machine-to-machine communication.
    *   **Implementation:** Shows how to securely store SHA-256 hashes of API keys (never plaintext) in `app.db`. The middleware compares the hash of the incoming key with the stored hashes for a fast and secure lookup.
*   **2.3: Strategy 2: JSON Web Tokens (JWT)**
    *   **Use Case:** The standard for user-facing web applications.
    *   **Implementation:**
        1.  Create a new `/auth/login` endpoint that accepts a username/password, authenticates against a user database, and then generates a signed JWT.
        2.  The JWT payload will contain claims like `user_id`, `email`, and an expiration time (`exp`).
        3.  The authentication middleware will validate the JWT's signature using a shared secret and check that it has not expired.
*   **2.4: Strategy 3 (Advanced): Preparing for OAuth 2.0**
    *   A conceptual overview of how our authentication middleware could be extended to support OAuth 2.0, allowing users to "Log in with Google" or an enterprise SSO provider like Okta. This involves handling redirects and validating tokens issued by a third-party identity provider.



## **Part II: The Authorization Engine**

*This part details how we check the permissions of an already authenticated user.*

## **Chapter 3: Designing the Role-Based Access Control (RBAC) Model**
*   **3.1: Defining Roles and Permissions**
    *   Formally defines the standard roles for `go-quadgit`: `reader`, `writer`, `maintainer`, and `admin`.
    *   Creates a comprehensive mapping of which API actions (e.g., `repo:read`, `repo:write`, `mr:create`, `namespace:admin`) are granted to each role. This table serves as our policy specification.
*   **3.2: The Policy Data Model in `app.db`**
    *   Details the key structure `app:authz:policy:<namespace>` and the JSON value that maps user IDs to their role within that namespace.
*   **3.3: An Introduction to Attribute-Based Access Control (ABAC)**
    *   Contrasts RBAC with the more powerful ABAC model. It provides examples of policies that RBAC cannot express (e.g., time-based access) and discusses how a system like Open Policy Agent (OPA) could be integrated in the future for more complex needs. For our implementation, we will stick with the simpler, more common RBAC model.

## **Chapter 4: The Authorization Middleware in Action**
*   **4.1: Building the Middleware**
    *   A full Go implementation of the authorization middleware. It runs *after* the authentication middleware.
*   **4.2: The Authorization Logic**
    *   The middleware performs the following steps:
        1.  Extracts the authenticated `user_id` from the request context.
        2.  Determines the required permission for the target endpoint (e.g., the `POST /merges` handler requires `repo:merge`).
        3.  Determines the resource being accessed (e.g., `namespace: "production"`).
        4.  Loads the RBAC policy for the "production" namespace from `app.db`.
        5.  Finds the user's role in that policy.
        6.  Checks if the user's role is granted the `repo:merge` permission.
        7.  If yes, it calls `c.Next()`. If no, it aborts the request with `HTTP 403 Forbidden`.
*   **4.3: Handling Caching**
    *   Discusses performance optimization. Loading and parsing the policy JSON from the database on every single API request is inefficient. The implementation will use an in-memory, time-based cache (like a simple map with a mutex) to store policies for a short duration (e.g., 1 minute) to reduce database load.



## **Part III: Advanced Authorization Policies**

*This part covers special-case authorization rules that go beyond the simple RBAC model.*

## **Chapter 5: Implementing Branch Protection Rules**
*   **5.1: The Need for Granular Branch Control**
    *   Explains why a global `repo:write` permission is too broad. Important branches like `main` need an extra layer of protection.
*   **5.2: The `BranchProtection` Data Model and API**
    *   Details the `BranchProtection` struct (`allow_direct_push`, `require_merge_request`, etc.) and its storage key (`app:bpr:<namespace>:<branch_name>`).
    *   Defines the REST API (`PUT /ns/:ns/branches/:br/protection`) for administrators to manage these rules dynamically, emphasizing why this must be an API and not a static config file.
*   **5.3: Enhancing the Authorization Middleware**
    *   Shows how to extend the authorization middleware. After checking the user's role, if the target is a branch, the middleware performs a second check to see if a protection rule exists.
    *   It then applies the specific rules. For a `push` operation, it checks `allow_direct_push`. For a `merge` operation via the MR API, it checks `require_merge_request`. This demonstrates how to compose multiple layers of policy.

By the end of this book, the `go-quadgit` server will be transformed from an open platform into a secure fortress. It will have a complete, production-ready security layer capable of identifying every user, enforcing granular, per-namespace permissions, and protecting its most critical assets with specific, configurable rules.