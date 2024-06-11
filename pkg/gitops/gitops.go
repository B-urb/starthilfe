package gitops

import (
	"context"
	"errors"
	"fmt"
	"github.com/B-urb/starthilfe/pkg/projectconfig"
	"github.com/B-urb/starthilfe/pkg/state"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// AddSubtrees manages the addition of subtrees from a configuration.
func AddSubtrees(configPath string) error {
	_, err := git.PlainOpen("/Users/bjornurban/code/testing-stuff/testproject-starthilfe")
	if err != nil {
		return err
	}

	//if err := ensureGitRepoReady(repo); err != nil {
	//	return err
	//}

	cfg, err := projectconfig.LoadConfig(configPath)
	if err != nil {
		return err
	}

	return addSubtreesWithExec(cfg)
}

// UpdateSubtrees manages the updating of subtrees based on the state file.
func UpdateSubtrees(stateFilePath string) error {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return err
	}

	state, err := state.LoadState(stateFilePath) // Corrected to use projectconfig package
	if err != nil {
		return err
	}
	return updateSubtreesWithExec(repo, state)
}

// ensureGitRepoReady checks if the repository is ready for operations.
func ensureGitRepoReady(repo *git.Repository) error {
	_, err := repo.Head()
	if err != nil {
		return fmt.Errorf("repository does not have any commits yet: %v", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	status, err := w.Status()
	if err != nil {
		return err
	}

	if !status.IsClean() {
		return fmt.Errorf("working tree has modifications or untracked files")
	}

	return nil
}

func addSubtreeWithGoGit(repo *git.Repository, cfg *projectconfig.Config) error {
	// Create the remote if not exists
	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "actions",
		URLs: []string{"TEST"},
	})
	if err != nil {
		return err
	}

	// Fetch from the remote
	remote.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*"},
		Depth:    1,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	// Get the commit objects involved in the merge
	_, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	ref, err := repo.Reference(plumbing.ReferenceName("refs/remotes/actions/main"), true)
	if err != nil {
		return err
	}

	_, err = repo.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Merge the commits
	err = w.Pull(&git.PullOptions{
		RemoteName:    "actions",
		ReferenceName: plumbing.ReferenceName("refs/heads/main"),
		Auth:          nil, // Add authentication details if necessary
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	return nil
}

func addSubtreesWithExec(config *projectconfig.Config) error {

	for _, repo := range config.Repos {
		for _, path := range repo.Paths {
			slog.Info("Adding subtree", "path", path)
			err := DownloadFromGit(repo.URL, repo.Branch, "/Users/bjornurban/code/testing-stuff/testproject-starthilfe", path)
			if err != nil {
				slog.Error("Error downloading from git %s", "path", path, err.Error())
			}
			//cmd := exec.Command("git", "subtree", "add", "--prefix", path, repo.URL, repo.Branch, "--squash")
			//cmd.Dir = "/Users/bjornurban/code/testing-stuff/testproject-starthilfe"
			//cmd.Stdout = os.Stdout
			//cmd.Stderr = os.Stderr
			//err := cmd.Run()
			//if err != nil {
			//	if errors.Is(err, os.ErrExist) {
			//		// If error is because the directory already exists, continue without error
			//	} else {
			//		// If the error is due to another reason, return it
			//		return err
			//	}
			//}
		}
	}
	// Get the latest commit hash
	//TODO: get the commits of the added files and not the target. Use them for update
	//hashCmd := exec.Command("git", "rev-parse", "HEAD")
	//hashCmd.Stderr = os.Stderr
	//output, err := hashCmd.Output()
	//if err != nil {
	//	return err
	//}

	// Load or initialize state
	stateData, err := state.LoadState(".starthilfe_state")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			stateData = &state.State{}
		} else {
			return err
		}
	}

	// Add new entry
	stateData.Entries = append(stateData.Entries, struct {
		Commit string `yaml:"commit"`
		Path   string `yaml:"path"`
	}{
		Commit: "Haerwre",
		Path:   config.Repos[0].Paths[0],
	})

	return stateData.SaveState(".starthilfe_state")
}

// updateSubtreesWithExec updates the subtrees using the exec command to interface directly with the system's git.
func updateSubtreesWithExec(repo *git.Repository, state *state.State) error {
	for _, entry := range state.Entries {
		slog.Info("Updating subtree", "path", entry.Path)
		cmd := exec.Command("git", "subtree", "pull", "--prefix", entry.Path, "main", "--squash")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// DownloadFromGit downloads a specific file or directory from a repository.
func DownloadFromGit(repoURL, branch, targetDir, repoPath string) error {
	// Clone the repository into memory to access the commit tree
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           repoURL,
		Depth:         1, // Use depth 1 for a shallow clone
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to clone: %v", err)
	}

	// Get the HEAD commit
	ref, err := r.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %v", err)
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return fmt.Errorf("failed to get commit: %v", err)
	}

	// Get the tree from the commit and the path
	tree, err := commit.Tree()
	if err != nil {
		return fmt.Errorf("failed to get tree: %v", err)
	}

	// Get the specific entry at the given path
	entry, err := tree.FindEntry(repoPath)
	if err != nil {
		return fmt.Errorf("failed to find entry at path %s: %v", repoPath, err)
	}

	// Process the entry (file or directory)
	return processEntry(entry, tree, targetDir)
}

func processEntry(entry *object.TreeEntry, tree *object.Tree, targetDir string) error {

	if entry.Mode == filemode.Dir {
		// It's a directory, recurse into it
		subtree, err := tree.Tree(entry.Name)
		if err != nil {
			return fmt.Errorf("failed to access subtree: %v", err)
		}
		return downloadDirectory(subtree, filepath.Join(targetDir, entry.Name))
	} else {
		file, err := tree.TreeEntryFile(entry)
		if err != nil {
			slog.Error("error getting file: %v", err)
			return err
		}
		// It's a file, download it
		return downloadFile(file, targetDir)
	}
}

func downloadDirectory(tree *object.Tree, targetDir string) error {
	// Ensure the directory exists
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	// Iterate over the contents
	return tree.Files().ForEach(func(f *object.File) error {
		return downloadFile(f, targetDir)
	})
}

func downloadFile(file *object.File, targetDir string) error {
	reader, err := file.Reader()
	if err != nil {
		return fmt.Errorf("failed to open file reader: %v", err)
	}
	defer reader.Close()

	filePath := filepath.Join(targetDir, file.Name)
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, reader)
	return err
}
