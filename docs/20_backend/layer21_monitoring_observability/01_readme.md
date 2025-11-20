
## **Part III: Monitoring and Observability**

*This part focuses on giving administrators the visibility they need to understand the health and performance of the system.*

## **Chapter 5: Building the Insights Dashboard**
*   **5.1: The Goal: A Single Pane of Glass**
    *   Lays out the design for a comprehensive monitoring dashboard (e.g., in Grafana).
*   **5.2: Exposing Metrics for Prometheus**
    *   Implements a `/metrics` HTTP endpoint in the `go-quadgit-server` using the Go Prometheus client library.
*   **5.3: Key Metrics to Watch: A Checklist**
    *   **Application Metrics:** `http_requests_total`, `http_request_duration_seconds` (by endpoint), `active_websocket_connections`.
    *   **Go Runtime Metrics:** `go_goroutines`, `go_memstats_heap_alloc_bytes` (for detecting memory leaks).
    *   **BadgerDB Metrics:** Cache hit ratios (`badger_cache_hits / badger_cache_misses`), LSM-tree size, and Value Log size for each database instance. This is critical for validating the tuning of the multi-instance architecture.
    *   **Business Metrics:** `quadgit_commits_total`, `quadgit_merges_total`, `quadgit_namespaces_total`.
*   **5.4: Using the JSON Stats API for Deeper Insights**
    *   Shows how to set up a separate job that periodically runs `go-quadgit stats history --format=json` and ingests this data into a time-series database to track long-term trends in contributor activity and data churn.
