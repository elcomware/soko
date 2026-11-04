package migrations

import (
	"elcomware/app/facades"
	"github.com/goravel/framework/contracts/database/schema"
)

type M20261104000003CreateAccessSettings struct{}

func (*M20261104000003CreateAccessSettings) Signature() string {
	return "20261104000003_create_access_settings"
}
func (*M20261104000003CreateAccessSettings) Up() error {
	for _, name := range []string{"user_access_profiles", "access_audits", "brand_settings", "merchant_payment_settings", "access_mutex"} {
		if facades.Schema().HasTable(name) {
			continue
		}
		if err := facades.Schema().Create(name, func(t schema.Blueprint) {
			switch name {
			case "access_mutex":
				t.String("id", 32)
				t.Primary("id")
				t.Integer("version").Default(0)
			case "user_access_profiles":
				t.String("user_id", 64)
				t.Primary("user_id")
				t.LongText("grants_json")
				t.Integer("version").Default(1)
			case "access_audits":
				t.String("id", 64)
				t.Primary("id")
				t.String("actor_id", 64)
				t.String("target_id", 64)
				t.LongText("before_json")
				t.LongText("after_json")
				t.BigInteger("created_at")
			case "brand_settings":
				t.String("id", 32)
				t.Primary("id")
				t.LongText("settings_json")
				t.Integer("version").Default(1)
			case "merchant_payment_settings":
				t.String("merchant_id", 64)
				t.Primary("merchant_id")
				t.LongText("settings_json")
				t.Integer("version").Default(1)
			}
		}); err != nil {
			return err
		}
	}
	_, err := facades.Orm().Query().Exec("INSERT INTO access_mutex (id, version) SELECT ?, ? WHERE NOT EXISTS (SELECT 1 FROM access_mutex WHERE id = ?)", "permissions", 0, "permissions")
	return err
}
func (*M20261104000003CreateAccessSettings) Down() error {
	for _, n := range []string{"access_mutex", "merchant_payment_settings", "brand_settings", "access_audits", "user_access_profiles"} {
		if err := facades.Schema().DropIfExists(n); err != nil {
			return err
		}
	}
	return nil
}
