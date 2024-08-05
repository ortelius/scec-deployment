// Ortelius v11 Deployment Microservice that handles creating and retrieving Deployments
package main

import (
	"context"
	"encoding/json"

	_ "github.com/ortelius/scec-deployment/docs"

	driver "github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/arangodb/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/ortelius/scec-commons/database"
	"github.com/ortelius/scec-commons/model"
)

var logger = database.InitLogger()
var dbconn = database.InitializeDatabase()

// GetDeployments godoc
// @Summary Get a List of Deployments
// @Description Get a list of deploymentss.
// @Tags deployments
// @Accept */*
// @Produce json
// @Success 200
// @Router /msapi/deployment [get]
func GetDeployments(c *fiber.Ctx) error {

	var cursor driver.Cursor       // db cursor for rows
	var err error                  // for error handling
	var ctx = context.Background() // use default database context

	// query all the deployments in the collection
	aql := `FOR deployment in evidence
			FILTER (deployment.objtype == 'Deployment')
			RETURN deployment`

	// execute the query with no parameters
	if cursor, err = dbconn.Database.Query(ctx, aql, nil); err != nil {
		logger.Sugar().Errorf("Failed to run query: %v", err) // log error
	}

	defer cursor.Close() // close the cursor when returning from this function

	var deployments []*model.Deployment // define a list of deployments to be returned

	for cursor.HasMore() { // loop thru all of the documents

		deployment := model.NewDeployment() // fetched deployment
		var meta driver.DocumentMeta        // data about the fetch

		// fetch a document from the cursor
		if meta, err = cursor.ReadDocument(ctx, deployment); err != nil {
			logger.Sugar().Errorf("Failed to read document: %v", err)
		}
		deployments = append(deployments, deployment)                        // add the deployment to the list
		logger.Sugar().Infof("Got doc with key '%s' from query\n", meta.Key) // log the key
	}

	return c.JSON(deployments) // return the list of deployments in JSON format
}

// GetDeployment godoc
// @Summary Get a Deployment
// @Description Get a deployment based on the _key or name.
// @Tags deployment
// @Accept */*
// @Produce json
// @Success 200
// @Router /msapi/deployment/:key [get]
func GetDeployment(c *fiber.Ctx) error {

	var cursor driver.Cursor       // db cursor for rows
	var err error                  // for error handling
	var ctx = context.Background() // use default database context

	key := c.Params("key")                // key from URL
	parameters := map[string]interface{}{ // parameters
		"key": key,
	}

	// query the deployments that match the key or name
	aql := `FOR deployment in evidence
			FILTER (deployment.name == @key or deployment._key == @key)
			RETURN deployment`

	// run the query with patameters
	if cursor, err = dbconn.Database.Query(ctx, aql, &driver.QueryOptions{BindVars: parameters}); err != nil {
		logger.Sugar().Errorf("Failed to run query: %v", err)
	}

	defer cursor.Close() // close the cursor when returning from this function

	deployment := model.NewDeployment() // define a deployment to be returned

	if cursor.HasMore() { // deployment found
		var meta driver.DocumentMeta // data about the fetch

		if meta, err = cursor.ReadDocument(ctx, deployment); err != nil { // fetch the document into the object
			logger.Sugar().Errorf("Failed to read document: %v", err)
		}
		logger.Sugar().Infof("Got doc with key '%s' from query\n", meta.Key)

	} else { // not found so get from NFT Storage
		if jsonStr, exists := database.MakeJSON(key); exists {
			if err := json.Unmarshal([]byte(jsonStr), deployment); err != nil { // convert the JSON string from LTF into the object
				logger.Sugar().Errorf("Failed to unmarshal from LTS: %v", err)
			}
		}
	}

	return c.JSON(deployment) // return the deployment in JSON format
}

// NewDeployment godoc
// @Summary Create a Deployment
// @Description Create a new Deployment and persist it
// @Tags deployment
// @Accept application/json
// @Produce json
// @Success 200
// @Router /msapi/deployment [post]
func NewDeployment(c *fiber.Ctx) error {

	var err error                       // for error handling
	var meta driver.DocumentMeta        // data about the document
	var ctx = context.Background()      // use default database context
	deployment := new(model.Deployment) // define a deployment to be returned

	if err = c.BodyParser(deployment); err != nil { // parse the JSON into the deployment object
		return c.Status(503).Send([]byte(err.Error()))
	}

	cid, dbStr := database.MakeNFT(deployment) // normalize the object into NFTs and JSON string for db persistence

	logger.Sugar().Infof("%s=%s\n", cid, dbStr) // log the new nft

	var resp driver.CollectionDocumentCreateResponse
	// add the deployment to the database.  Ignore if it already exists since it will be identical
	if resp, err = dbconn.Collection.CreateDocument(ctx, deployment); err != nil && !shared.IsConflict(err) {
		logger.Sugar().Errorf("Failed to create document: %v", err)
	}
	meta = resp.DocumentMeta
	logger.Sugar().Infof("Created document in collection '%s' in db '%s' key='%s'\n", dbconn.Collection.Name(), dbconn.Database.Name(), meta.Key)

	return c.JSON(deployment) // return the deployment object in JSON format.  This includes the new _key
}

// setupRoutes defines maps the routes to the functions
func setupRoutes(app *fiber.App) {

	app.Get("/swagger/*", swagger.HandlerDefault)    // handle displaying the swagger
	app.Get("/msapi/deployment", GetDeployments)     // list of deployments
	app.Get("/msapi/deployment/:key", GetDeployment) // single deployment based on name or key
	app.Post("/msapi/deployment", NewDeployment)     // save a single deployment
}

// @title Ortelius v11 deployment Microservice
// @version 11.0.0
// @description RestAPI for the Deployment Object
// @description ![Release](https://img.shields.io/github/v/release/ortelius/scec-deployment?sort=semver)
// @description ![license](https://img.shields.io/github/license/ortelius/.github)
// @description
// @description ![Build](https://img.shields.io/github/actions/workflow/status/ortelius/scec-deployment/build-push-chart.yml)
// @description [![MegaLinter](https://github.com/ortelius/scec-deployment/workflows/MegaLinter/badge.svg?branch=main)](https://github.com/ortelius/scec-deployment/actions?query=workflow%3AMegaLinter+branch%3Amain)
// @description ![CodeQL](https://github.com/ortelius/scec-deployment/workflows/CodeQL/badge.svg)
// @description [![OpenSSF-Scorecard](https://api.securityscorecards.dev/projects/github.com/ortelius/scec-deployment/badge)](https://api.securityscorecards.dev/projects/github.com/ortelius/scec-deployment)
// @description
// @description ![Discord](https://img.shields.io/discord/722468819091849316)

// @termsOfService http://swagger.io/terms/
// @contact.name Ortelius Google Group
// @contact.email ortelius-dev@googlegroups.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /msapi/deployment
func main() {
	port := ":" + database.GetEnvDefault("MS_PORT", "8080") // database port
	app := fiber.New()                                      // create a new fiber application
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowOrigins: "*",
	}))

	setupRoutes(app) // define the routes for this microservice

	if err := app.Listen(port); err != nil { // start listening for incoming connections
		logger.Sugar().Fatalf("Failed get the microservice running: %v", err)
	}
}
