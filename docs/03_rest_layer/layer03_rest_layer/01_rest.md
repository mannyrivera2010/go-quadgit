
#### **Chapter 2: The REST API Server**
*   **2.1: Choosing a Framework: `net/http` vs. Gin**
    *   A discussion of the pros and cons, ultimately choosing a lightweight but powerful framework like Gin for its routing, middleware, and JSON binding capabilities.
*   **2.2: Designing the RESTful API Endpoints**
    *   A full breakdown of the API design, organized by resource (`/refs`, `/commits`, `/mergerequests`, etc.). We'll define the HTTP verb, URL path, request body, and success/error responses for each core operation.
*   **2.3: Middleware: The Power of Composition**
    *   Introduces the concept of middleware for handling cross-cutting concerns. We'll design middleware for logging, request ID generation, authentication, and authorization.
*   **2.4: Dependency Injection in Handlers**
    *   Shows the pattern for injecting the shared `quadstore.Store` instance and the configuration object into the request context so that all HTTP handlers have access to them.



## **Part I: Universal Data Exchange**

*This part focuses on ensuring data can flow into and out of `go-quadgit` using standard, widely-adopted formats.*

## **Chapter 1: The RDF Serialization Landscape**
*   **1.1: Why Format Matters: Interoperability and Usability**
    *   Compares and contrasts the most common RDF formats, explaining their strengths and weaknesses:
        *   **N-Quads/N-Triples:** Best for streaming and bulk loads.
        *   **Turtle/TriG:** Best for human readability and authoring.
        *   **JSON-LD:** Best for web developers and API integration.
        *   **RDF/XML:** The original, verbose format, necessary for legacy compatibility.
*   **1.2: Choosing a Go RDF Library**
    *   A review of available open-source Go libraries for handling RDF (e.g., `knakk/rdf`). We'll select a library and explain how to integrate it as a core dependency for all serialization/deserialization tasks.

## **Chapter 2: Implementing Multi-Format Parsers**
*   **2.1: Enhancing the `add` Command**
    *   Refactors the `go-quadgit add` command. Instead of only accepting `.nq` files, it will now inspect the file extension (`.ttl`, `.trig`, `.jsonld`) or a `--format` flag to select the correct parser from our chosen library.
*   **2.2: The `POST /importer` API Endpoint**
    *   Builds a new, dedicated REST endpoint for data ingestion. This endpoint will use HTTP content negotiation (the `Content-Type` header) to automatically detect the format of the uploaded RDF data.
    *   This handler will stream the uploaded file directly to the appropriate parser, making it memory-efficient for large uploads.

## **Chapter 3: Implementing Multi-Format Serializers**
*   **3.1: The "Export" Feature**
    *   Builds a new `go-quadgit export <ref> --format <format>` command that can dump the entire state of a commit into a specified RDF format.
*   **3.2: Enhancing the REST API with Content Negotiation**
    *   Refactors all existing `GET` endpoints that return RDF data. They will now inspect the client's `Accept` header.
    *   For example, a client can now request `Accept: application/ld+json` on a `/commits/:hash/data` endpoint and receive the graph data as JSON-LD instead of a custom format. This makes the API instantly more useful to a wider range of web clients.

