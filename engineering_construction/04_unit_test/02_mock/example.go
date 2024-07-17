package example

import "context"

type DB interface {
	GetUser(ctx context.Context, id int) (string, error)
}

func GetUserName(ctx context.Context, db DB, id int) (string, error) {
	return db.GetUser(ctx, id)
}
