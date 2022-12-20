package executor

import (
	"commiter/internal/config"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

type Executor struct {
	cm          *exec.Cmd
	PathProject string
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

	path, err := os.Getwd()

	_, err = os.Stat(path + `\` + gitCfg.CurrPath)
	if err == nil {
		return err
	}

	if os.IsNotExist(err) == false {
		return nil
	}

	log.Info().Msg("Cloning repo strart...")

	pathRepo := gitCfg.Url + "/" + gitCfg.Project_url + ".git"
	txcmd := "git clone %s %s"

	allCommand := fmt.Sprintf(txcmd, pathRepo, gitCfg.CurrPath)
	log.Info().Msg(allCommand)

	log.Info().Msg("Current path " + path)

	arg := strings.Split(allCommand, " ")
	cm := exec.Command(arg[0], arg[1:]...)
	cm.Dir = path
	b, err := cm.CombinedOutput()

	if err != nil {
		log.Error().Err(err).Msg("GIT CLONE_Error")
		return err
	}

	log.Info().Msg("Cloning repo Finish")
	log.Info().Msg("Cloning repo Finish " + string(b))
	
	
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
	arg := strings.Split(str_cmd, " ")
	cm := exec.Command(arg[0], arg[1:]...)
	cm.Dir = ex.PathProject

	b, err := cm.CombinedOutput()
	if err != nil {
		cmdText := "git status"
		err := ex.System_ex(cmdText)
		if err != nil {
			log.Error().Err(err).Msg("Error CombinedOutput: ")
			return err
		}
	} else {
		log.Info().Msg("Done commit! : " + string(b))
	}
	return nil
}

func system_exec(str_cmd string) error {
	//run command system
	c := exec.Command("cmd", "/C", str_cmd)

	if err := c.Run(); err != nil {
		log.
			Error().
			Err(err).
			Msg("Error: " + str_cmd)

		return err
	} else {
		log.
			Info().
			Msg("Done! " + str_cmd)
		return nil
	}

}
