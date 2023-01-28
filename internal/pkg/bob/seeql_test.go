package bob

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string `db:"name" mysql:"type:varchar"`
}

func (User) TableName() string {
	return "users"
}

type Profile struct {
	Name string `db:"profile" mysql:"type:varchar"`
}

func (Profile) TableName() string {
	return "user_profiles"
}

type PhoneNumber struct {
	Number int `db:"number"`
}

func TestBuilder(t *testing.T) {
	t.Run("replaces :table with derived table name", func(t *testing.T) {
		q := Table(&PhoneNumber{}).Build(`SELECT * FROM :table`)
		assert.Equal(t, q, "SELECT * FROM phone_numbers")
	})

	t.Run("replaces :table with table name", func(t *testing.T) {
		q := Table(&User{}).Build(`SELECT * FROM :table`)
		assert.Equal(t, q, "SELECT * FROM users")
	})

	t.Run("replaces multiple table names", func(t *testing.T) {
		q := Table(&User{}).Build(`SELECT * FROM :user_table INNER JOIN :user_table_profile`, "user_table")
		q = Table(&Profile{}).Build(q, "user_table_profile")

		assert.Equal(t, q, "SELECT * FROM users INNER JOIN user_profiles")
	})

	t.Run("does not replace partial matches", func(t *testing.T) {
		q := Table(&User{}).Build(`SELECT * FROM :table INNER JOIN :profile_table`, "table")
		q = Table(&Profile{}).Build(q, "profile_table")

		assert.Equal(t, q, "SELECT * FROM users INNER JOIN user_profiles")
	})

	t.Run("replaces oop query", func(t *testing.T) {
		q := Table(&User{}).Build(`SELECT * FROM :user_table WHERE :user_table.name = ?`, "user_table")

		assert.Equal(t, q, "SELECT * FROM users WHERE users.name = ?")
	})

	t.Run("builds query with clause", func(t *testing.T) {
		b := Table(&User{}).BuildWithClause(`UPDATE :table SET %s WHERE id = ?`)
		q := b.Assemble(map[string]interface{}{"name": "fooz", "description": ""})

		assert.Equal(t, q, "UPDATE users SET name=:name, description=:description WHERE id = ?")
		// assert.Equal(t, q, "UPDATE users SET description=:description, name=:name WHERE id = ?")
	})
}
