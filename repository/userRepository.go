package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"lunchorder/queries"
)

type UserRepository struct {
	db *sqlx.DB
}

var userRepo *UserRepository

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(user *User) error {
	result, err := r.db.Exec(queries.CreateUser, user.Name)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = uint(id)
	return nil
}

func (r *UserRepository) GetUserByName(name string) (*User, error) {
	var user User
	err := r.db.Get(&user, queries.GetUserByName, name)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByGoogleID(googleID string) (*User, error) {
	var user User
	err := r.db.Get(&user, queries.GetUserByGoogleID, googleID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id uint) (*User, error) {
	var user User
	err := r.db.Get(&user, queries.GetUserByID, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.Get(&user, queries.GetUserByEmail, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpsertUser(user *User) error {
	// 1. Check if user exists by Google ID
	existingUser, err := r.GetUserByGoogleID(*user.GoogleID)
	if err == nil {
		// User exists, update
		user.ID = existingUser.ID

		// Check for name conflict if name changed
		if user.Name != existingUser.Name {
			// Check if new name is taken
			otherUser, err := r.GetUserByName(user.Name)
			if err == nil && otherUser.ID != user.ID {
				// Name taken by someone else. Append suffix.
				originalName := user.Name
				for i := 1; i <= 10; i++ {
					newName := fmt.Sprintf("%s%d", originalName, i)
					u, err := r.GetUserByName(newName)
					if err != nil {
						// Name is free (assuming err is not found)
						user.Name = newName
						break
					}
					if u.ID == user.ID {
						// We already own this name
						user.Name = newName
						break
					}
				}
			}
		}

		_, err = r.db.NamedExec(queries.UpdateUserGoogle, user)
		return err
	}

	// 2. Check if user exists by email (might be a legacy user without Google ID)
	if user.Email != nil && *user.Email != "" {
		existingUser, err := r.GetUserByEmail(*user.Email)
		if err == nil {
			// User exists by email, update with Google ID
			user.ID = existingUser.ID

			// Check for name conflict if name changed
			if user.Name != existingUser.Name {
				// Check if new name is taken
				otherUser, err := r.GetUserByName(user.Name)
				if err == nil && otherUser.ID != user.ID {
					// Name taken by someone else. Append suffix.
					originalName := user.Name
					for i := 1; i <= 10; i++ {
						newName := fmt.Sprintf("%s%d", originalName, i)
						u, err := r.GetUserByName(newName)
						if err != nil {
							// Name is free (assuming err is not found)
							user.Name = newName
							break
						}
						if u.ID == user.ID {
							// We already own this name
							user.Name = newName
							break
						}
					}
				}
			}

			_, err = r.db.NamedExec(queries.UpdateUserGoogle, user)
			return err
		}
	}

	// 3. User does not exist, insert
	originalName := user.Name
	for i := 0; i <= 10; i++ {
		if i > 0 {
			user.Name = fmt.Sprintf("%s%d", originalName, i)
		}

		existing, err := r.GetUserByName(user.Name)
		if err == nil && existing != nil {
			// Name taken, try next
			continue
		}

		// Name free (or DB error), try insert
		_, err = r.db.NamedExec(queries.InsertUserGoogle, user)
		if err != nil {
			// If it's a duplicate entry error (race condition), we might want to retry loop?
			// But for now, return error if insert fails
			return err
		}

		// Success
		updatedUser, err := r.GetUserByGoogleID(*user.GoogleID)
		if err != nil {
			return err
		}
		user.ID = updatedUser.ID
		return nil
	}

	return fmt.Errorf("failed to find unique name for user")
}
