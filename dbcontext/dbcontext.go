package dbcontext

import (
	"allopopot-interconnect-service/config"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var UserModel *mongo.Collection
var TokenModel *mongo.Collection
var ExpencierProjectsModel *mongo.Collection
var ExpencierTransactionsModel *mongo.Collection

func InitDb() {
	bsonOpts := &options.BSONOptions{
		UseJSONStructTags: true,
		NilMapAsEmpty:     true,
		NilSliceAsEmpty:   true,
	}
	log.Println("Database Connection Initializing")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MONGODB_URI).SetBSONOptions(bsonOpts))
	if err != nil {
		log.Panicln("Database Connection Failed")
	}
	DB = client.Database(config.MONGODB_DATABASE_NAME)
	log.Println("Database Connection Successful")
	SetCollections()
	CreateIndexes()
}

func SetCollections() {
	UserModel = DB.Collection("users")
	TokenModel = DB.Collection("tokens")
	ExpencierProjectsModel = DB.Collection("expencier_projects")
	ExpencierTransactionsModel = DB.Collection("expencier_transactions")
}

func CreateIndexes() {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := UserModel.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Database Indexes Created Successfully")
}
