package auth_repository

import (
	"github.com/jmoiron/sqlx"
	"strings"
	"tok-core/src/data/entities"
	"tok-core/src/data/models"
)

func (repo *Repository) Signup(model *models.Signup, encryptedPassword string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	// entity
	entity := &entities.Signup{
		Username: strings.ToLower(strings.TrimSpace(model.Username)),
		Password: strings.ToLower(strings.TrimSpace(encryptedPassword)),
	}

	// query
	signupQuery := repo.Script("auth", "signup")

	if err := repo.Transaction(repo.connection, func(tx *sqlx.Tx) error {
		if _, err := tx.NamedExecContext(ctx, signupQuery, entity); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Login(username, encryptedPassword string) (*entities.UserGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	// entity
	entity := &entities.Login{
		Username: strings.ToLower(strings.TrimSpace(username)),
		Password: strings.ToLower(strings.TrimSpace(encryptedPassword)),
	}

	// query
	query := repo.Script("auth", "login")

	rows, err := repo.connection.NamedQueryContext(ctx, query, entity)
	if err != nil {
		return nil, err
	}
	defer repo.CloseRows(rows)

	if rows.Next() {
		user := entities.UserGet{}
		if err = rows.StructScan(&user); err != nil {
			return nil, err
		}
		return &user, nil
	}

	return nil, nil
}
