package internal

import (
	user "movie/system/internal/user"
)

func Models() []interface{} {
	return []interface{}{
		&user.User{},
	}
}
