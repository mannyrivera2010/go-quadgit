
## 


# **Book 13: Web-Native Integration**

**Layer III: The Application & Ecosystem**
**Sub-System 6: Ecosystem Integration & Trust**

**Subtitle:** *From a Database API to a Crawlable Web of Versioned Knowledge*

**Book Description:** In the previous book, we built a powerful SPARQL endpoint, making our data accessible to the semantic web community. However, a query endpoint still treats the database as a monolithic black box that you ask questions of. The vision of the web, and of Linked Data, is more profound: a network of *addressable resources* that link to one another.

This book is the final step in `go-quadgit`'s journey to becoming a true citizen of the web. We will first equip our system to "speak" the most common RDF serialization formats, making data import and export seamless. Then, we will embark on a major architectural enhancement: re-designing parts of our REST API to be compliant with the **Linked Data Platform (LDP)**. You will learn how to make every entity in your versioned knowledge graph an addressable, interactive, and discoverable resource with its own URL. This transforms our API from a simple query interface into a truly "web-native" platform.

**Prerequisites:** Readers must have a complete understanding of the `go-quadgit` REST API and versioning model. A strong grasp of RESTful principles, HTTP semantics, and the core ideas of Linked Data (IRIs, resources, RDF) is essential.



## **Part II: The Linked Data Platform (LDP)**

*This part is a deep dive into re-architecting our API around resource-centric, version-aware URLs.*

## **Chapter 4: Principles of the Linked Data Platform**
*   **4.1: Beyond Query Endpoints: The Resource-Centric Web**
    *   Explains the core LDP philosophy: instead of a single `/query` endpoint, every "thing" gets its own URL.
*   **4.2: HATEOAS: The Engine of Application State**
    *   Introduces the concept of "Hypermedia as the Engine of Application State." This means that API responses should not just contain data; they must contain *links* to related resources and possible actions, allowing a client to "crawl" the API and discover functionality without prior knowledge.
*   **4.3: Our Challenge: Merging LDP with Versioning**
    *   Formally defines the central problem: How to make a resource addressable *at a specific version*. This section re-introduces our chosen URL structure as the solution: `/ns/<ns>/refs/<type>/<ref>/resources/<iri>`.

## **Chapter 5: Implementing a Resource-Centric Read API**
*   **5.1: The Resource Endpoint: `GET .../resources/<iri>`**
    *   A full implementation of the HTTP handler for this endpoint. It will:
        1.  Parse the namespace, ref, and resource IRI from the URL.
        2.  Resolve the ref to a commit hash.
        3.  Use the internal query engine to `DESCRIBE` the resource at that commit.
        4.  Serialize the resulting graph into rich JSON-LD.
*   **5.2: Implementing HATEOAS**
    *   Shows how to modify the JSON-LD serializer. When it encounters an object that is an IRI, instead of just outputting the IRI string, it must transform it into a full, linked object with an `@id` pointing to its own LDP URL. This is what makes the graph crawlable.
*   **5.3: Implementing LDP Containers for Named Graphs**
    *   Implements the `/ns/.../graphs/<graph_iri>` endpoint. A `GET` request here will act as an LDP Container, returning a list of links to all the resources contained within that named graph.

## **Chapter 6: `PUT`/`PATCH`/`DELETE` as Versioned Commits**
*   **6.1: The Ultimate Integration: RESTful Writes as Atomic Commits**
    *   This is the capstone chapter of the book. It implements the handlers for the LDP write methods.
*   **6.2: `PATCH .../resources/<iri>` (Partial Update)**
    *   **Implementation:** The handler accepts a standard format like JSON Merge Patch or SPARQL Update. It performs the crucial `If-Match` optimistic locking check. If successful, it constructs a changeset (quads to delete, quads to add) and calls the core `store.Commit()` method with an auto-generated commit message.
*   **6.3: `PUT .../resources/<iri>` (Full Replacement)**
    *   **Implementation:** The handler takes a full RDF document as the request body. After the `If-Match` check, it calculates the diff between the existing state of the resource and the new state, and uses that to create a commit.
*   **6.4: `POST .../graphs/<graph_iri>` (Creating New Resources)**
    *   **Implementation:** The handler for `POST`ing to a container. It accepts an RDF document, mints a new IRI for the resource, creates a commit that adds the resource's triples to the named graph, and returns `HTTP 201 Created` with a `Location` header pointing to the URL of the newly created resource.

By the end of this book, `go-quadgit` will have achieved its final form. It will not only be a powerful, version-controlled database but also a true web platform. Its data will be accessible not just through specialized queries but as a browseable, interactive, and discoverable network of resources, fully integrated with the standards and principles that underpin the World Wide Web itself.



### **20.4 (Detailed Expansion): Linked Data Platform (LDP) - Making Resources Addressable**

The Linked Data Platform (LDP) is a W3C recommendation that provides a standard way to build REST APIs for RDF data. Its philosophy is simple but powerful: treat every "thing" in your knowledge graph not as a row in a database to be queried, but as a **web resource** with its own stable URL that can be directly interacted with using standard HTTP methods.

While our existing REST API is functional, it is "query-centric" (`GET /diff?from=...`, `POST /merges`). Adopting LDP principles allows us to create a more intuitive, "resource-centric" API that is crawlable, discoverable, and aligns with the fundamental architecture of the World Wide Web.

#### **The Crucial Insight: Integrating Versioning with LDP**

LDP itself does not specify versioning. Our primary challenge is to merge LDP's resource-centric model with `quad-db`'s core versioning model. A request like `GET /resources/ex:product1` is ambiguous. Does it mean the version on `main`? The version on a feature branch? The version from last year?

The solution is to make the version **an explicit part of the resource URL**. Our LDP-compliant URL structure will be:

`/api/v1/ns/<namespace>/refs/<type>/<ref_name>/resources/<resource_iri>`

*   **`<namespace>`**: The `quad-db` namespace (e.g., `production`).
*   **`<type>`**: `heads` for branches or `tags` for tags.
*   **`<ref_name>`**: The URL-encoded name of the branch or tag (e.g., `main`, `v1.0.0`).
*   **`<resource_iri>`**: The URL-encoded IRI of the subject we want to interact with.

**Example URL:**
`/api/v1/ns/production/refs/heads/main/resources/http:%2F%2Fexample.org%2Fproducts%2Fprod123`

This URL unambiguously identifies the resource `ex:prod123` as it exists at the `HEAD` of the `main` branch in the `production` namespace.

#### **Reading Resources: `GET` and HATEOAS**

A `GET` request to a resource URL fetches its description.

*   **Request:** `GET /api/v1/ns/production/refs/heads/main/resources/ex:prod123`
*   **Backend Implementation:**
    1.  The handler parses the URL to get the namespace, ref, and resource IRI.
    2.  It resolves the ref `main` to a commit hash.
    3.  It uses the query engine from Chapter 17 to execute a query against that specific version of the graph:
        ```sparql
        DESCRIBE <http://example.org/products/prod123>
        # Or, more simply:
        # SELECT ?p ?o WHERE { <http://example.org/products/prod123> ?p ?o }
        ```
    4.  The results are serialized into a rich format like JSON-LD.
*   **Response (JSON-LD with HATEOAS):**
    The response doesn't just contain data; it contains *links*. This principle is called **Hypermedia as the Engine of Application State (HATEOAS)**.

    ```json
    // HTTP/1.1 200 OK
    // ETag: "a1b2c3d..."  (The commit hash of main's HEAD)
    // Content-Type: application/ld+json

    {
      "@context": {
        "ex": "http://example.org/ns#",
        "rdf": "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
      },
      "@id": "http://example.org/products/prod123",
      "rdf:type": { "@id": "ex:Product" },
      "ex:hasPrice": 129.99,
      "ex:manufacturer": {
        // This is HATEOAS: the object is not just a string, but a link to another resource.
        "@id": "/api/v1/ns/production/refs/heads/main/resources/ex:companyA"
      }
    }
    ```
    A client can now "crawl" the graph by following the `ex:manufacturer` link to discover information about the manufacturer.

#### **Modifying Resources: `PUT`, `PATCH`, `DELETE` as Commits**

This is the most powerful part of the integration. Every LDP write operation becomes a **new `quad-db` commit**, giving us a full versioned history of every change made via the LDP interface. This requires optimistic locking.

*   **`PATCH /.../resources/ex:product123` (Partial Update)**
    *   **Use Case:** The user wants to change only the price of the product.
    *   **Request:** The `PATCH` body would contain a description of the change, for example, using the SPARQL `UPDATE` format or a JSON Patch format. The request **must** include the `If-Match` header with the ETag of the resource's version.
    *   **Backend Implementation:**
        1.  Perform the optimistic locking check (compare `If-Match` with the branch's current `HEAD` hash).
        2.  Read the *current* set of triples for `ex:product123`.
        3.  Apply the patch: delete the old `ex:hasPrice` triple and add the new one.
        4.  Call `store.Commit()` with the new set of triples for the relevant named graph. The commit message could be automatically generated: `PATCH on resource <ex:product123>`.
        5.  Return `200 OK` with the new ETag (the new commit hash).

*   **`PUT /.../resources/ex:product123` (Full Replacement)**
    *   **Use Case:** The user wants to completely replace the description of the product with a new set of triples.
    *   **Request:** The `PUT` body contains the *full, new* RDF description. It also requires the `If-Match` header.
    *   **Backend Implementation:**
        1.  Perform the optimistic locking check.
        2.  Instead of patching, the implementation deletes *all* existing triples for `ex:product123` and replaces them with the triples from the request body.
        3.  Call `store.Commit()` to save this new state.
        4.  Return `201 Created` or `200 OK`.

*   **`DELETE /.../resources/ex:product123`**
    *   **Use Case:** Delete the product and all its defining triples.
    *   **Backend Implementation:**
        1.  Perform the optimistic locking check with `If-Match`.
        2.  Find all triples where `ex:product123` is the subject.
        3.  Call `store.Commit()`, providing a change set that deletes all these triples.
        4.  Return `204 No Content`.

#### **LDP Containers and Named Graphs**

LDP defines "Containers" (`ldp:Container`) which are resources that hold other resources. This maps perfectly to our **Named Graphs**.

*   **Endpoint:** `/api/v1/ns/production/refs/heads/main/graphs/ex:productCatalog`
*   **`GET` on a Container:** A `GET` request to this URL would list all the resources contained within that named graph. The response would be an RDF document listing the IRIs of the subjects.
    ```json
    {
      "@id": "...",
      "ldp:contains": [
        { "@id": ".../resources/ex:product123" },
        { "@id": ".../resources/ex:product124" }
      ]
    }
    ```
*   **`POST` to a Container:** This is how you create a new resource. A `POST` request to the container URL with the RDF description of a new product in the body would:
    1.  Have the server mint a new, unique IRI for the product.
    2.  Create a new commit that adds the new product's triples to the `ex:productCatalog` named graph.
    3.  Return `201 Created` with a `Location` header pointing to the URL of the newly created resource (e.g., `Location: .../resources/ex:product125`).

By adopting these LDP patterns, `quad-db` exposes its versioned data in a way that is profoundly "webby." It moves beyond being a database with a custom API and becomes a true platform for linked data, where every piece of information, at any point in its history, is an addressable, interactive, and discoverable resource on the web.
