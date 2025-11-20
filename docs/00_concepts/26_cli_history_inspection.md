# History and Inspection

These commands allow you to explore the history of the graph.

## `quad-db log`
*   **Function:** Shows the commit history for the current branch.
*   **Implementation:**
    1.  Reads the commit hash from the current `HEAD`.
    2.  Traverses backward through the commit graph by recursively reading the `parent` hash from each commit object and printing its metadata (hash, author, date, message).

## `quad-db diff <commit1> <commit2>`
*   **Function:** Shows the difference in quads between two commits.
*   **Implementation:**
    1.  Resolves both commit arguments to their respective tree hashes.
    2.  Recursively compares the trees and blobs to generate a set of quads that were added, modified, or deleted between the two states.

## `quad-db show <commit-hash>`
*   **Function:** Shows the metadata and changes for a specific commit.
*   **Implementation:**
    1.  Reads the specified commit object from BadgerDB.
    2.  Prints the commit metadata.
    3.  Performs a `diff` between that commit and its parent to display the changes introduced by that commit.
