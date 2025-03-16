package fixture

import (
	"context"
	"os"

	"github.com/osmanunal/commerce-oms/product-service/internal/model"
	"github.com/osmanunal/commerce-oms/product-service/pkg/config"
	"github.com/osmanunal/commerce-oms/product-service/pkg/database"
	"github.com/uptrace/bun/dbfixture"
)

func InsertDefaultData() error {
	cfg := config.Read()
	db := database.ConnectDB(cfg.DBConfig)

	db.RegisterModel(&model.Product{})

	fixture := dbfixture.New(db)
	err := fixture.Load(context.Background(), os.DirFS("fixture"), "fixture.yml")
	return err
}
