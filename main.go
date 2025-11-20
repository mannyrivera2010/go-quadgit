// main.go
package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/spf13/cobra"
)

// --- 1. DATA MODEL ---
// These structs represent the Git-like objects we store.

type Commit struct {
	Tree      string    `json:"tree"` // SHA-1 hash of the tree object
	Parents   []string  `json:"parents"` // SHA-1 hashes of parent commits
	Author    string    `json:"author"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// A Tree is simplified to map a graph name to a blob hash containing its quads.
type Tree map[string]string

// A Blob is simply the content; in our case, a list of quad strings.
type Blob []string

// --- 2. STORE ---
// Manages all interaction with the BadgerDB database.

const (
	dbPath    = ".quad-db"
	indexPath = ".quad-db/index"
)

var db *badger.DB

// openDB opens the BadgerDB database in the .quad-db directory.
func openDB() (*badger.DB, error) {
	if db != nil {
		return db, nil
	}
	opts := badger.DefaultOptions(dbPath).WithLogger(nil) // Suppress Badger logger
	var err error
	db, err = badger.Open(opts)
	return db, err
}

// closeDB closes the database connection.
func closeDB() {
	if db != nil {
		db.Close()
	}
}

// writeObject serializes an object (Commit, Tree), computes its hash,
// and saves it to the database.
func writeObject(obj interface{}) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	hashBytes := sha1.Sum(data)
	hash := hex.EncodeToString(hashBytes[:])
	key := []byte("obj:" + hash)

	err = db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err == nil {
			return nil // Object already exists
		}
		if err != badger.ErrKeyNotFound {
			return err
		}
		return txn.Set(key, data)
	})
	return hash, err
}

// readCommit reads and deserializes a commit object from its hash.
func readCommit(hash string) (*Commit, error) {
	var commit Commit
	key := []byte("obj:" + hash)
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return fmt.Errorf("commit with hash %s not found", hash)
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &commit)
		})
	})
	return &commit, err
}

// setReference points a reference (like a branch or HEAD) to a commit hash.
func setReference(ref, hash string) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("ref:"+ref), []byte(hash))
	})
}

// getReference resolves a reference to a commit hash.
func getReference(ref string) (string, error) {
	var hash string
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("ref:" + ref))
		if err != nil {
			return fmt.Errorf("reference %s not found", ref)
		}
		return item.Value(func(val []byte) error {
			hash = string(val)
			return nil
		})
	})
	return hash, err
}

// resolveHead gets the commit hash that HEAD points to.
func resolveHead() (string, error) {
	headVal, err := getReference("HEAD")
	if err != nil {
		return "", err
	}
	// HEAD points to a branch ref, e.g., "ref:head:main"
	return getReference(strings.TrimPrefix(headVal, "ref:"))
}

// --- 3. CLI COMMANDS ---

var rootCmd = &cobra.Command{
	Use:   "quad-db",
	Short: "A git-like quad store CLI using BadgerDB",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Don't open DB for 'init' command if the directory doesn't exist yet
		if cmd.Name() == "init" {
			return nil
		}
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			return errors.New("repository not initialized, run 'quad-db init'")
		}
		_, err := openDB()
		return err
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		closeDB()
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new quad-db repository",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
			log.Fatal("Repository already initialized.")
		}
		os.Mkdir(dbPath, 0755)

		var err error
		db, err = openDB()
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}

		// 1. Create an empty tree
		emptyTree := make(Tree)
		treeHash, err := writeObject(emptyTree)
		if err != nil {
			log.Fatalf("Failed to create initial tree: %v", err)
		}

		// 2. Create the root commit
		rootCommit := Commit{
			Tree:      treeHash,
			Parents:   []string{}, // No parents
			Author:    "System",
			Message:   "Initial commit",
			Timestamp: time.Now(),
		}
		commitHash, err := writeObject(rootCommit)
		if err != nil {
			log.Fatalf("Failed to create root commit: %v", err)
		}

		// 3. Create the 'main' branch and point HEAD to it
		if err := setReference("head:main", commitHash); err != nil {
			log.Fatalf("Failed to create main branch: %v", err)
		}
		if err := setReference("HEAD", "ref:head:main"); err != nil {
			log.Fatalf("Failed to set HEAD: %v", err)
		}

		fmt.Printf("Initialized empty quad-db repository in %s\n", dbPath)
	},
}

var addCmd = &cobra.Command{
	Use:   "add <file.nq>",
	Short: "Add quads from a file to the staging area",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		content, err := os.ReadFile(args[0])
		if err != nil {
			log.Fatalf("Failed to read file %s: %v", args[0], err)
		}

		// Simple staging: append to an index file.
		f, err := os.OpenFile(indexPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open index: %v", err)
		}
		defer f.Close()

		if _, err := f.WriteString(string(content)); err != nil {
			log.Fatalf("Failed to write to index: %v", err)
		}
		fmt.Printf("Staged changes from %s\n", args[0])
	},
}

var commitCmd = &cobra.Command{
	Use:   "commit -m <message>",
	Short: "Record staged changes to the repository",
	Run: func(cmd *cobra.Command, args []string) {
		message, _ := cmd.Flags().GetString("message")
		if message == "" {
			log.Fatal("Commit message is required. Use -m.")
		}

		// 1. Read staged quads from index
		stagedQuads, err := os.ReadFile(indexPath)
		if err != nil || len(stagedQuads) == 0 {
			log.Fatal("Nothing to commit. Stage changes with 'add' first.")
		}

		// 2. Create a blob from the staged quads
		quads := strings.Split(strings.TrimSpace(string(stagedQuads)), "\n")
		blobHash, err := writeObject(quads)
		if err != nil {
			log.Fatalf("Failed to create blob object: %v", err)
		}

		// 3. Create a new tree
		// Simplified: This commit will only contain the new blob.
		// A real implementation would merge with the parent's tree.
		newTree := Tree{
			"default": blobHash, // Using a default graph name
		}
		treeHash, err := writeObject(newTree)
		if err != nil {
			log.Fatalf("Failed to create tree object: %v", err)
		}

		// 4. Get parent commit
		parentHash, err := resolveHead()
		if err != nil {
			log.Fatalf("Could not resolve HEAD: %v", err)
		}

		// 5. Create the new commit object
		newCommit := Commit{
			Tree:      treeHash,
			Parents:   []string{parentHash},
			Author:    "user@example.com", // Should be configurable
			Message:   message,
			Timestamp: time.Now(),
		}
		commitHash, err := writeObject(newCommit)
		if err != nil {
			log.Fatalf("Failed to write commit object: %v", err)
		}

		// 6. Update the branch reference
		headRef, _ := getReference("HEAD")
		if err := setReference(strings.TrimPrefix(headRef, "ref:"), commitHash); err != nil {
			log.Fatalf("Failed to update branch reference: %v", err)
		}

		// 7. Clear the index
		os.Truncate(indexPath, 0)

		fmt.Printf("[%s] %s\n", commitHash[:7], message)
	},
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit history",
	Run: func(cmd *cobra.Command, args []string) {
		hash, err := resolveHead()
		if err != nil {
			log.Fatalf("Could not resolve HEAD: %v", err)
		}

		for {
			commit, err := readCommit(hash)
			if err != nil {
				log.Fatalf("Failed to read commit history: %v", err)
			}

			fmt.Printf("commit %s\n", hash)
			fmt.Printf("Author: %s\n", commit.Author)
			fmt.Printf("Date:   %s\n", commit.Timestamp.Format(time.RFC1123Z))
			fmt.Printf("\n\t%s\n\n", commit.Message)

			if len(commit.Parents) == 0 {
				break
			}
			hash = commit.Parents[0] // Follow the first parent
		}
	},
}

func main() {
	// Add commands to root
	rootCmd.AddCommand(initCmd, addCmd, logCmd)

	// Add flags
	commitCmd.Flags().StringP("message", "m", "", "Commit message")
	rootCmd.AddCommand(commitCmd)

	// Execute the CLI
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}