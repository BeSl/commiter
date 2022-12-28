package storage

import (
	"webcommitter/internal/config"
	"webcommitter/internal/model"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"gorm.io/gorm"
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

	var cc model.NewCommit

	err = json.Unmarshal(b, &cc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return "", http.StatusBadRequest, err
	}

	// err = s.regMessageDB(&upl)
	if err != nil {
		return "", http.StatusFailedDependency, err
	} else {
		return "OK, ready", http.StatusOK, nil
	}
}

func (s *Storage) CheckedStatusQueues(w http.ResponseWriter, r *http.Request) (string, error) {
	st := CommitStatus{}
	var err error
	//s.DB.Get(&st, "Select COUNT(*)as CommitCount FROM commit_tasks WHERE processed=false")

	if err != nil {
		return "", err
	}

	return st.CommitCount, nil
}

func (s *Storage) CreateTablesDB() {

	// q := `CREATE TABLE if not exists public.users (
	// 	id bigserial NOT NULL,
	// 	extid uuid NULL,
	// 	"name" varchar(150) NULL,
	// 	is_admin bool NULL,
	// 	gitlogin varchar(100) NOT NULL DEFAULT "",
	// 	tgid int4 NULL,
	// 	first_name varchar(75) NULL,
	// 	last_name varchar(75) NULL,
	// 	CONSTRAINT users_pkey PRIMARY KEY (id)
	// );`
	// s.DB.MustExec(q)

	// q = `CREATE TABLE if not exists public.commit_tasks (
	// 	id int8 NOT NULL GENERATED ALWAYS AS IDENTITY,
	// 	"name" varchar(200) NULL,
	// 	extid uuid NULL,
	// 	"type" varchar NULL,
	// 	base64data text NULL,
	// 	textcommit text NULL,
	// 	dataevent date NULL,
	// 	userid int8 NULL,
	// 	processed bool NOT NULL DEFAULT false,
	// 	error text NULL,
	// 	CONSTRAINT commit_tasks_pk PRIMARY KEY (id)
	// );`
	// s.DB.MustExec(q)
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

func (s *Storage) FindLastCommit() (*model.DataWork, error) {

	// q :=
	// 	`SELECT
	// 		coalesce(u.gitlogin, '') as gitlogin,
	// 		coalesce(u.name, '') as username,
	// 		ct.name as name,
	// 		ct.id as id,
	// 		ct.type as type,
	// 		ct.base64data as base64data,
	// 		coalesce(ct.textcommit, 'not text') as commit
	// 	FROM
	// 		commit_tasks ct
	// 			left join users u
	// 			on u.id= ct.userid
	// 	WHERE
	// 	 		ct.processed =false
	// 	ORDER BY ct.id
	// 	LIMIT 1`

	dw := model.DataWork{}
	// err := s.DB.GetContext(context.Background(), &dw, q)

	// if err != nil {
	// 	return nil, err
	// }
	// if len(dw.GitLogin) == 0 {
	// 	return nil, errors.New("Не заполнен пользователь " + dw.UserName)
	// }

	// if err == sql.ErrNoRows {
	// 	return nil, nil
	// }

	// if err != nil {
	// 	return nil, errors.New("Ошибка в FindLastCommit " + err.Error())
	// }

	return &dw, nil
}
