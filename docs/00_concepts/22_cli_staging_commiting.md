# Staging and Committing Changes
This set of commands allows users to add, remove, and save changes to the graph. To manage this, we introduce the concept of a "staging area" or "index," which, just like in Git, is a key in BadgerDB (`index`) that lists pending changes (quad additions/deletions).

## `quad-db add <file.nq>`
*   **Function:** Adds the contents of an N-Quads file to the staging area for the next commit.
*   **Implementation:**
    1.  Parses the N-Quads file.
    2.  For each quad in the file, it adds an "add" operation entry to the `index` key in BadgerDB. This entry would contain the full quad data.

## `quad-db rm <file.nq>`
*   **Function:** Stages the deletion of quads specified in a file.
*   **Implementation:**
    1.  Parses the N-Quads file.
    2.  For each quad, it adds a "delete" operation entry to the `index` key.

## `quad-db status`
*   **Function:** Shows the current state of the working directory and staging area.
*   **Implementation:**
    1.  Reads the `index` key to list all staged additions and deletions.
    2.  It could also compare the current working set of quads (if maintained outside the DB) against the `HEAD` commit to show unstaged changes.

## `quad-db commit -m "A descriptive message"`
*   **Function:** Records the staged changes to the repository. This is the core versioning operation.
*   **Implementation:**
    1.  Reads the set of quad changes from the `index`.
    2.  Creates new "blob" objects for these changes.
    3.  Creates a new "tree" object by taking the tree from the parent commit (`HEAD`) and applying the changes from the blobs.
    4.  Creates a new "commit" object containing the hash of the new tree, the parent commit's hash (read from `HEAD`), the author's metadata, and the commit message.
    5.  Updates the current branch reference (e.g., `ref:head:main`) to point to the new commit's hash.
    6.  Clears the `index`.
