
## **Part III: The Semantic Intelligence Layer**

*This part introduces the features that make `go-quadgit` truly "graph-aware."*

## **Chapter 6: Schema-Aware Merging**
*   **6.1: From Syntactic to Semantic Conflicts**
    *   Explains the limitations of a standard diff. A merge might be textually clean but logically invalid (e.g., creating two birth dates for one person).
*   **6.2: The Schema Index: Pre-computing Constraints for Performance**
    *   Details the implementation of the `index:schema:...` keys. This involves a process that runs during `commit` to parse the ontology file, extract key constraints (functional, cardinality, etc.), and save them into a fast, queryable index in `index.db`.
*   **6.3: A Comprehensive List of Semantic Checks**
    *   Provides a checklist and implementation guide for various OWL/RDFS constraints:
        *   `owl:FunctionalProperty`
        *   `owl:InverseFunctionalProperty`
        *   `owl:SymmetricProperty`
        *   `owl:maxCardinality`
        *   `rdfs:domain` and `rdfs:range`
        *   `owl:disjointWith`

## **Chapter 7: Implementing the Semantic Merge**
*   **7.1: Modifying the `Store.Merge()` Method**
    *   Refactors the `merge` method to include a new "Semantic Validation" step after the three-way diff is calculated but before the final merge commit is created.
*   **7.2: The Validation Logic**
    *   Shows how the validation step simulates the "proposed" merged state and then runs a series of fast checks against the Schema Index. For each proposed quad, it can quickly look up whether its predicate is functional, what its domain/range is, etc.
*   **7.3: Handling Schema Evolution: The Two-Phase Merge**
    *   A deep dive into the most complex scenario: what if the schema itself is what changed in conflicting ways on two branches? This section details the "Schema Reconciliation First" strategy, where schema changes are merged and validated first, and only then is the data validated against the *newly reconciled* schema.

