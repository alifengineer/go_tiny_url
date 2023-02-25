package postgres

import (
	"context"
	pb "go_auth_api_gateway/genproto/auth_service"
	"go_auth_api_gateway/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type scopeRepo struct {
	db *pgxpool.Pool
}

func NewScopeRepo(db *pgxpool.Pool) storage.ScopeRepoI {
	return &scopeRepo{
		db: db,
	}
}

func (r *scopeRepo) Upsert(ctx context.Context, entity *pb.UpsertScopeRequest) (pKey *pb.ScopePrimaryKey, err error) {
	query := `INSERT INTO "scope" (
		client_platform_id,
		path,
		method,
		requests
	) values (
		$1,
		$2,
		$3,
		$4
	) ON CONFLICT (
		client_platform_id,
		path,
		method
	) DO UPDATE SET requests = "scope".requests + $4, updated_at = NOW()`

	_, err = r.db.Exec(ctx, query,
		entity.ClientPlatformId,
		entity.Path,
		entity.Method,
		1,
	)

	pKey = &pb.ScopePrimaryKey{
		ClientPlatformId: entity.ClientPlatformId,
		Path:             entity.Path,
		Method:           entity.Method,
	}

	return pKey, err
}

func (r *scopeRepo) GetByPK(ctx context.Context, pKey *pb.ScopePrimaryKey) (res *pb.Scope, err error) {
	res = &pb.Scope{}
	query := `SELECT
		client_platform_id,
		path,
		method,
		requests
	FROM
		"scope"
	WHERE
		client_platform_id = $1 AND path = $2 AND method = $3`

	err = r.db.QueryRow(ctx, query, pKey.ClientPlatformId, pKey.Path, pKey.Method).Scan(
		&res.ClientPlatformId,
		&res.Path,
		&res.Method,
		&res.Requests,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *scopeRepo) GetList(ctx context.Context, req *pb.GetScopeListRequest) (res *pb.GetScopesResponse, err error) {
	res = &pb.GetScopesResponse{}
	query := `SELECT
			client_platform_id,
			COALESCE(path, ''),
			COALESCE(method, ''),
			COALESCE(requests, 0) AS requests
  		FROM
			"scope"
  		WHERE client_platform_id = $1`

	if req.OrderBy == "" && req.OrderType == "" {
		query += " ORDER BY " + req.OrderBy + " " + req.OrderType
	}

	rows, err := r.db.Query(ctx, query, req.ClientPlatformId)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		scope := &pb.Scope{}
		err = rows.Scan(
			&scope.ClientPlatformId,
			&scope.Path,
			&scope.Method,
			&scope.Requests,
		)
		if err != nil {
			return res, err
		}

		res.Scopes = append(res.Scopes, scope)
	}

	return res, nil
}
