package generate

import (
	"checkin/config"
	"checkin/query/model"

	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Generate() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	dialector := config.EnvConfig.GetGormDialector()

	gormdb, _ := gorm.Open(dialector, &gorm.Config{

		NamingStrategy: schema.NamingStrategy{
			TablePrefix: config.EnvConfig.DB_PREFIX,
		},
	})

	g.UseDB(gormdb) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(model.CheckinDevice{}, model.CheckinDeviceUser{}, model.CheckinDeviceRecord{})
	// Generate Type Safe API with Dynamic SQL defined on Querier interface
	// g.ApplyInterface(func(Querier) {}, model.UserToken{}, model.Video{}, model.User{}, model.Comment{})

	// Generate the code
	g.Execute()
}
