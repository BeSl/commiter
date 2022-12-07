package qworker

import (
	"commiter/internal/config"
	"commiter/internal/executor"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"os"
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
	Base64data string
	Name       string
	ExtID      string
	TypeProc   string
}

func NewQWorker(gitcfg *config.Gitlab) *QWorker {
	return &QWorker{
		GitConf: *gitcfg,
	}
}

func (qw *QWorker) ListenNewJob() error {
	if qw.shuttingDown() {
		return ErrWorkerClosed
	}
	var sleepMinute = 1

	for {
		worked := qw.itswork()
		if worked == false {
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

		time.Sleep(time.Minute * time.Duration(sleepMinute))
	}
	return nil

}
func saveFileRepository(dw *DataWork, cfg *config.Gitlab) error {
	data := dw.Base64data

	file, _ := os.Create(pathFileFromData(dw, cfg))
	defer file.Close()

	sDec, _ := base64.StdEncoding.DecodeString(data)
	_, err := file.Write(sDec)
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

	return fmt.Sprintf("%s/%s/%s/%s.%s",
		cfg.CurrPath,
		pathBase_ExtProcessor,
		pathType,
		dw.Name, extP)
}

func selectDataFromWork(db *sqlx.DB) (*DataWork, error) {
	//table id, extID, data, autor, textcommit, nameObject, typeObj, isComplete
	txtQuery := "SELECT TOP 1 qw.*, u.gitname FROM queue_work as qw left join users as u on u.id = qw.autor WHERE qw.isComplete = false order by qw.id ask"

	dw := DataWork{}
	err := db.Get(dw, txtQuery)
	if err != nil {
		return &DataWork{}, err
	}

	return &dw, nil
}

func createCommitDataProc(dw *DataWork) error {

	ex := executor.NewExecutor()
	err := ex.AddIndexFile()
	if err != nil {
		return err
	}

	err = ex.CommitRepo("testUser", "tst@tt.ru", "demo commit")

	if err != nil {
		return err
	}
	err = ex.PushToRepo()
	if err != nil {
		return err
	}
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
