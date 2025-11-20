# Querying
This is the primary way users will retrieve data from the graph at a specific point in time.

## `quad-db query -v <version> "SELECT ..."`
*   **Function:** Executes a SPARQL-like query against the graph as it existed at a specific version (commit hash, branch, or tag).
*   **Implementation:**
    1.  **Version Resolution:** Resolves the `<version>` argument to a specific commit hash, and from there to a root tree hash.
    2.  **State Reconstruction:** Traverses the tree and blob objects to reconstruct the set of all quads that existed at that version.
    3.  **Query Execution:** Uses the reconstructed quad set and the powerful SPOG/POSG/etc. indexes in BadgerDB to efficiently execute the query and return the results.
