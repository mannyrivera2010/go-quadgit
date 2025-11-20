// Package quadstore defines the public, embeddable API for interacting with a
// go-quadgit repository. It provides a stable interface for all core
// versioning and query operations.
package quadstore

import (
	"context"
	"io"
)

// OpenOptions provides configuration for opening a repository.
type OpenOptions struct {
	// Path to the root directory where all database instances are stored.
	Path string
	// The namespace to operate on. If empty, uses a default namespace.
	Namespace string
}

// Store defines the public API for interacting with a versioned quad store repository.
// All implementations of this interface must be safe for concurrent use from multiple goroutines.
type Store interface {
	// --- Core Object Read/Write ---

	// ReadCommit retrieves a complete commit object by its SHA-1 hash.
	ReadCommit(ctx context.Context, hash string) (*Commit, error)

	// Commit creates a new commit object in the repository. It is an atomic operation.
	//   - parentHash: The hash of the commit this new commit will be on top of.
	//   - author: The author metadata for the commit.
	//   - message: The commit message.
	//   - graphData: A map where keys are named graph IRIs and values are the complete
	//     set of quads for that graph in the new state. Providing an empty slice for a
	//     graph IRI will delete that graph. Graphs not included in the map will be
	//     inherited from the parent commit.
	//   - sign: An optional callback function to perform GPG signing. If nil, the
	//     commit will be unsigned. The function receives the canonical commit data
	//     and should return an ASCII-armored signature.
	// It returns the hash of the newly created commit.
	Commit(ctx context.Context, parentHash string, author Author, message string, graphData map[string][]Quad, sign func(data []byte) (string, error)) (string, error)

	// --- Reference Management ---

	// SetReference creates or updates a reference (like a branch or tag) to point to a specific commit hash.
	// The name should be a full reference name, e.g., "refs/heads/main".
	SetReference(ctx context.Context, name string, hash string) error

	// GetReference retrieves the commit hash a specific full reference name points to.
	GetReference(ctx context.Context, name string) (string, error)

	// ResolveRef resolves a user-friendly name (e.g., "main", "v1.0", "HEAD", "a1b2c3d") to a full commit hash.
	ResolveRef(ctx context.Context, name string) (string, error)

	// ListReferences returns a list of all references matching a given prefix (e.g., "refs/heads/").
	ListReferences(ctx context.Context, prefix string) ([]Reference, error)

	// DeleteReference removes a reference from the repository.
	DeleteReference(ctx context.Context, name string) error

	// --- History & State Inspection ---

	// Log retrieves a slice of commits by walking the history backwards from a starting hash.
	Log(ctx context.Context, startHash string, limit int) ([]*Commit, error)

	// Blame annotates each quad in a named graph at a specific commit with the commit that last introduced it.
	// It returns a read-only channel from which the caller can stream the results. This is a
	// memory-efficient way to handle potentially large graphs. The channel will be closed when the operation is complete.
	Blame(ctx context.Context, graphIRI string, atCommitHash string) (<-chan BlameResult, error)

	// Diff generates the changes (additions/deletions) between the states of two commits.
	// It returns a read-only channel for streaming results to handle large diffs efficiently.
	// The channel will be closed when the operation is complete.
	Diff(ctx context.Context, fromCommitHash, toCommitHash string) (<-chan Change, error)

	// --- Advanced Operations ---

	// Merge attempts to perform a three-way merge.
	// It takes the commit hashes for the target branch head, the source branch head,
	// and their calculated common ancestor. If the merge is clean, it returns an empty
	// slice of conflicts and no error. If conflicts are detected, it returns a slice
	// of Conflict objects and no error, indicating a manual resolution is required.
	Merge(ctx context.Context, baseCommitHash, targetCommitHash, sourceCommitHash string) ([]Conflict, error)
	
	// Revert creates a new commit on top of a given branch head that is the inverse of a specified commit.
	// This provides a safe way to undo changes. Returns the hash of the new revert commit.
	Revert(ctx context.Context, branchHeadHash, commitToRevertHash string, author Author) (string, error)

	// Backup performs a full or incremental backup of the entire repository to a writer.
	// `sinceVersion` is obtained from a previous backup's manifest for incrementals. A value of 0
	// indicates a full backup. It returns a manifest with metadata about the completed backup.
	Backup(ctx context.Context, writer io.Writer, sinceVersion uint64) (*BackupManifest, error)

	// Restore populates a database from a backup stream. This is a destructive operation and
	// should be performed on an empty repository.
	Restore(ctx context.Context, reader io.Reader) error

	// Close closes the connection to the underlying database store(s) and releases any resources.
	// It must be called when the application is done with the Store instance.
	Close() error
}

// Open is the main entry point to the quadstore library.
// It initializes and returns a Store instance for a given repository path and namespace.
// The concrete implementation is in the internal/datastore package and is not exposed publicly.
func Open(ctx context.Context, opts OpenOptions) (Store, error) {
	// This function's body will be implemented in a separate, internal package.
	// It will call an internal constructor, e.g., `datastore.NewRepository(opts)`.
	// This is a common Go pattern to hide the concrete implementation type.
	panic("unimplemented")
}
