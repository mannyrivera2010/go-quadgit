# Branching and Tags
These commands manage different lines of development or versions of the graph.

## `quad-db branch`
*   **Function:** Lists, creates, or deletes branches.
*   **Implementation:**
    *   `quad-db branch`: Scans BadgerDB for keys with the prefix `ref:head:`.
    *   `quad-db branch <branch-name>`: Creates a new key `ref:head:<branch-name>` and sets its value to the current `HEAD` commit's hash.
    *   `quad-db branch -d <branch-name>`: Deletes the key `ref:head:<branch-name>`.

## `quad-db checkout <branch-name>`
*   **Function:** Switches the `HEAD` to a different branch.
*   **Implementation:**
    1.  Updates the `HEAD` key to contain the reference `ref:head:<branch-name>`. This effectively switches the active line of history for the next commit.

## `quad-db tag <tag-name>`
*   **Function:** Creates a permanent, named pointer to a specific commit.
*   **Implementation:**
    1.  Creates a new key `ref:tag:<tag-name>` and sets its value to the current `HEAD` commit's hash.
