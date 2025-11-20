# Repository Management

These commands manage the database repository itself.

## `quad-db init`
*   **Function:** Initializes a new, empty versioned quad store.
*   **Implementation:**
    1.  Creates a new directory for the BadgerDB database.
    2.  Opens a connection to the new database.
    3.  Creates the initial commit, representing an empty state. This involves creating an empty tree object and a root commit object with no parent.
    4.  Creates the `main` branch by creating a key `ref:head:main` that points to the initial commit's hash.
    5.  Creates the `HEAD` key with the value `ref:head:main`, making `main` the active branch.

