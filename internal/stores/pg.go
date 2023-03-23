package stores

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/e-faizov/GophKeeper/internal/models"
)

func NewPgStore(conn string) (*PGStore, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("error open db: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err = initTables(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("error init tables: %w", err)
	}

	return &PGStore{
		db: db,
	}, nil
}

type PGStore struct {
	db *sql.DB
}

func (p *PGStore) Register(ctx context.Context, login, password string) (bool, string, error) {
	uid := uuid.New()

	hash, err := calcHash(password)
	if err != nil {
		return false, "", err
	}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return false, "", err
	}

	rollback := func(err error) error {
		errRoll := tx.Rollback()
		if errRoll != nil {
			err = multierror.Append(err, fmt.Errorf("error on rollback %w", errRoll))
		}
		return err
	}

	sqlString := `insert into users (uuid, login, hash) values ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, sqlString, uid.String(), login, hash)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Constraint == "users_login_uindex" {
			return false, "", nil
		}
		return false, "", rollback(err)
	}
	err = tx.Commit()
	if err != nil {
		return false, "", err
	}
	return true, uid.String(), nil
}
func (p *PGStore) Login(ctx context.Context, login, password string) (bool, string, error) {
	sqlString := `select uuid, hash from users where login=$1`

	rows, err := p.db.QueryContext(ctx, sqlString, login)
	if err != nil {
		return false, "", err
	}
	defer rows.Close()
	count := 0
	resUudi := ""
	userHash := ""
	for rows.Next() {
		count++
		err = rows.Scan(&resUudi, &userHash)
		if err != nil {
			return false, "", err
		}
	}
	if err = rows.Err(); err != nil {
		return false, "", err
	}
	if count != 1 {
		return false, "", nil
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userHash), []byte(password)); err != nil {
		return false, "", err
	}

	return true, resUudi, nil
}

func (p *PGStore) NewSecret(ctx context.Context, userID, uid string, s models.Secret) error {
	sqlString := `insert into secrets (uid, user_uid, data1, data2, data3, type, version)
				values ($1, $2, $3, $4, $5, $6, 1)`
	_, err := p.db.ExecContext(ctx, sqlString, uid, userID, s.Name, s.Data, s.Meta, s.Type)
	if err != nil {
		return err
	}
	return nil
}

func (p *PGStore) RemoveSecret(ctx context.Context, userID string, s models.Secret) error {
	sqlString := `delete from secrets where user_uid=$1 and uid=$2 and version=$3`
	_, err := p.db.ExecContext(ctx, sqlString, userID, s.Uid, s.Version)
	return err
}

func (p *PGStore) GetSecret(ctx context.Context, userID string, s models.Secret) (models.Secret, error) {
	sqlString := `select data2, type, version from secrets where user_uid=$1 and uid=$2`

	row := p.db.QueryRowContext(ctx, sqlString, userID, s.Uid)

	var res models.Secret
	err := row.Scan(&res.Data, &res.Type, &res.Version)
	if err != nil {
		return models.Secret{}, err
	}
	res.Uid = s.Uid

	if res.Version != s.Version {
		return models.Secret{}, errors.New("wrong version")
	}

	return res, nil
}

func (p *PGStore) GetSecretsList(ctx context.Context, userID string) ([]models.Secret, error) {
	sqlString := `select uid, data1, data3, version, type from secrets where user_uid=$1`
	rows, err := p.db.QueryContext(ctx, sqlString, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Secret
	for rows.Next() {
		var tmp models.Secret
		err = rows.Scan(&tmp.Uid, &tmp.Name, &tmp.Meta, &tmp.Version, &tmp.Type)
		if err != nil {
			return nil, err
		}
		res = append(res, tmp)
	}

	return res, nil
}

func (p *PGStore) EditSecret(ctx context.Context, userID string, s models.Secret) error {
	sqlString := `update secrets set
					version=version+1, data2=$1, data3=$2
					where user_uid=$3 and uid=$4 and version=$5`

	_, err := p.db.ExecContext(ctx, sqlString, s.Data, s.Meta, userID, s.Uid, s.Version)
	return err
}

func calcHash(s string) (string, error) {
	bt, err := bcrypt.GenerateFromPassword([]byte(s), 8)
	if err != nil {
		return "", err
	}
	return string(bt), nil
}
