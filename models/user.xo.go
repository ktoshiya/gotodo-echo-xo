package models

// Code generated by xo. DO NOT EDIT.

import (
	"context"
	"time"
)

// User represents a row from 'todo.users'.
type User struct {
	ID       uint64    `json:"id"`       // id
	Name     string    `json:"name"`     // name
	Password string    `json:"password"` // password
	Created  time.Time `json:"created"`  // created
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the User exists in the database.
func (u *User) Exists() bool {
	return u._exists
}

// Deleted returns true when the User has been marked for deletion from
// the database.
func (u *User) Deleted() bool {
	return u._deleted
}

// Insert inserts the User to the database.
func (u *User) Insert(ctx context.Context, db DB) error {
	switch {
	case u._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case u._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO users (` +
		`name, password, created` +
		`) VALUES (` +
		`?, ?, ?` +
		`)`
	// run
	logf(sqlstr, u.Name, u.Password, u.Created)
	res, err := db.ExecContext(ctx, sqlstr, u.Name, u.Password, u.Created)
	if err != nil {
		return logerror(err)
	}
	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return logerror(err)
	} // set primary key
	u.ID = uint64(id)
	// set exists
	u._exists = true
	return nil
}

// Update updates a User in the database.
func (u *User) Update(ctx context.Context, db DB) error {
	switch {
	case !u._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case u._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE todo.users SET ` +
		`name = ?, password = ?, created = ? ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, u.Name, u.Password, u.Created, u.ID)
	if _, err := db.ExecContext(ctx, sqlstr, u.Name, u.Password, u.Created, u.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the User to the database.
func (u *User) Save(ctx context.Context, db DB) error {
	if u.Exists() {
		return u.Update(ctx, db)
	}
	return u.Insert(ctx, db)
}

// Upsert performs an upsert for User.
func (u *User) Upsert(ctx context.Context, db DB) error {
	switch {
	case u._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO todo.users (` +
		`id, name, password, created` +
		`) VALUES (` +
		`?, ?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`name = VALUES(name), password = VALUES(password), created = VALUES(created)`
	// run
	logf(sqlstr, u.ID, u.Name, u.Password, u.Created)
	if _, err := db.ExecContext(ctx, sqlstr, u.ID, u.Name, u.Password, u.Created); err != nil {
		return logerror(err)
	}
	// set exists
	u._exists = true
	return nil
}

// Delete deletes the User from the database.
func (u *User) Delete(ctx context.Context, db DB) error {
	switch {
	case !u._exists: // doesn't exist
		return nil
	case u._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM todo.users ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, u.ID)
	if _, err := db.ExecContext(ctx, sqlstr, u.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	u._deleted = true
	return nil
}

// UserByName retrieves a row from 'todo.users' as a User.
//
// Generated from index 'uix_name'.
func UserByName(ctx context.Context, db DB, name string) (*User, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, name, password, created ` +
		`FROM todo.users ` +
		`WHERE name = ?`
	// run
	logf(sqlstr, name)
	u := User{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, name).Scan(&u.ID, &u.Name, &u.Password, &u.Created); err != nil {
		return nil, logerror(err)
	}
	return &u, nil
}

// UserByID retrieves a row from 'todo.users' as a User.
//
// Generated from index 'users_id_pkey'.
func UserByID(ctx context.Context, db DB, id uint64) (*User, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, name, password, created ` +
		`FROM todo.users ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, id)
	u := User{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, id).Scan(&u.ID, &u.Name, &u.Password, &u.Created); err != nil {
		return nil, logerror(err)
	}
	return &u, nil
}
