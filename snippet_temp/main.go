package main

import (
	"fmt"
	"os/exec"
)

var (
	branch_commit = "develop"
	url_git       = "https://ee.ru"
	project       = "t/b"
	path_clone    = "D:/../tempRepo_"
)

func main() {
	sys_clone_repo()
}

func gitStatus() {

	str_command := "git status"
	system_exec(str_command)

}
func sys_clone_repo() {
	//если нет репозитория - клонируем, переходим наветку
	//иначе перейти на ветку и pull из origin

	fullCommand := fmt.Sprintf("git clone %s/%s.git %s", url_git, project, path_clone)

	system_exec(fullCommand)
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
