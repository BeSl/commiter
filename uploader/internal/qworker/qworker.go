package qworker

import (
	"commiter/internal/config"
	"context"
	"encoding/base64"
	"errors"
	"net"
	"os"
	"sync/atomic"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrWorkerClosed = errors.New("qwork: worker closed")

type GitConfig struct {
}

type atomicBool int32

func (b *atomicBool) isSet() bool { return atomic.LoadInt32((*int32)(b)) != 0 }

type QWorker struct {
	GitConf     config.Gitlab
	ConnContext func(ctx context.Context, c net.Conn) context.Context
	db          *sqlx.DB
	inShutdown  atomicBool
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
		//запрос к бд
		//старт обработки

		dw, err := selectDataFromWork(qw.db)
		if err != nil {
			//отправлять боту
			continue
		}
		data, err := base64.StdEncoding.DecodeString(dw.Base64data)
		if err != nil {
			//фикс ошибки
		}

		f, _ := os.Create("data.epf")
		f.Write(data)
		
		createCommitDataProc(dw)

		time.Sleep(time.Minute * time.Duration(sleepMinute))
	}
	return nil

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
