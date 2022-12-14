package comittworker

import (
	"commiter/internal/config"
	"commiter/internal/errorwrapper"
	"commiter/internal/executor"
	"commiter/internal/model"
	"commiter/internal/storage"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
)

const TimeSleepMinute = 1

var ErrWorkerClosed = errors.New("qwork: worker closed")

var (
	pathBase_ExtProcessor = "DataProcessorsExt"
	path_extReport        = "Отчет"
	path_extProc          = "Обработка"
)

type CommitCreator struct {
	DB         *sqlx.DB
	Bot        *tgbotapi.BotAPI
	GitCfg     *config.Gitlab
	inShutdown atomicBool
}

type atomicBool int32

func (b *atomicBool) isSet() bool { return atomic.LoadInt32((*int32)(b)) != 0 }

func (cm *CommitCreator) Shutdown(ctx context.Context) error {
	return nil
}

func (cm *CommitCreator) shuttingDown() bool {
	return cm.inShutdown.isSet()
}

func (cm *CommitCreator) itswork() bool {
	return cm.inShutdown.isSet()
}

func NewCommitCreator(db *sqlx.DB, bot *tgbotapi.BotAPI, cfg *config.Gitlab) *CommitCreator {
	return &CommitCreator{
		DB:     db,
		Bot:    bot,
		GitCfg: cfg,
	}
}

func (cc *CommitCreator) ListenNewTasks() error {

	st := storage.NewStorage(cc.DB, cc.Bot, cc.GitCfg)
	adminUser, err := st.FindAdmin()
	if err != nil {
		return err
	}

	for {
		dataCommit, err := st.FindLastCommit()
		if err != nil {
			errorwrapper.HandError(err, cc.DB, cc.Bot, adminUser.TGid)
			time.Sleep(time.Minute * time.Duration(TimeSleepMinute))
			continue
		}

		err = createCommit(dataCommit, st)
		if err != nil {
			errorwrapper.HandError(err, cc.DB, cc.Bot, adminUser.TGid)
			time.Sleep(time.Minute * time.Duration(TimeSleepMinute))
			continue
		}

		time.Sleep(time.Minute * time.Duration(TimeSleepMinute))
	}

}

func createCommit(dc *model.DataWork, st *storage.Storage) error {

	err := saveFileRepository(dc, st.GitConf)
	if err != nil {
		return err
	}
	err = commitRepo(dc, st.GitConf)
	if err != nil {
		return err
	}

	txtQ := `UPDATE
	 	commit_tasks AS tgt
	 	SET  processed=true  WHERE tgt.id=$1`

	_, err = st.DB.Exec(txtQ, dc.ID)
	if err != nil {
		return err
	}
	return nil
}

func commitRepo(dw *model.DataWork, cfg *config.Gitlab) error {

	arg1 := strings.Split("git status", " ")
	cm := exec.Command(arg1[0], arg1[1:]...)
	cm.Dir = cfg.CurrPath
	stat, err := cm.CombinedOutput()
	if err != nil {
		return err

	} else {
		fmt.Println("res = " + string(stat))
		fmt.Println("Done! status")
	}

	ex := executor.NewExecutor()
	cmdText := "git add *"
	err = ex.System_ex(cmdText)
	if err != nil {
		return err
	}

	var result = "git commit --author=%s<%s> -m"
	cmdText = fmt.Sprintf(result, dw.UserName, dw.GitLogin)

	arg := strings.Split(cmdText, " ")
	arg = append(arg, dw.Commit)
	cm = exec.Command(arg[0], arg[1:]...)
	cm.Dir = cfg.CurrPath

	b, err := cm.CombinedOutput()
	if err != nil {
		cmdText := "git status"
		err := ex.System_ex(cmdText)
		if err != nil {
			fmt.Println("Error CombinedOutput: ", err.Error())
			return err
		}
	} else {
		fmt.Println("res = " + string(b))
		fmt.Println("Done commit!")
	}

	cmdText = "git push -u origin develop"
	err = ex.System_ex(cmdText)
	if err != nil {
		return err
	}

	return nil
}

func saveFileRepository(dw *model.DataWork, cfg *config.Gitlab) error {

	check, err := PathRepoExist(cfg.CurrPath)

	if err != nil {
		return err
	}

	if check == false {
		return os.ErrProcessDone
	}

	ex := executor.NewExecutor()

	cmdText := "git reset"
	err = ex.System_ex(cmdText)
	if err != nil {
		return err
	}
	cmdText = "git checkout develop"
	err = ex.System_ex(cmdText)
	if err != nil {
		return err
	}

	cmdText = "git pull"
	err = ex.System_ex(cmdText)
	if err != nil {
		return err
	}

	data := dw.Base64data
	pathFile := pathFileFromData(dw, cfg)
	file, _ := os.Create(pathFile)

	defer file.Close()

	sDec, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}

	_, err = file.Write(sDec)
	if err != nil {
		return err
	}
	return nil
}

func pathFileFromData(dw *model.DataWork, cfg *config.Gitlab) string {

	pathType := path_extProc
	extP := "epf"

	if dw.TypeProc == "Отчет" {
		pathType = path_extReport
		extP = "erf"
	}

	res := fmt.Sprintf("%s/%s/%s/%s.%s",
		cfg.CurrPath,
		pathBase_ExtProcessor,
		pathType,
		strings.Replace(dw.Name, "/", "_", -1),
		extP)
	return res
}

func PathRepoExist(path string) (bool, error) {

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
