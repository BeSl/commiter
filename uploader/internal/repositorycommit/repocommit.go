package repositorycommit

import (
	"commiter/internal/config"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/jmoiron/sqlx"
	"go-micro.dev/v4/runtime/local/source/git"
)

type GitWorker struct {
	Db      *sqlx.DB
	GitCfg  *config.Gitlab
	GitRepo *git.Repository
}

func NewGitWorker(cfgGit *config.Gitlab, g *git.Repository) *GitWorker {
	return &GitWorker{
		GitCfg:  cfgGit,
		GitRepo: g,
	}

}

func (gw *GitWorker) CloneRepo() error {

	_, err := git.PlainClone(gw.GitCfg.CurrPath, false, &git.CloneOptions{
		URL:      gw.GitCfg.Url,
		Progress: os.Stdout,
	})

	return err

}

func (gw *GitWorker) GitPull() error {
	r, err := git.PlainOpen(gw.GitCfg.CurrPath)

	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	return w.Pull(&git.PullOptions{RemoteName: "origin"})

}

func (gw *GitWorker) CheckOut(branchName string) error {
	w, err := gw.GitRepo.Worktree()
	if err != nil {
		return err
	}

	w.Checkout(&git.CheckoutOptions{Branch: plumbing.ReferenceName(branchName)})
	return nil
}

func (gw *GitWorker) CreateCommit() error {

	return nil

}
