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
	"github.com/rs/zerolog/log"
)

const TimeSleepMinute = 1
const TimeSleep30Minute = 30

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

	st := storage.NewStorage(cc.DB, cc.GitCfg)
	adminUser, err := st.FindAdmin()
	if err != nil {
		return err
	}

	for {
		dataCommit, err := st.FindLastCommit()
		if err != nil {
			errorwrapper.HandError(err, cc.DB, cc.Bot, adminUser.TGid)
			if strings.Contains(err.Error(), "sql: no rows in result set") {
				time.Sleep(time.Minute * time.Duration(TimeSleepMinute))
			} else {
				errorwrapper.HandError(err, cc.DB, cc.Bot, adminUser.TGid)
				time.Sleep(time.Minute * time.Duration(TimeSleepMinute))
			}
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

	fileName, err := saveFileRepository(dc, st.GitConf)
	if err != nil {
		return err
	}
	err = commitRepo(fileName, dc, st.GitConf)
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

func commitRepo(fileName string, dw *model.DataWork, cfg *config.Gitlab) error {

	arg1 := strings.Split("git status", " ")
	cm := exec.Command(arg1[0], arg1[1:]...)
	cm.Dir = cfg.CurrPath
	log.Info().Msg("Path: " + cfg.CurrPath)
	stat, err := cm.CombinedOutput()
	if err != nil {
		return err

	} else {
		log.
			Info().
			Msg("git status! " + string(stat))
	}
	path, err := os.Getwd()

	ex := executor.New()
	ex.PathProject = path + `\` + cfg.CurrPath
	cmdText := "git add -A"
	cmdText = strings.Replace(cmdText, "\\", "/", -1)
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
			log.Error().Err(err).Msg("Error CombinedOutput: ")
			return err
		}
	} else {
		log.Info().Msg("Done commit! : " + string(b))
	}

	cmdText = "git push -u origin develop"
	err = ex.System_ex(cmdText)
	if err != nil {
		log.Error().Err(err).Msg(cmdText)
		return err
	}

	return nil
}

func saveFileRepository(dw *model.DataWork, cfg *config.Gitlab) (string, error) {

	check, err := PathRepoExist(cfg.CurrPath)

	if err != nil {
		return "", err
	}

	if check == false {
		return "", os.ErrProcessDone
	}

	path, err := os.Getwd()
	ex := executor.New()
	ex.PathProject = path + `\` + cfg.CurrPath

	cmdText := "git reset"
	err = ex.System_ex(cmdText)
	if err != nil {
		return "", err
	}
	cmdText = "git checkout develop"
	err = ex.System_ex(cmdText)
	if err != nil {
		return "", err
	}

	cmdText = "git pull"
	err = ex.System_ex(cmdText)
	if err != nil {
		return "", err
	}

	data := dw.Base64data
	pathFile := pathFileFromData(dw, cfg)
	file, _ := os.Create(pathFile.FullName)
	log.Info().Msg("Create " + pathFile.FullName)
	defer file.Close()

	sDec, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	_, err = file.Write(sDec)
	if err != nil {
		return "", err
	}
	return pathFile.Name, nil
}

type FileInfo struct {
	FullName string
	Name     string
}

func pathFileFromData(dw *model.DataWork, cfg *config.Gitlab) FileInfo {

	pathType := path_extProc
	extP := "epf"

	if dw.TypeProc == "Отчет" {
		pathType = path_extReport
		extP = "erf"
	}

	fInf := FileInfo{}

	res := fmt.Sprintf("%s/%s/%s/%s.%s",
		cfg.CurrPath,
		pathBase_ExtProcessor,
		pathType,
		strings.Replace(dw.Name, "/", "_", -1),
		extP)
	fInf.FullName = res
	fInf.Name = fmt.Sprintf("%s.%s", strings.Replace(dw.Name, "/", "_", -1), extP)
	return fInf
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
