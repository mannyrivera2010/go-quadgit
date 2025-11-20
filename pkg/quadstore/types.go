// Package quadstore defines the public, embeddable API for interacting with a
// go-quadgit repository. It provides a stable interface for all core
// versioning and query operations.
package quadstore

import (
	"time"
)

// Quad represents a single, atomic RDF statement within a named graph.
// It is the fundamental unit of data in the system.
type Quad struct {
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object"`
	Graph     string `json:"graph"`
}

// Author contains metadata about the person who created a commit.
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CommitStats holds pre-computed statistics for a commit, allowing for
// fast retrieval of historical metrics without expensive on-the-fly calculations.
type CommitStats struct {
	TotalQuads int64 `json:"total_quads"`
	Added      int   `json:"added"`
	Deleted    int   `json:"deleted"`
}

// Commit represents a single, versioned point in the repository's history.
// It is an immutable, content-addressable object.
type Commit struct {
	Hash      string      `json:"hash"`
	Tree      string      `json:"tree"` // SHA-1 hash of the Tree object
	Parents   []string    `json:"parents"`
	Author    Author      `json:"author"`
	Message   string      `json:"message"`
	Timestamp time.Time   `json:"timestamp"`

	// Signature holds the detached, ASCII-armored PGP signature of the
	// marshalled commit data (excluding this field itself). It is empty
	// for unsigned commits.
	Signature string `json:"signature,omitempty"`

	// Stats contains pre-computed metrics about the state of the graph
	// at the time of this commit.
	Stats CommitStats `json:"stats"`
}

// Reference is a named, mutable pointer to a commit. It represents a branch or a tag.
type Reference struct {
	Name string `json:"name"` // The full reference name (e.g., "refs/heads/main" or "refs/tags/v1.0")
	Hash string `json:"hash"` // The commit hash this reference points to.
}

// ChangeType defines whether a change is an addition or deletion.
type ChangeType bool

const (
	Addition ChangeType = true
	Deletion ChangeType = false
)

// Change represents a single quad addition or deletion in a diff operation.
// This is used for streaming diff results.
type Change struct {
	Quad     Quad       `json:"quad"`
	Type     ChangeType `json:"type"`
}

// BlameResult associates a single quad with the commit that last introduced it.
// This is used for streaming the results of a blame operation.
type BlameResult struct {
	Quad   Quad    `json:"quad"`
	Commit *Commit `json:"commit"`
}

// Conflict represents a single point of contention found during a merge that
// prevents the merge from being completed automatically.
type Conflict struct {
	Type        string   `json:"type"`        // e.g., "SEMANTIC_CONFLICT_FUNCTIONAL_PROPERTY"
	Description string   `json:"description"` // A human-readable explanation of the conflict.
	Conflicting []string `json:"conflicting_quads"` // The string representations of the conflicting quads.
}

// BackupManifest contains metadata about a completed backup, required for
// performing subsequent incremental backups.
type BackupManifest struct {
	Timestamp       time.Time `json:"timestamp"`
	DatabaseVersion uint64    `json:"database_version"` // The BadgerDB version at the time of backup.
	IsIncremental   bool      `json:"is_incremental"`
}