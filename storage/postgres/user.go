package postgres

import (
	"context"
	"database/sql"
	"go_auth_api_gateway/config"
	pb "go_auth_api_gateway/genproto/auth_service"
	"go_auth_api_gateway/pkg/helper"
	"go_auth_api_gateway/storage"
	"time"

	"github.com/saidamir98/udevs_pkg/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) storage.UserRepoI {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, entity *pb.CreateUserRequest) (pKey *pb.UserPrimaryKey, err error) {
	query := `INSERT INTO "users" (
		id,
		first_name,
		last_name,
		phone,
		username,
		password,
		created_at,
		updated_at
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8
	)`

	uuid, err := uuid.NewRandom()
	if err != nil {
		return pKey, err
	}

	_, err = r.db.Exec(ctx, query,
		uuid.String(),
		entity.GetFirstName(),
		entity.GetLastName(),
		entity.GetPhone(),
		entity.GetUsername(),
		entity.GetPassword(),
		time.Now().UTC().Format(time.RFC3339),
		time.Now().UTC().Format(time.RFC3339),
	)

	pKey = &pb.UserPrimaryKey{
		Id: uuid.String(),
	}

	return pKey, err
}

func (r *userRepo) GetByPK(ctx context.Context, pKey *pb.UserPrimaryKey) (res *pb.User, err error) {
	res = &pb.User{}
	query := `SELECT
		id,
		first_name,
		last_name,
		phone,
		username,
		password,
		TO_CHAR(created_at, ` + config.DatabaseQueryTimeLayout + `) AS created_at,
		TO_CHAR(updated_at, ` + config.DatabaseQueryTimeLayout + `) AS updated_at
	FROM
		"users"
	WHERE deleted_at = 0 AND
		id = $1`

	err = r.db.QueryRow(ctx, query, pKey.Id).Scan(
		&res.Id,
		&res.FirstName,
		&res.LastName,
		&res.Phone,
		&res.Username,
		&res.Password,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *userRepo) GetList(ctx context.Context, queryParam *pb.GetUserListRequest) (res *pb.GetUserListResponse, err error) {
	res = &pb.GetUserListResponse{}
	params := make(map[string]interface{})
	var arr []interface{}
	query := `SELECT
	id,
	first_name,
	last_name,
	phone,
	username,
	password,
	TO_CHAR(created_at, ` + config.DatabaseQueryTimeLayout + `) AS created_at,
	TO_CHAR(updated_at, ` + config.DatabaseQueryTimeLayout + `) AS updated_at
	FROM
		"users"`
	filter := " WHERE deleted_at = 0"
	order := " ORDER BY created_at"
	arrangement := " DESC"
	offset := " OFFSET 0"
	limit := " LIMIT 10"

	if len(queryParam.Search) > 0 {
		params["search"] = queryParam.Search
		filter += " AND ((name || phone || user_name ) ILIKE ('%' || :search || '%'))"
	}

	if queryParam.Offset > 0 {
		params["offset"] = queryParam.Offset
		offset = " OFFSET :offset"
	}

	if queryParam.Limit > 0 {
		params["limit"] = queryParam.Limit
		limit = " LIMIT :limit"
	}

	cQ := `SELECT count(1) FROM "users"` + filter
	cQ, arr = helper.ReplaceQueryParams(cQ, params)
	err = r.db.QueryRow(ctx, cQ, arr...).Scan(
		&res.Count,
	)
	if err != nil {
		return res, err
	}

	q := query + filter + order + arrangement + offset + limit

	q, arr = helper.ReplaceQueryParams(q, params)
	rows, err := r.db.Query(ctx, q, arr...)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		obj := &pb.User{}
		var (
			createdAt sql.NullString
			updatedAt sql.NullString
		)

		err = rows.Scan(
			&obj.Id,
			&obj.FirstName,
			&obj.LastName,
			&obj.Phone,
			&obj.Username,
			&obj.Password,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return res, err
		}
		if createdAt.Valid {
			obj.CreatedAt = createdAt.String
		}

		if updatedAt.Valid {
			obj.UpdatedAt = updatedAt.String
		}

		res.Users = append(res.Users, obj)
	}

	return res, nil
}

func (r *userRepo) Update(ctx context.Context, entity *pb.UpdateUserRequest) (rowsAffected int64, err error) {
	query := `UPDATE "users" SET
		first_name = :first_name,
		last_name = :last_name,
		phone = :phone,
		username = :username,
		updated_at = now()
	WHERE
		id = :id`

	params := map[string]interface{}{
		"id":         entity.Id,
		"first_name": entity.FirstName,
		"last_name":  entity.LastName,
		"phone":      entity.Phone,
		"username":   entity.Username,
	}

	q, arr := helper.ReplaceQueryParams(query, params)
	result, err := r.db.Exec(ctx, q, arr...)
	if err != nil {
		return 0, err
	}

	rowsAffected = result.RowsAffected()

	return rowsAffected, err
}

func (r *userRepo) Delete(ctx context.Context, pKey *pb.UserPrimaryKey) (rowsAffected int64, err error) {
	query := `UPDATE "users" SET deleted_at = date_part('epoch', CURRENT_TIMESTAMP)::int`

	result, err := r.db.Exec(ctx, query, pKey.Id)
	if err != nil {
		return 0, err
	}

	rowsAffected = result.RowsAffected()

	return rowsAffected, err
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (res *pb.User, err error) {
	res = &pb.User{}

	query := `SELECT
				id,
				first_name,
				last_name,
				phone,
				username,
				password,
				TO_CHAR(created_at, ` + config.DatabaseQueryTimeLayout + `) AS created_at,
				TO_CHAR(updated_at, ` + config.DatabaseQueryTimeLayout + `) AS updated_at
	FROM
		"users"
	WHERE`

	if util.IsValidEmail(username) {
		query = query + ` email = $1`
	} else if util.IsValidPhone(username) {
		query = query + ` phone = $1`
	} else {
		query = query + ` username = $1`
	}

	err = r.db.QueryRow(ctx, query, username).Scan(
		&res.Id,
		&res.FirstName,
		&res.LastName,
		&res.Phone,
		&res.Username,
		&res.Password,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *userRepo) ResetPassword(ctx context.Context, user *pb.ResetPasswordRequest) (rowsAffected int64, err error) {
	query := `UPDATE "users" SET
		password = :password,
		updated_at = now()
	WHERE
		id = :id`

	params := map[string]interface{}{
		"id":       user.UserId,
		"password": user.Password,
	}

	q, arr := helper.ReplaceQueryParams(query, params)
	result, err := r.db.Exec(ctx, q, arr...)
	if err != nil {
		return 0, err
	}

	rowsAffected = result.RowsAffected()

	return rowsAffected, err
}
