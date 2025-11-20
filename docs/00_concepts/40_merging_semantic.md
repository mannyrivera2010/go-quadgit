
## The Core Principle: From Syntactic to Semantic Conflicts

A standard merge tool operates **syntactically**. It sees quads as unique strings of text. If branch `A` adds `<S> <P> <O1> <G>` and branch `B` adds `<S> <P> <O2> <G>`, a syntactic merge tool has no inherent reason to see this as a conflict. It would simply add both quads to the final state.

A **schema-aware** merge tool operates **semantically**. It understands that the predicate `<P>` (e.g., `<ex:hasAge>`) might have constraints defined in an ontology (the schema). If the schema states that `<ex:hasAge>` is an `owl:FunctionalProperty`, the tool knows that a subject can only have *one* age. Therefore, the presence of two different ages for the same subject is a logical contradiction and must be flagged as a conflict.

## Prerequisite: The Schema is Part of the Versioned Data

For the merge tool to be "aware" of the schema, the schema itself must be accessible within the database. The standard practice is to store the ontology (the RDFS/OWL file defining the classes and properties) in a dedicated named graph.

*   **Example:** All schema-defining triples could be stored in the graph `<urn:quad-db:schema>`.

When a merge operation begins, the first step for the merge tool is to load and parse the triples from this schema graph into an efficient, in-memory lookup structure.

## How the Merge Algorithm Changes

The schema-aware merge process enhances the standard three-way merge algorithm:

1.  **Find Common Ancestor:** (No change) Find the merge base between the two branches.
2.  **Load Schema:** (New Step) Load the schema definition from the `urn:quad-db:schema` graph of the **target branch**. This ensures that changes are validated against the most current version of the ontology.
3.  **Calculate Diffs:** (No change) Generate the set of quad additions and deletions for both branches since the common ancestor.
4.  **Apply Enhanced Conflict Detection:** (Major Change) In addition to direct quad conflicts (add vs. delete), the tool now applies a series of semantic validation rules based on the loaded schema.

## Key Schema Constructs and Their Role in Conflict Detection

Here are the specific OWL and RDFS constructs the merge tool would use:

## 1. Functional Properties (`owl:FunctionalProperty`)

This is the most critical and common check.
*   **Rule:** A predicate declared as an `owl:FunctionalProperty` can only have one unique value (object) for a given subject.
*   **Merge Scenario:**
    *   **Schema:** `<ex:hasSSN> rdf:type owl:FunctionalProperty .`
    *   **`main` branch:** Adds `<person:Bob> <ex:hasSSN> "123" .`
    *   **`feature` branch:** Adds `<person:Bob> <ex:hasSSN> "456" .`
*   **Schema-Aware Outcome:** **High-priority conflict.** The system immediately flags this as a logical impossibility. The conflict report would explicitly state: *"Conflict: The predicate <ex:hasSSN> is functional, but branches provide conflicting values ('123' and '456') for the subject <person:Bob>."*

## 2. Cardinality Constraints (`owl:maxCardinality`)

This is a more general version of functional properties.
*   **Rule:** Constrains how many values a predicate can have for a subject. `owl:maxCardinality "1"` is equivalent to a functional property.
*   **Merge Scenario:**
    *   **Schema:** A constraint on a class, stating that instances have a `ex:hasChild` property with a maximum cardinality of 2.
    *   **State in Common Ancestor:** `<person:Carol>` already has one child: `<person:Carol> <ex:hasChild> <person:David> .`
    *   **`main` branch:** Adds `<person:Carol> <ex:hasChild> <person:Eve> .`
    *   **`feature` branch:** Adds `<person:Carol> <ex:hasChild> <person:Frank> .`
*   **Schema-Aware Outcome:** **Conflict.** Individually, each branch's change is valid. However, merging them would result in Carol having three children, violating the `maxCardinality` of 2. The merge tool, by simulating the final state, can detect this violation.

## 3. Class Disjointness (`owl:disjointWith`)

This prevents an individual from belonging to two mutually exclusive classes.
*   **Rule:** If `ClassA` is `owl:disjointWith` `ClassB`, an individual cannot be an instance of both.
*   **Merge Scenario:**
    *   **Schema:** `<class:Child> owl:disjointWith <class:Adult> .`
    *   **`main` branch:** Asserts `<person:George> rdf:type <class:Child> .`
    *   **`feature` branch:** Asserts `<person:George> rdf:type <class:Adult> .`
*   **Schema-Aware Outcome:** **Conflict.** A naive merge would result in George being both a child and an adult, a logical inconsistency. The schema-aware tool identifies this as a direct conflict based on the disjointness axiom.

## 4. Domain and Range Validation (`rdfs:domain`, `rdfs:range`)

This helps maintain data quality and type safety.
*   **Rule:** `rdfs:domain` specifies the class of the subject; `rdfs:range` specifies the class of the object.
*   **Merge Scenario:**
    *   **Schema:** `<ex:hasAge> rdfs:range xsd:integer .`
    *   **`main` branch:** Adds `<person:Alice> <ex:hasAge> "30"^^xsd:integer .` (Correct)
    *   **`feature` branch:** Adds `<person:Alice> <ex:hasAge> "thirty"^^xsd:string .` (Incorrect type)
*   **Schema-Aware Outcome:** This could be treated as either a **hard conflict** or a **high-priority warning**. The system knows that "thirty" is not in the value space of `xsd:integer`. The conflict report can be very specific: *"Warning/Conflict: The value 'thirty' for predicate <ex:hasAge> violates its defined range of xsd:integer."*

## Post-Merge Consistency Check (An Even More Advanced Step)

Some schema rules don't create direct conflicts on individual quads but can render the entire merged graph logically inconsistent. These are best handled by an optional consistency check after the initial merge logic runs.

*   **Example with `owl:SymmetricProperty`:**
    *   **Schema:** `<foaf:knows> rdf:type owl:SymmetricProperty .`
    *   **Merge Result:** The merge peacefully adds the quad `<person:Alice> <foaf:knows> <person:Bob> .` However, the corresponding quad `<person:Bob> <foaf:knows> <person:Alice> .` does not exist.
    *   **Outcome:** The merge itself doesn't have a "conflict," but the resulting state is inconsistent with the schema. A post-merge check could run a lightweight reasoner to find such violations and warn the user: *"Post-merge warning: The merged graph is inconsistent. The predicate <foaf:knows> is symmetric, but a corresponding statement for <person:Bob> knowing <person:Alice> is missing."*

## Benefits and Challenges

## Benefits:
1.  **Higher Accuracy:** Conflicts are semantically meaningful, reducing false positives (where changes are syntactically different but semantically compatible) and catching false negatives (where changes seem compatible but are logically contradictory).
2.  **Improved User Experience:** Conflict messages are far more descriptive and actionable, guiding the user on *why* a conflict occurred (e.g., "functional property violation").
3.  **Guaranteed Data Integrity:** The system actively prevents the creation of commits that are known to be logically inconsistent with the governing schema.

## Challenges:
1.  **Performance:** Loading a large schema and performing these checks on every merge can be computationally expensive. The implementation must use efficient data structures.
2.  **Complexity:** The merge logic becomes significantly more complex to implement and test.
3.  **Schema Evolution:** How do you handle a merge where the branches have conflicting *changes to the schema itself*? This requires an even more sophisticated, multi-layered merge strategy.
# User

Performance: Loading a large schema and performing these checks on every merge can be computationally expensive. The implementation must use efficient data structures.Schema Evolution: How do you handle a merge where the branches have conflicting changes to the schema itself? This requires an even more sophisticated, multi-layered merge strategy.

## Key Concepts in this Implementation

1.  **Performance (`SchemaIndex` and `SchemaCache`):**
    *   `SchemaIndex`: A pre-computed, serializable struct that holds ready-to-use schema constraints (`FunctionalProperties`, `DisjointClasses`, etc.).
    *   `SchemaCache`: A thread-safe, in-memory LRU cache (simulated here with a map and mutex) that stores `SchemaIndex` objects to make schema lookups instantaneous during a merge.
    *   `BuildSchemaIndex`: The function that performs the expensive work of parsing schema quads, but it's only called when a schema changes.

2.  **Schema Evolution (`Merge` function):**
    *   **Phase 1 (Schema Reconciliation):** The `Merge` function first isolates and merges *only* the schema graphs. If there's a conflict (e.g., a predicate's range is changed differently in both branches), it stops and reports the schema conflict.
    *   **Phase 2 (Data Validation):** If the schema is reconciled successfully, the `Merge` function proceeds. It uses the *newly merged schema* as the single source of truth to validate all data changes from both branches. This correctly catches data that was valid under an old schema but is invalid under the new one.


