package qworker

import (
	"commiter/internal/config"
	"commiter/internal/repositorycommit"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"os"
	"sync/atomic"
	"time"

	"github.com/go-git/go-git/v5"
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

func NewQWorker(gitcfg *config.Gitlab, db *sqlx.DB) *QWorker {
	return &QWorker{
		GitConf: *gitcfg,
		db:      db,
	}
}

func (qw *QWorker) ListenNewJob() error {
	if qw.shuttingDown() {
		return ErrWorkerClosed
	}
	var sleepMinute = 1
	tempC := 0

	for {
		tempC++
		if tempC > 2 {
			break
		}
		dw, err := selectDataFromWork(qw.db)
		if err != nil {
			qw.sendError(err)
			continue
		}

		err = saveFileRepository(dw, &qw.GitConf)
		if err != nil {
			qw.sendError(err)
			continue
		}

		createCommitDataProc(dw)
		txtQ := `UPDATE 
		commit_tasks AS tgt 
		SET  processed=true  WHERE tgt.id=$1`

		_, err = qw.db.Exec(txtQ, dw.ID)
		if err != nil {
			return err
		}
		time.Sleep(time.Minute * time.Duration(sleepMinute))
	}
	return nil

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

func saveFileRepository(dw *DataWork, cfg *config.Gitlab) error {

	check, err := PathRepoExist(cfg.CurrPath)

	if err != nil {
		return err
	}

	if check == false {
		return os.ErrProcessDone
	}
	
	r, err := git.PlainOpen(cfg.CurrPath)
	
	gw := repositorycommit.NewGitWorker(cfg, r)
	err = gw.GitPull()
	
	if err != nil {
		return err
	}

	gw.CheckOut("develop")
	data := dw.Base64data
	pathFile:=pathFileFromData(dw, cfg)
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
	
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	
	_,err = w.Add(pathFile)
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
		dw.Name, extP)
	return res
}

func selectDataFromWork(db *sqlx.DB) (*DataWork, error) {
	//table id, extID, data, autor, textcommit, nameObject, typeObj, isComplete
	// txtQ := "SELECT u.gitlogin gitlogin, u.name as UserName,ct.name as name,ct.type as type,ct.base64data as base64data,coalesce(ct.textcommit, 'not text') as commit FROM commit_tasks ct left join users u on u.id= ct.userid WHERE ct.processed =false ORDER BY ct.id DESC LIMIT 1"
	txtQ := `SELECT 
				u.gitlogin as gitlogin,
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
			ORDER BY ct.id DESC 
			LIMIT 1`

	dw := DataWork{}
	err := db.Get(&dw, txtQ)
	if err != nil {
		return &DataWork{}, err
	}

	return &dw, nil
}

func createCommitDataProc(dw *DataWork) error {

	git.PlainOpen(dw.)
	return nil
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
