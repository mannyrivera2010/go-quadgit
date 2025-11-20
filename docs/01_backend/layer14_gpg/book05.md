# **Book 5: Interoperability and Trust**

By the end of this book, `go-quadgit` is no longer an island. It is a fully-fledged citizen of the web, capable of proving its own history, speaking the standard languages of its peers, and exposing its versioned knowledge in a way that is both powerful and profoundly intuitive.

**Layer III: The Application & Ecosystem**

**Sub-System 6: Ecosystem Integration & Trust**

**Subtitle:** *From a Closed System to an Open, Verifiable, and Connected Platform*

**Book Description:** A database, no matter how powerful, becomes a silo if it cannot speak the language of other systems or provide verifiable proof of its own integrity. This book is dedicated to transforming `go-quadgit` from a self-contained system into a trusted and interoperable citizen of the data ecosystem.

First, we will build the cryptographic foundation of trust by implementing GPG signing and verification, creating a non-repudiable chain of custody for every change. We will leverage this to build the ultimate provenance tool, the `blame` command. Next, we will construct a powerful SPARQL query engine, allowing the vast world of existing RDF tools to connect to and query our versioned data. Finally, we will embrace the principles of the web itself, implementing support for standard RDF serializations and designing a Linked Data Platform (LDP) compliant API that makes every piece of versioned knowledge an addressable, interactive, and discoverable resource.

**Prerequisites:** Readers should be familiar with the core `go-quadgit` architecture, including the multi-instance datastore and the REST server. A basic understanding of public-key cryptography (GPG), SPARQL, and RDF serialization formats will be beneficial.

## **Part I: Establishing Cryptographic Trust**

*This part focuses on building features that provide mathematical proof of data integrity and authorship.*

## **Chapter 1: GPG Signing for Commit Authenticity and Integrity**
*   **1.1: The Trust Deficit: Why the Author Field Isn't Enough**
    *   A detailed discussion of the threat model: How can a malicious actor forge a commit? Why is proving authorship critical for compliance and security?
*   **1.2: A Practical Introduction to GPG Signing**
    *   Explains the concepts of public/private keys, digital signatures, and the guarantees of Authenticity and Integrity they provide.
*   **1.3: Extending the `Commit` Data Model**
    *   Details the addition of the `Signature` field to the `Commit` struct and the importance of the `omitempty` tag.
*   **1.4: The Signing Workflow: Orchestrating the GPG Executable**
    *   A full implementation guide for the `sign` callback function. It shows how to use Go's `os/exec` to safely call the user's local `gpg` agent, pipe the canonical commit data to it, and capture the resulting ASCII-armored signature.

## **Chapter 2: Integrating Signature Verification into Core Workflows**
*   **2.1: The Verification Workflow**
    *   Implements the reverse process: a `verify` helper function that extracts the signature and the original signed data from a `Commit` object and passes them to `gpg --verify`.
*   **2.2: At-a-Glance Trust: `log --show-signature`**
    *   Shows how to modify the `log` command to call the `verify` function for each commit and display a "Good signature," "BAD signature," or "Untrusted key" status in the output.
*   **2.3: The Security Gate: Automatic Verification on `merge` and `push`**
    *   Details the modification of the `Store.Merge()` and the server-side `push` logic to automatically verify incoming commits. Explains the logic for aborting on a bad signature versus warning on an untrusted one.

## **Chapter 3: The `blame` Command: Providing Line-by-Line Data Provenance**
*   **3.1: The Goal: Answering "Where Did This Fact Come From?"**
    *   Defines the user story and the expected output of the `go-quadgit blame <graph-iri>` command.
*   **3.2: The Scalable Reverse History Walk Algorithm**
    *   A detailed, step-by-step implementation of the efficient `blame` algorithm. It explains how to use the scalable `StreamDiff` function to inspect the changes introduced by each historical commit without running out of memory.
*   **3.3: Building the `BlameViewer`: A UI-Friendly JSON Output**
    *   Implements a `--format=json` flag for the `blame` command, which produces a structured output ideal for consumption by a web frontend, turning the command-line tool into a powerful backend for a visual provenance explorer.
