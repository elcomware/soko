package bootstrap

import (
	contractsseeder "github.com/goravel/framework/contracts/database/seeder"

	"elcomware/database/seeders"
)

func Seeders() []contractsseeder.Seeder {
	return []contractsseeder.Seeder{
		&seeders.DatabaseSeeder{},
		&seeders.UserSeeder{},
		&seeders.ProductSeeder{},
		&seeders.RfqSeeder{},
		&seeders.OfferSeeder{},
	}
}

