package executor

import (
	"commiter/internal/config"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Executor struct {
	cm *exec.Cmd
}

func New() *Executor {
	return &Executor{}
}

func (ex *Executor) Check_env() error {
	err := check_oscript()
	if err != nil {
		return err
	}
	return nil
}

func (ex *Executor) CloneRepo(gitCfg *config.Gitlab) error {

	_, err := os.Stat(gitCfg.CurrPath)
	if err == nil {
		return err
	}

	if os.IsNotExist(err) {
		return err
	}

	pathRepo := gitCfg.Url + "/" + gitCfg.Project_url + ".git"
	txcmd := "git clone %s %s"

	arg := strings.Split(fmt.Sprintf(txcmd, pathRepo, gitCfg.CurrPath), " ")
	cm := exec.Command(arg[0], arg[1:]...)
	cm.Dir = gitCfg.CurrPath
	b, err := cm.CombinedOutput()

	if err != nil {
		fmt.Println("GIT CLONE_Error CombinedOutput: ", b)
		return err
	}

	return nil
}

func check_oscript() error {

	str_command := "oscript -v"
	err := system_exec(str_command)

	if err != nil {
		return err
	}

	return nil
}

func (ex *Executor) AddIndexFile() error {

	addIndex := "git add *"
	err := system_exec(addIndex)
	if err != nil {
		return err
	}

	return nil
}

func (ex *Executor) PushToRepo() error {
	cmdPush := "git push origin"
	err := system_exec(cmdPush)
	if err != nil {
		return err
	}

	return nil

}

func (ex *Executor) System_ex(str_cmd string) error {
	//run command system
	cmdPath := "D:/tempRepo_"
	cm := exec.Command("cmd", "/C", str_cmd)
	cm.Dir = cmdPath
	if err := cm.Run(); err != nil {
		fmt.Println("Error: ", err)
		return err
	} else {
		fmt.Println("Done!")
		return nil
	}

}

func system_exec(str_cmd string) error {
	//run command system
	c := exec.Command("cmd", "/C", str_cmd)

	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
		return err
	} else {
		fmt.Println("Done!")
		return nil
	}

}
