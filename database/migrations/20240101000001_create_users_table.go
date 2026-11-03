package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"elcomware/app/facades"
)

type M20240101000001CreateUsersTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000001CreateUsersTable) Signature() string {
	return "20240101000001_create_users_table"
}

// Up Run the migrations.
func (r *M20240101000001CreateUsersTable) Up() error {
	if !facades.Schema().HasTable("users") {
		return facades.Schema().Create("users", func(table schema.Blueprint) {
			table.String("id", 64)
			table.Primary("id")
			table.String("name")
			table.String("email").Nullable()
			table.String("password").Nullable()
			table.String("role", 32).Default("buyer")
			table.String("company_name").Nullable()
			table.String("country", 64).Nullable()
			table.String("account_type", 32).Default("individual")
			table.Boolean("verified").Default(false)
			table.Text("avatar").Nullable()
			table.Timestamps()

			table.Index("email")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20240101000001CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
