package qworker

import (
	"commiter/internal/api"
	"commiter/internal/config"
	"commiter/internal/errorwrapper"
	"commiter/internal/executor"
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
)

var ErrWorkerClosed = errors.New("qwork: worker closed")

var (
	pathBase_ExtProcessor = "DataProcessorsExt"
	path_extReport        = "Отчет"
	path_extProc          = "Обработка"
)

type TypeDataProcess struct {
	extFile  string
	pathRepo string
}
type ExtendetDataProcessors struct {
	extID      string
	name       string
	dataBase64 string
	Type       TypeDataProcess
}
type GitConfig struct {
}

type atomicBool int32

func (b *atomicBool) isSet() bool { return atomic.LoadInt32((*int32)(b)) != 0 }

type QWorker struct {
	GitConf     config.Gitlab
	ConnContext func(ctx context.Context, c net.Conn) context.Context
	db          *sqlx.DB
	inShutdown  atomicBool
	Bot         *tgbotapi.BotAPI
}

type DataWork struct {
	Base64data string `db:"base64data"`
	Name       string `db:"name"`
	ID         int64  `db:"id"`
	TypeProc   string `db:"type"`
	UserName   string `db:"username"`
	GitLogin   string `db:"gitlogin"`
	Commit     string `db:"commit"`
}

func NewQWorker(gitcfg *config.Gitlab, db *sqlx.DB, bot *tgbotapi.BotAPI) *QWorker {
	return &QWorker{
		GitConf: *gitcfg,
		db:      db,
		Bot:     bot,
	}
}

func (qw *QWorker) ListenNewJob() error {
	if qw.shuttingDown() {
		return ErrWorkerClosed
	}
	var sleepMinute = 1
	//	tempC := 0
	extConnec := &api.ExternalConnection{DB: qw.db, Bot: qw.Bot}
	for {

		dw, err := selectDataFromWork(qw.db)

		if err != nil {
			errorwrapper.HandError(err, extConnec, fmt.Sprintf("%s - %s", dw.Name, dw.Commit))
			time.Sleep(time.Minute * time.Duration(sleepMinute))
			continue
		}

		err = saveFileRepository(dw, &qw.GitConf)
		if err != nil {
			errorwrapper.HandError(err, extConnec, fmt.Sprintf("%s - %s", dw.Name, dw.Commit))
			time.Sleep(time.Minute * time.Duration(sleepMinute))
			continue
		}
		err = commitRepo(dw, &qw.GitConf)
		if err != nil {
			errorwrapper.HandError(err, extConnec, fmt.Sprintf("%s - %s", dw.Name, dw.Commit))
			time.Sleep(time.Minute * time.Duration(sleepMinute))
			continue
		}
		txtQ := `UPDATE 
		commit_tasks AS tgt 
		SET  processed=true  WHERE tgt.id=$1`

		_, err = qw.db.Exec(txtQ, dw.ID)
		if err != nil {
			return errorwrapper.HandError(err, extConnec, "")
		}
		time.Sleep(time.Minute * time.Duration(sleepMinute))
	}
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

func saveFileRepository_old(dw *DataWork, cfg *config.Gitlab) error {

	// check, err := PathRepoExist(cfg.CurrPath)

	// if err != nil {
	// 	return err
	// }

	// if check == false {
	// 	return os.ErrProcessDone
	// }

	// r, err := git.PlainOpen(cfg.CurrPath)

	// gw := repositorycommit.NewGitWorker(cfg, r)
	// err = gw.GitPull()

	// if err != nil {
	// 	return err
	// }

	// gw.CheckOut("develop")
	// data := dw.Base64data
	// pathFile := pathFileFromData(dw, cfg)
	// file, _ := os.Create(pathFile)

	// defer file.Close()

	// sDec, err := base64.StdEncoding.DecodeString(data)
	// if err != nil {
	// 	return err
	// }

	// _, err = file.Write(sDec)
	// if err != nil {
	// 	return err
	// }

	// w, err := r.Worktree()
	// if err != nil {
	// 	return err
	// }

	// _, err = w.Add(pathFile)
	// if err != nil {
	// 	return err
	// }

	// commit, err := w.Commit(dw.Commit, &git.CommitOptions{
	// 	Author: &object.Signature{
	// 		Name:  dw.UserName,
	// 		Email: dw.GitLogin,
	// 		When:  time.Now(),
	// 	},
	// })

	// if err != nil {
	// 	return err
	// }
	// err = r.Push(&git.PushOptions{})
	// if err != nil {
	// 	return err
	// }
	// log.Info().Msgf("Create commit " + commit.String())

	return nil
}

func commitRepo(dw *DataWork, cfg *config.Gitlab) error {

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

func saveFileRepository(dw *DataWork, cfg *config.Gitlab) error {

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

func pathFileFromData(dw *DataWork, cfg *config.Gitlab) string {

	pathType := path_extProc
	extP := "epf"

	if dw.TypeProc == "Отчет" {
		pathType = path_extReport
		extP = "erf"
	}
	//"D:/tempRepo_/DataProcessorsExt/Обработка/Сайт выгрузка картинок.epf"

	res := fmt.Sprintf("%s/%s/%s/%s.%s",
		cfg.CurrPath,
		pathBase_ExtProcessor,
		pathType,
		strings.Replace(dw.Name, "/", "_", -1),
		extP)
	return res
}

func selectDataFromWork(db *sqlx.DB) (*DataWork, error) {
	txtQ := `SELECT 
				coalesce(u.gitlogin, '') as gitlogin,
		 		u.name as username,
				ct.name as name,
				ct.id as id,
				ct.type as type,
				ct.base64data as base64data,
		 		coalesce(ct.textcommit, 'not text') as commit 
		 	FROM 
		 		commit_tasks ct 
				left join users u 
				on u.id= ct.userid 
		 	WHERE 
		 		ct.processed =false 
			ORDER BY ct.id 
			LIMIT 1`

	dw := DataWork{}
	err := db.Get(&dw, txtQ)
	if len(dw.GitLogin) == 0 {
		return &DataWork{}, errors.New("Не заполнен пользователь " + dw.UserName)
	}

	if err == sql.ErrNoRows {
		return &DataWork{}, nil
	}

	if err != nil {
		return &DataWork{}, err
	}

	return &dw, nil
}

func (qw *QWorker) Shutdown(ctx context.Context) error {
	return nil
}

func (qw *QWorker) shuttingDown() bool {
	return qw.inShutdown.isSet()
}

func (qw *QWorker) itswork() bool {
	return qw.inShutdown.isSet()
}

func (qw *QWorker) sendError(e error) {

}
