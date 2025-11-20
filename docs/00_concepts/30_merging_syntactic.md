# Merging

## `quad-db merge <branch-name>`

1.  **Find Common Ancestor:** The CLI walks the history of the current branch (`HEAD`) and the `<branch-name>` to find the most recent shared commit hash.

2.  **Calculate Diffs:** It generates two lists of changes (quad additions/deletions) for each branch since the common ancestor.

3.  **Detect Conflicts:** The CLI iterates through the changes.
    *   It uses a map where keys are quad identifiers. It checks if the same quad was added in one diff and deleted in another.
    *   It uses a second map where keys are `Subject:Predicate:Graph`. It checks if both branches add a quad with this key but with a *different object*.
    *   If any such conditions are met, a conflict is flagged.

4.  **Handle the Outcome:**
    *   **No Conflicts:** The merge is performed automatically. A new commit is created with two parents. The command reports a successful merge.
    *   **Conflicts Found:** The merge is halted. The CLI does *not* create a commit. Instead, it does the following:
        a.  **Enters a "Merging" State:** It writes a special file to the database directory (e.g., `.quad-db/MERGE_HEAD`) containing the hash of the source commit being merged. This tells other commands that a merge is in progress.
        b.  **Generates a Conflict Report:** It creates a human-readable file, for instance, `.quad-db/MERGE_MSG`, that details every conflict. This is our equivalent of Git's `<<<<<<<`, `=======`, `>>>>>>>` markers.

## The Conflict File (`MERGE_MSG`)

The conflict report would be structured to be clear and actionable.

```
# CONFLICT: Contradictory changes for the same quad.
# The quad '<person:Bob> <foaf:knows> <person:Charlie> .' was kept in 'main' but deleted in 'feature'.
# To resolve, add one of the following lines to your staged changes:
#
# ADD <person:Bob> <foaf:knows> <person:Charlie> .
# DEL <person:Bob> <foaf:knows> <person:Charlie> .

# ---

# CONFLICT: Different objects for the same subject-predicate pair.
# Subject: <person:Alice>, Predicate: <ex:hasAge>
#
# Value from 'main':
# ADD <person:Alice> <ex:hasAge> "30"^^xsd:integer .
#
# Value from 'feature':
# ADD <person:Alice> <ex:hasAge> "31"^^xsd:integer .
#
# To resolve, stage the correct quad(s). You might choose one, both, or neither.
```

## The User's Resolution Workflow

1.  User runs `quad-db merge feature`. The command fails and prints:
    ```
    Automatic merge failed; fix conflicts and then commit the result.
    Conflicts reported in .quad-db/MERGE_MSG
    ```

2.  The user inspects the `MERGE_MSG` file to understand the conflicts.

3.  The user creates a new file, `resolution.nq`, to fix the conflicts. For the examples above, they might add the following lines to `resolution.nq`:
    ```n-quads
    # Decided that Bob no longer knows Charlie.
    DEL <person:Bob> <foaf:knows> <person:Charlie> .

    # Decided that Alice's real age is 31.
    ADD <person:Alice> <ex:hasAge> "31"^^xsd:integer .
    ```

4.  The user stages this resolution using the `add` or a new `resolve` command:
    ```bash
    quad-db add resolution.nq
    ```
    This command would need to be enhanced to handle `DEL` syntax.

5.  The user finalizes the merge by committing:
    ```bash
    quad-db commit -m "Merge branch 'feature'; resolved age and acquaintance conflicts"
    ```
    When `commit` is run, it detects the `.quad-db/MERGE_HEAD` file. Instead of creating a regular commit, it creates a **merge commit** with two parents (the target `HEAD` and the hash stored in `MERGE_HEAD`). It then deletes the `MERGE_HEAD` and `MERGE_MSG` files, concluding the merging state.

## Advanced Consideration: Schema-Aware Merging

For an even more powerful system, the merge logic could be schema-aware. If the database also stored RDFS or OWL ontology information, it could:
*   **Automatically flag conflicts** for any predicate defined as an `owl:FunctionalProperty`.
*   **Avoid flagging conflicts** for predicates that are explicitly non-functional.
*   Understand class disjointness and other constraints to detect more subtle logical inconsistencies introduced by a merge.

## The Merge Process: A Three-Way Merge

The process mirrors Git's three-way merge. When you command `quad-db merge feature`, the system needs to reconcile three versions of the graph:

1.  **The Target Commit (`main` branch HEAD):** The state of the graph you are merging *into*.
2.  **The Source Commit (`feature` branch HEAD):** The state of the graph you are merging *from*.
3.  **The Common Ancestor (or Merge Base):** The most recent commit that both `main` and `feature` share in their history.

The core logic is to calculate two diffs:
*   **Diff A:** The changes between the `Common Ancestor` and the `Target Commit`.
*   **Diff B:** The changes between the `Common Ancestor` and the `Source Commit`.

The system then attempts to apply `Diff B` onto the `Target Commit`'s state.

## Automatic Merging (The Happy Path)

An automatic merge is possible when the changes in the two branches are non-overlapping.

*   **Example:**
    *   **Common Ancestor:** Contains quad `Q1`.
    *   **`main` branch:** Adds a new quad, `Q2`. (`Diff A` = `ADD Q2`)
    *   **`feature` branch:** Adds a different new quad, `Q3`. (`Diff B` = `ADD Q3`)

The result is trivial: the final merged state will contain `Q1`, `Q2`, and `Q3`. The system can create a new **merge commit** on `main` that has two parents (the old `main` HEAD and the `feature` HEAD) and points to a new tree representing this combined state.

## What Constitutes a Conflict?

A conflict arises when `Diff A` and `Diff B` make contradictory changes to the same part of the graph. In a quad store, this is more nuanced than conflicting lines in a text file.

## 1. Direct Quad Conflict

This is the most straightforward conflict. The same quad (identical Subject, Predicate, Object, and Graph) is handled differently in each branch.

*   **Scenario:**
    *   **`main`:** Deletes the quad `<S> <P> <O> <G>`.
    *   **`feature`:** Modifies the graph, which might be perceived as keeping the quad alive or modifying it (e.g., if we supported quad modification directly). More simply, the quad was *not* deleted.
    *   **Conflict:** One history says the quad should be gone, the other says it should remain. The system cannot decide.

## 2. Semantic Conflict (Functional Predicates)

This is a more complex and common scenario in knowledge graphs. It occurs when two branches provide different objects for the same subject-predicate pair, especially if that predicate is conceptually "functional" (i.e., should only have one value).

*   **Scenario:** The predicate `<ex:hasAge>` should only have one value for any given subject.
    *   **Common Ancestor:** Does not contain an age for `<person:Alice>`.
    *   **`main`:** Adds the quad `<person:Alice> <ex:hasAge> "30"^^xsd:integer .`
    *   **`feature`:** Adds the quad `<person:Alice> <ex:hasAge> "31"^^xsd:integer .`

**This is a critical conflict.** The system cannot know if Alice is 30 or 31. Even if the predicate is not strictly functional (e.g., `<foaf:knows>`), adding two different values for the same subject-predicate pair from different branches often represents a semantic disagreement that a human needs to resolve.

# Conflict Resolution Workflow

Here’s how the `quad-db merge` command would handle this.

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

## Production-Ready Solution for Performance

The bottleneck is the repeated loading and parsing of schema quads on every merge. The solution is a two-layer caching and pre-computation strategy to ensure schema validation is nearly instantaneous for the vast majority of merge operations.

## Layer 1: Pre-computed Schema Index in BadgerDB (Commit-Time Optimization)

Instead of processing the raw schema quads at merge time, we pre-process them whenever the schema itself changes and store the result in a dedicated index within BadgerDB.

**Implementation Strategy:**

1.  **Identify Schema Changes:** During the `commit` operation, the system checks if the commit modifies the dedicated schema graph (e.g., `<urn:quad-db:schema>`). This is done by comparing the tree hash for that graph between the new commit and its parent.

2.  **Trigger Indexing Job:** If the schema has changed, a background or inline job is triggered to build a "Schema Index." The new commit will not be finalized until the index is built.

3.  **Structure the Schema Index:** This index is a set of key-value pairs in BadgerDB, keyed by the SHA-1 hash of the schema's content blob. This makes the index content-addressable and automatically versioned.

    *   **Key:** `index:schema:<schema-blob-hash>:functional`
    *   **Value:** A marshalled (e.g., JSON or Protobuf) list of all property IRIs that are `owl:FunctionalProperty`.
        `["http://example.org/hasSSN", "http://example.org/hasPrimaryEmail"]`

    *   **Key:** `index:schema:<schema-blob-hash>:maxCardinality:<property-iri>`
    *   **Value:** The integer cardinality limit.
        `2`

    *   **Key:** `index:schema:<schema-blob-hash>:disjoint:<class-iri>`
    *   **Value:** A marshalled list of classes that are disjoint with the key's class.
        `["http://example.org/Adult", "http://example.org/Corporation"]`

4.  **Benefit (Write-Time Cost for Read-Time Gain):** The computational cost is moved from every merge operation to the rare event of a schema commit. Merges, which are far more frequent, no longer need to parse RDF. They simply read a few pre-computed keys from BadgerDB.

## Layer 2: In-Memory LRU Cache (Merge-Time Optimization)

Even reading from BadgerDB has overhead. To make schema validation virtually free, we add an in-memory cache for the fully deserialized Schema Index objects.

**Implementation Strategy:**

1.  **Instantiate an LRU Cache:** When the `quad-db` application starts, it initializes a global, thread-safe LRU (Least Recently Used) cache. The size can be configured (e.g., cache the last 50 most used schemas).

2.  **Modify the Merge Workflow:**
    a. When `quad-db merge` begins, it identifies the required `schema-blob-hash` from the target branch's tree.
    b. **Check In-Memory Cache:** It first queries the LRU cache with the `schema-blob-hash`.
    c. **Cache Hit:** If found, the fully parsed, ready-to-use schema constraint objects are returned instantly. This is the optimal path.
    d. **Cache Miss:** If not in the LRU cache:
        i.  The system reads the pre-computed Schema Index from BadgerDB (using the keys from Layer 1).
        ii. It deserializes these values into ready-to-use Go objects (maps, slices, etc.).
        iii. It stores these objects in the LRU cache with the `schema-blob-hash` as the key.
        iv. It then proceeds with the merge validation.

**Production-Ready Outcome:**
The first time a specific schema version is used in a merge, there's a small, one-time cost to read the index from BadgerDB and populate the in-memory cache. Every subsequent merge operation anywhere in the repository that uses that *same schema version* will have its schema validation step completed in nanoseconds via a direct memory read. This architecture is extremely fast, scalable, and robust.

---

## Production-Ready Solution for Schema Evolution

Handling conflicting schema changes requires formalizing the merge process into a strict, user-guided, two-phase operation. The system must refuse to merge data until all ambiguity about the governing schema is resolved.

This is the **"Schema Reconciliation First"** strategy.

## Phase 1: Reconcile the Schema

1.  **Isolate Schema Changes:** The `quad-db merge <branch>` command first ignores all data graphs. It performs a three-way merge exclusively on the quads within the schema graph (`<urn:quad-db:schema>`).

2.  **Detect Direct Schema Conflicts:** It checks for conflicts at the axiom level.
    *   **Direct Contradiction:** `main` sets `<P> rdfs:range xsd:integer`, while `feature` sets `<P> rdfs:range xsd:string`. This is a hard conflict.
    *   **Additive Contradiction:** `main` makes `<P>` an `owl:FunctionalProperty`. `feature` makes it an `owl:SymmetricProperty`. While not a direct contradiction on a single quad, these are fundamental changes to the same entity (`<P>`) that require human intervention.

3.  **Halt on Conflict and Report:** If a schema conflict is found, the merge process halts immediately. It will *not* proceed to data merging. The CLI provides a specific report:
    ```
    CONFLICT (Schema): Branches contain conflicting changes for the predicate '<ex:hasAge>'.
    'main' sets the range to 'xsd:integer'.
    'feature' sets the range to 'xsd:string'.
    Please resolve this schema conflict first.
    ```

4.  **User Resolution:** The user must manually create a `schema-resolution.nq` file that defines the desired final state of the schema and `quad-db add` it.

5.  **Create an Intermediate Schema Commit:** The user runs `quad-db commit`. The system detects it's in a merge state and creates a commit that contains **only the merged schema**. This commit has two parents and a clear message, e.g., "Merge commit (schema only) for branch 'feature'". The repository is now in a state where the schema is consistent, but the data merge is still pending.

## Phase 2: Validate Data Against the Reconciled Schema

1.  **Automatic Resumption:** After the intermediate schema commit is created, the system can automatically (or via a `quad-db merge --continue` command) resume the merge process.

2.  **Re-evaluate Data Diffs:** Now, the system uses the schema from the **newly created intermediate commit** as the single source of truth.

3.  **Detect Data-vs-Schema Conflicts:** It re-validates the data changes from both branches against this new, unified schema. This is where it will catch the classic evolution problem:
    *   **Scenario:** `main` made `<P>` functional (Phase 1). `feature`, based on the old schema, added two values for `<P>`.
    *   **Outcome:** During Phase 2, the system now sees that adding two values for `<P>` violates the *newly merged* schema. It flags this as a data conflict.

4.  **Standard Data Conflict Resolution:** The user is now presented with the familiar data-level conflicts ("Functional property violation for `<P>`"). They resolve these conflicts by staging the correct quads and running `quad-db commit` one last time.

5.  **Final Merge Commit:** This final commit "amends" or replaces the intermediate schema commit, resulting in a single, clean merge commit that contains the fully reconciled schema and data.

**Production-Ready Outcome:**
This two-phase process transforms a chaotic, unpredictable problem into a structured, deterministic workflow.
*   **It enforces correctness:** No inconsistent data can be committed because the schema is agreed upon first.
*   **It's understandable:** It separates concerns for the user. They first focus only on fixing the ontology, then they focus on fixing the instance data based on that fixed ontology.
*   **It's auditable:** The history clearly shows how schema and data conflicts were resolved, which is crucial for regulated environments.



