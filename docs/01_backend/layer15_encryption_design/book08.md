# **Book 8: Securing Data at Rest with Encryption**

This book provides the complete blueprint for building a multi-tenant encryption layer that is secure, user-friendly, and operationally sound, adding the final pillar of confidentiality to the `go-quadgit` platform.

**Layer II: The Networked Platform**

**Sub-System 4: Security & Access Control**

**Subtitle:** *Implementing Per-Namespace, Password-Based Encryption with a Cryptographic Key Hierarchy*

**Book Description:** In the previous book, we built a powerful system for authenticating users and authorizing their actions. But what about the data itself? By default, all `Commit`, `Tree`, and `Blob` objects are stored as plaintext JSON in our database. Anyone who gains access to the server's filesystem can read the entire history of every knowledge graph. For any organization handling sensitive, private, or regulated data, this is an unacceptable risk.

This book provides a masterclass in building a zero-trust, end-to-end encryption layer for `go-quadgit`. We will move beyond a simple, single master key and design a sophisticated **cryptographic key hierarchy** that provides per-namespace confidentiality. You will learn the theory and practice of using a Key Encryption Key (KEK) to protect multiple Data Encryption Keys (DEKs), and how to use strong Key Derivation Functions (KDFs) like Argon2id to safely turn user-provided passwords into powerful encryption keys.

Through detailed Go implementations, we will build a secure "unlock" workflow for the REST API, manage the lifecycle of cryptographic keys, and establish operational procedures for password rotation and secure data migration. After completing this book, you will have transformed `go-quadgit` into a hardened platform where data confidentiality is a core, verifiable, and non-negotiable feature.

**Prerequisites:** This is a highly advanced book. Readers must have a complete understanding of the `go-quadgit` architecture, a strong command of Go, and a solid conceptual understanding of modern symmetric-key cryptography (AES), authenticated encryption (AEAD), and the purpose of hashing and salting.

## **Part I: The Cryptographic Design**

*This part lays out the theoretical foundation for our multi-tenant encryption model.*

## **Chapter 1: The Limits of a Single Master Key**
*   **1.1: A Review of the Simple Encryption Model**
    *   Briefly revisits a naive implementation where a single `QUADGIT_MASTER_KEY` encrypts all data.
*   **1.2: The "Single Point of Failure" Problem**
    *   Details the threat model: if the single master key is compromised, *all* namespaces are compromised instantly.
*   **1.3: The Need for Tenant Isolation**
    *   Explains why, in a multi-tenant environment, the security of one tenant's data must not be dependent on the security practices of another. A password leak for "Project A" should have zero impact on "Project B".
*   **1.4: Our Goal: The KEK/DEK Model**
    *   Formally introduces the "envelope encryption" pattern we will build: a master Key Encryption Key (KEK) that "wraps" or encrypts multiple, per-namespace Data Encryption Keys (DEKs).

## **Chapter 2: A Practical Guide to the Cryptographic Primitives**
*   **2.1: AEAD Ciphers: Why AES-256-GCM is the Right Choice**
    *   Explains Authenticated Encryption with Associated Data (AEAD). It details why the GCM mode is critical, as it provides not just confidentiality but also integrity and authenticity, preventing an attacker from tampering with ciphertext on disk.
*   **2.2: The Problem with Passwords: Entropy and Brute-Forcing**
    *   Explains why a user's password (e.g., "password123") is not a cryptographically secure key. It has low entropy and is vulnerable to dictionary attacks.
*   **2.3: Key Derivation Functions (KDFs): Turning Passwords into Keys**
    *   Introduces KDFs as the solution. Explains the purpose of a **salt** (to make every hash unique) and the concept of **memory-hard** functions that resist brute-force attacks on GPUs.
*   **2.4: A Deep Dive into Argon2id**
    *   Details why Argon2id, the winner of the Password Hashing Competition, is the modern, recommended choice for a KDF. Explains its tuning parameters (memory, iterations, parallelism).

## **Chapter 3: Designing the Key Hierarchy and "Unlock" Workflow**
*   **3.1: The Lifecycle of a Namespace Key**
    *   A detailed diagram illustrating the full lifecycle:
        1.  **Creation:** A user provides a password. The system generates a random salt. `Password + Salt -> KDF -> DEK`.
        2.  **Encryption:** The system uses the master KEK to encrypt the DEK. `DEK + KEK -> EncryptedDEK`.
        3.  **Storage:** The system stores the `EncryptedDEK` and the `Salt` in the database.
        4.  **Unlocking:** A user provides a password. The system re-derives a key, decrypts the stored DEK, and compares them to verify the password is correct.
*   **3.2: The New Data Model in `app.db`**
    *   Defines the `app:ns_crypto:<namespace>` key and the JSON structure of its value: `{"salt": "...", "encrypted_dek": "..."}`.
*   **3.3: The Server's State: Caching Unlocked Keys**
    *   Explains the need for a thread-safe, in-memory map on the server to hold the plaintext DEKs for namespaces that users have recently unlocked, avoiding the need to re-derive the key on every single request.

## **Part II: Implementation**

*This part provides the practical Go code for building the entire encryption system.*

## **Chapter 4: Refactoring the Datastore and Core API**
*   **4.1: The Global KEK and the Cipher Cache**
    *   Shows the Go code for initializing the master `cipher.AEAD` from the `QUADGIT_MASTER_KEY` environment variable at server startup.
    *   Implements the thread-safe map (`map[string]cipher.AEAD` with a `sync.RWMutex`) for caching the unlocked per-namespace ciphers.
*   **4.2: The `getCipher()` Method**
    *   Implements a helper method on the `Repository` that looks up the correct AEAD cipher for the current namespace context from the global cache. It returns an error if the namespace is locked.
*   **4.3: Updating `writeObject` and `readObject`**
    *   A full, line-by-line refactoring of these two critical methods. They are modified to call `getCipher()` and, if a cipher is returned, use it to transparently encrypt or decrypt the data. The logic must gracefully handle unencrypted namespaces by simply passing the plaintext through.

## **Chapter 5: The User-Facing Workflow: Creation and Unlocking**
*   **5.1: The `ns create --password` Command**
    *   The full implementation of the CLI command. It includes using Go's `golang.org/x/crypto/argon2` package to derive the DEK, and then using the global KEK to encrypt it before saving it to the new namespace's crypto key.
*   **5.2: The `POST /unlock` REST Endpoint**
    *   Implements the HTTP handler that performs the full unlock workflow described in Chapter 3.
*   **5.3: Session Management with JWT**
    *   Shows how, upon a successful unlock, the server can issue a specialized, short-lived JWT to the client. This token might contain a claim like `"unlocked_ns": ["proj_alpha"]`. The authentication middleware can then use this token to know that the user has permission to access the cached cipher without needing the password again for subsequent requests.

## **Chapter 6: Building a Secure Frontend Experience**
*   **6.1: Handling the `HTTP 423 Locked` Response**
    *   Shows the frontend (React) logic for catching a `423` status code.
*   **6.2: The "Unlock Namespace" Modal**
    *   Builds the React component for a modal dialog that prompts the user for their namespace password.
*   **6.3: Securely Storing the Session Token**
    *   Discusses best practices for storing the returned JWT on the client-side (e.g., in memory or in a secure, `HttpOnly` cookie).

## **Part III: Operations and Advanced Security**

*This part covers the day-to-day management and long-term security of the encrypted system.*

## **Chapter 7: Password and Key Lifecycle Management**
*   **7.1: The Password Rotation Workflow**
    *   Implements the `POST /namespaces/:name/change-password` endpoint. The user must provide their old password and a new one. The logic involves unlocking the DEK with the old password and then re-encrypting it with a key derived from the new password.
*   **7.2: Master Key (KEK) Rotation**
    *   Details the high-security operational procedure for rotating the master `QUADGIT_MASTER_KEY`. This involves a controlled maintenance window where the server is run in a special "re-keying" mode to decrypt all stored DEKs with the old KEK and re-encrypt them with a new KEK.

## **Chapter 8: Data Migration and Auditing**
*   **8.1: Encrypting an Existing Namespace**
    *   Provides a step-by-step administrator's guide for migrating an unencrypted namespace. This involves creating a new, encrypted namespace and using the `go-quadgit cp` command to transfer the data, which will be encrypted on write.
*   **8.2: The `go-quadgit security-scan` Command**
    *   Builds a new administrative tool that can be run to audit the security posture of the entire instance. It iterates through all namespaces and reports on their encryption status, the strength of their KDF parameters, and other security-related metadata.

