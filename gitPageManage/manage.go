package gitPageManage

import (
	"errors"
	"github.com/RellwNote/RellwNote/config"
	"github.com/RellwNote/RellwNote/log"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
	"path/filepath"
	"strings"
)

var repo *git.Repository
var worktree *git.Worktree
var auth transport.AuthMethod

func init() {
	isExist := isInitRepo(config.LibraryPath)
	if !isExist {
		Clone()
	}
	repo, err := git.PlainOpen(config.LibraryPath)
	if err != nil {
		log.Error.Println(err)
	}
	worktree, err = repo.Worktree()
	if err != nil {
		log.Error.Println(err)
	}
	if config.GitPushUseToken {
		auth = &http.TokenAuth{
			Token: config.GitToken,
		}
	} else {
		auth = &http.BasicAuth{
			Username: config.GitUserName,
			Password: config.GitPassword,
		}
	}

}

func isInitRepo(path string) bool {
	headPath := filepath.Join(path, ".git", "HEAD")
	_, err := os.Stat(headPath)
	return err == nil
}
func Clone() {
	_, err := git.PlainClone(config.LibraryPath, false, &git.CloneOptions{
		URL:      config.GitUrl,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Error.Println(err)
	}
}

func Pull() (status int) {
	err := worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		log.Error.Println(err)
		return -1
	}
	return 0
}

func Commit(files []string, commitMsg string) (status int) {
	_, err := worktree.Add(strings.Join(files, " "))
	if err != nil {
		log.Error.Println(err)
	}

	_, err = worktree.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name: commitMsg,
		},
	})
	if err != nil {
		log.Error.Println(err)
		return -1
	}

	return 0
}

func Push(commitStr string) (status int) {
	err := repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	if err != nil {
		log.Error.Println(err)
		return -1
	}
	return 0
}

func CheckChangesFile() []string {
	changeFiles := make([]string, 0)
	wrStatus, err := worktree.Status()
	if err != nil {
		log.Error.Println(err)
	}

	for file, stat := range wrStatus {
		if stat.Worktree != git.Unmodified {
			changeFiles = append(changeFiles, file)
		}
	}
	return changeFiles
}
