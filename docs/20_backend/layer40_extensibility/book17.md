
## 
**Layer IV: Operations & Reliability**
**Sub-System 8: Administration & Extensibility**

# **Book 17: The Platform Extensibility Guide**

**Subtitle:** *Customizing `go-quadgit` with Hooks, Plugins, and Application Integrations*

**Book Description:** No platform can anticipate every user's need. A truly powerful system is not a closed box, but an open framework that invites extension and integration. This book is the definitive guide to customizing and extending the `go-quadgit` platform. It is designed for the advanced user and systems integrator who needs to enforce custom business rules, connect `go-quadgit` to external tools, and build automated data governance workflows.

We will start by designing and implementing a robust hook system, inspired by Git Hooks, that allows you to execute custom scripts at critical points in the application's lifecycle. You will then walk through three detailed, practical use cases: building a `pre-commit` hook to enforce data quality with SHACL, a `post-commit` hook to trigger CI/CD pipelines and send notifications, and a `pre-merge` hook to integrate with external approval systems like Jira. Finally, we will revisit the high-level vision of building entire applications on top of the `go-quadgit` platform, providing architectural blueprints for creating sophisticated, data-aware software.

**Prerequisites:** This book is for advanced users. A strong understanding of the `go-quadgit` CLI, its core concepts (commit, merge), and general scripting (e.g., shell scripts, Python, or Go) is essential. Familiarity with CI/CD and external tools like SHACL validators is beneficial.



## **Part I: The Hook System**

*This part covers the design and implementation of the core extensibility mechanism.*

## **Chapter 1: Designing the Hook System**
*   **1.1: The Philosophy of Hooks: Safe, Externalized Logic**
    *   Explains why a hook system is a secure and maintainable way to add custom logic. It runs user code in a separate process, preventing a faulty script from crashing the main `go-quadgit` server. It also allows users to write hooks in any language they choose.
*   **1.2: The Hook Execution Model**
    *   Details the mechanics: `go-quadgit` looks for an executable file with a specific name (e.g., `pre-commit`) in the `.quadgit/hooks/` directory of a repository.
    *   **"Pre" Hooks (Veto Power):** Explains that `pre-` hooks are synchronous. The core command waits for the script to finish. If the script exits with a non-zero status code, the operation is aborted.
    *   **"Post" Hooks (Reactionary):** Explains that `post-` hooks are asynchronous ("fire and forget"). The core command does not wait for them to complete and their exit code is ignored. This ensures they don't slow down the user's workflow.
*   **1.3: The Hook Points: A Comprehensive List**
    *   A reference table detailing every available hook, when it's triggered, what data is passed to its `stdin`, and whether it has veto power. This includes `pre-commit`, `post-commit`, `pre-merge`, `pre-push`, `post-receive` (server-side), etc.
*   **1.4: Implementing the `HookRunner` in Go**
    *   Provides the Go code for the internal `HookRunner` utility. It shows how to use `os/exec` to securely run an external script, set a timeout, capture `stdout`, `stderr`, and the `exit code`, and pass data via `stdin`.

## **Chapter 2: Implementing Hooks in Core Workflows**
*   **2.1: Integrating `pre-commit` and `post-commit`**
    *   A line-by-line refactoring of the `Store.Commit()` method. It shows exactly where the `HookRunner` is called for the `pre-commit` hook (before the transaction) and how the `post-commit` hook is launched in a new goroutine after the transaction succeeds.
*   **2.2: Integrating `pre-merge`**
    *   Shows how the `Store.Merge()` method is modified to run the `pre-merge` hook after the merge base is found but before the final merge commit is created.
*   **2.3: Implementing Server-Side Hooks: `post-receive`**
    *   Details the implementation of the most important server-side hook. The `ReceivePack` handler, after successfully updating a reference from a `push`, will trigger the `post-receive` hook, passing it a list of the references that were changed, their old hashes, and their new hashes.



## **Part II: Practical Use Cases and Tutorials**

*This part provides three complete, end-to-end tutorials for building common and powerful integrations.*

## **Chapter 3: Use Case 1: Pre-Commit Validation with SHACL**
*   **3.1: The Goal: Enforcing Data Quality at the Source**
    *   Defines the user story: an organization wants to guarantee that no quad that violates their corporate SHACL shapes can ever enter the repository.
*   **3.2: Setting up the Environment**
    *   Includes instructions for versioning the `shapes.ttl` file within the `go-quadgit` repository itself.
    *   Provides a `Dockerfile` for a SHACL validation tool (like `pySHACL`) to ensure a consistent validation environment.
*   **3.3: Writing the `pre-commit` Script**
    *   A detailed, commented shell script (`pre-commit.sh`).
    *   **Step 1:** Read the staged diff from `stdin`.
    *   **Step 2:** Use `go-quadgit` commands to get the current state of the graph and apply the diff in-memory to create a "proposed state" file.
    *   **Step 3:** Execute the SHACL validator against the proposed state.
    *   **Step 4:** Check the validation report. If it's not conformant, print the report to `stderr` and `exit 1`. Otherwise, `exit 0`.
*   **3.4: The User Experience**
    *   Shows the exact console output a user would see when their commit is rejected by the hook, demonstrating the immediate, actionable feedback.

## **Chapter 4: Use Case 2: Post-Commit Notifications and CI/CD Integration**
*   **4.1: The Goal: Automating Downstream Processes**
    *   Defines the user story: on every push to the `main` branch, the team wants a notification in their Slack channel and a new build to be triggered in GitLab CI.
*   **4.2: Writing the Server-Side `post-receive` Script**
    *   Provides a script (e.g., in Python or Go) that is placed in the `.quadgit/hooks/` directory on the **server**.
    *   The script reads the old hash, new hash, and ref name from its arguments.
    *   It includes a conditional check: `if [ "$ref_name" = "refs/heads/main" ]`.
    *   Inside the `if` block, it uses `curl` to make two API calls:
        1.  A `POST` to the Slack Incoming Webhook URL with a formatted JSON message.
        2.  A `POST` to the GitLab "trigger pipeline" API endpoint.
*   **4.3: The Result: A Fully Automated DataOps Pipeline**
    *   Shows a screenshot of the resulting Slack notification and the new pipeline running in GitLab, demonstrating the end-to-end automation.

## **Chapter 5: Use Case 3: Integrating with External Approval Systems**
*   **5.1: The Goal: Enforcing Enterprise Governance**
    *   Defines the user story: a merge to a `production` namespace is only allowed if the corresponding Jira ticket is in the "Approved for Release" state.
*   **5.2: Writing the Server-Side `pre-merge` Hook**
    *   This hook is implemented as part of the `POST /mergerequests/:id/merge` API handler logic on the server.
    *   Before calling `store.Merge()`, the handler executes the hook.
    *   The hook script parses the source branch name to extract a ticket ID (e.g., from `feature/PROJ-456-fix`).
    *   It uses an API token (stored as a secure environment variable on the server) to make a `GET` request to the Jira API (`/rest/api/2/issue/PROJ-456`).
    *   It parses the JSON response and checks the value of the `status` field.
    *   If the status is not "Approved for Release," it exits with a non-zero code, causing the API request to fail with a `403 Forbidden` and a clear error message.
*   **5.3: Closing the Loop: The `post-merge` Hook**
    *   Shows how a corresponding `post-merge` hook can automatically add a comment to the Jira ticket and transition its status to "Deployed," creating a bi-directional integration.


