package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"elcomware/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateJobsTable{},
		&migrations.M20240101000001CreateUsersTable{},
		&migrations.M20240101000002CreateProductsTable{},
		&migrations.M20240101000003CreateOrdersTable{},
		&migrations.M20240101000004CreateRfqsTable{},
		&migrations.M20240101000005CreateOffersTable{},
	}
}
