package storage

import (
	"errors"
	"fmt"
	"strconv"
	"webcommitter/internal/config"
	"webcommitter/internal/model"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommitStatus struct {
	CommitCount string
}

type Storage struct {
	DB      *gorm.DB
	GitConf *config.Gitlab
}

func NewStorage(db *gorm.DB, git *config.Gitlab) *Storage {
	return &Storage{
		DB:      db,
		GitConf: git,
	}
}

func (s *Storage) AddNewRequest(w http.ResponseWriter, r *http.Request) (string, int, error) {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return "", http.StatusBadRequest, err
	}

	var cc model.Commit

	err = json.Unmarshal(b, &cc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return "", http.StatusBadRequest, err
	}

	var author model.Authorcommit
	var proc model.DataProccessor

	result := s.DB.FirstOrCreate(&author, cc.AuthorCommit)
	cc.AuthorCommitID = int(author.ID)

	result = s.DB.FirstOrCreate(&proc, cc.Proccessor)
	cc.ProccessorID = int(proc.ID)

	cc.Data = proc.Base64data
	result = s.DB.Omit(clause.Associations).Create(&cc)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		s.DB.Create(author)
	}

	if err != nil {
		return "", http.StatusFailedDependency, err
	} else {
		return "OK, ready", http.StatusOK, nil
	}
}

func (s *Storage) CheckedStatusQueues(w http.ResponseWriter, r *http.Request) string {
	var cm []model.Commit

	result := s.DB.Find(&cm, "its_done = ?", "false")
	return strconv.FormatInt(result.RowsAffected, 10)

}

func (s *Storage) regMessageDB(upl *model.DataCommit) error {

	// us := model.User{}
	// err := s.DB.GetContext(context.TODO(),
	// 	&us,
	// 	"Select u.id, u.extid, u.name from users as u where u.extid=$1 LIMIT 1",
	// 	upl.User.ExtID)

	// if err != nil && !strings.Contains(err.Error(),
	// 	"sql: no rows in result set") {
	// 	return err
	// }

	// if us.ID == 0 {
	// 	q := "INSERT INTO users (extId, name) VALUES ($1, $2)  RETURNING id"
	// 	err = s.DB.GetContext(context.Background(), &us, q,
	// 		string(upl.User.ExtID),
	// 		string(upl.User.FullName))
	// 	if err != nil {
	// 		return err
	// 	}

	// }
	// tx := s.DB.MustBegin()
	// tx.MustExec("INSERT INTO commit_tasks (extId,name,base64data,type,userid,textcommit) VALUES ($1, $2, $3, $4, $5, $6)",
	// 	upl.DataProccessor.ExtID,
	// 	upl.DataProccessor.Name,
	// 	upl.DataProccessor.Base64data,
	// 	upl.DataProccessor.Type,
	// 	strconv.FormatInt(us.ID, 10),
	// 	upl.TextCommit)

	// err = tx.Commit()

	// if err != nil {
	// 	return err
	// }

	return nil
}

func (s *Storage) FindUserByID(id string) (*model.User, error) {

	u := &model.User{}
	q := "Select id,extID,FullName,GitEmail,IsAdmin,TGGit FROM user WHERE id = $1"
	err := s.runSelectQuery(q, &u, u.ID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Storage) FindAdmin() (*model.User, error) {

	//q := `Select u.id,u.extid,u.name,u.is_admin,u.gitlogin, u.tgid from users as u where u.is_admin=true LIMIT 1`
	us := model.User{}
	// err := s.DB.GetContext(context.Background(), &us, q)

	// if err != nil {
	// 	return nil, err
	// }
	return &us, nil
}

func (s *Storage) runSelectQuery(query string, nO interface{}, args ...interface{}) error {

	// err := s.DB.QueryRowxContext(context.TODO(),
	// 	query, args[0:]...).Scan()

	// if err != nil {
	// 	return err
	// }

	return nil
}

func (s *Storage) FindLastCommit() (string, error) {

	var cm model.Commit

	s.DB.First(&cm, "its_done = ?", "false")
	res, err := json.Marshal(cm)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", res), nil
}
