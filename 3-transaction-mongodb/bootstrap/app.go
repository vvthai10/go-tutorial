package bootstrap

import (
	"github.com/vvthai10/transaction-mongodb/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	Env *config.Env
	DB *mongo.Database
}

func App() Application{
	app := &Application{}
	app.Env = config.NewEnv()
	app.DB = NewMongoDB(app.Env)

	return *app
}