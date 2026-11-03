package seeders

type DatabaseSeeder struct{}

// Signature The unique signature of the seeder.
func (s *DatabaseSeeder) Signature() string {
	return "DatabaseSeeder"
}

// Run executes the seeder logic.
func (s *DatabaseSeeder) Run() error {
	userSeeder := &UserSeeder{}
	if err := userSeeder.Run(); err != nil {
		return err
	}

	productSeeder := &ProductSeeder{}
	if err := productSeeder.Run(); err != nil {
		return err
	}

	rfqSeeder := &RfqSeeder{}
	if err := rfqSeeder.Run(); err != nil {
		return err
	}

	offerSeeder := &OfferSeeder{}
	if err := offerSeeder.Run(); err != nil {
		return err
	}

	return nil
}

