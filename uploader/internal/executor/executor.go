package executor

import (
	"fmt"
	"os/exec"
)

type Executor struct {
}

func NewExecutor() *Executor {
	return &Executor{}
}

func (ex *Executor) Check_env() error {
	err := check_oscript()
	if err != nil {
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

func (ex *Executor) CommitRepo(UserName, UserEmail, textcommit string) error {

	cmdCommit := "git commit --author=\"%s <%s>\" -m \"%s\""

	err := system_exec(fmt.Sprintf(cmdCommit, UserName, UserEmail, textcommit))
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
