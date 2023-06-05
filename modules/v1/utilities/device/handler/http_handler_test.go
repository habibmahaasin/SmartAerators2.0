package handler

import (
	"GuppyTech/modules/v1/utilities/device/repository"
	"GuppyTech/modules/v1/utilities/device/service"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	dsn := "host=satao.db.elephantsql.com user=dwbejsql password=Kb48I9w7spTcFsiPCP2tPHeR9mhm3Ds1 dbname=dwbejsql port=5432 TimeZone=Asia/Jakarta"
	Db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err.Error())
	}
	return Db
}

func SetupHandler() *deviceHandler {
	Repository := repository.NewRepository(SetupDB())
	Service := service.NewService(Repository)
	Handler := NewDeviceHandler(Service)
	return Handler
}

func SetUpRouter() *gin.Engine {
	app := gin.Default()
	return app
}

func Test_SubscribeWebhook(t *testing.T) {
	tests := []struct {
		name       string
		inputJSON  string
		statusCode int
		response   string
	}{
		{
			name: "Success Input Data",
			inputJSON: `{
				"m2m:sgn" : {
				   "m2m:nev" : {
					  "m2m:rep" : {
						 "m2m:cin" : {
							"rn" : "cin_b7NDksZhTsWEmLg0",
							"ty" : 4,
							"ri" : "/antares-cse/cin-b7NDksZhTsWEmLg0",
							"pi" : "/antares-cse/cnt-ps9t5UiX15TVLxYB",
							"ct" : "20220405T160104",
							"lt" : "20220405T160104",
							"st" : 0,
							"cnf" : "text/plain:0",
							"cs" : 266,
							"con" : "{\"aeratorMode\":0,\"temperature\":28.375,\"ph\":21.97006,\"dissolvedOxygen\":50.34506,\"statusDevice\":0}"
						 }
					  },
					  "m2m:rss" : 1
				   },
				   "m2m:sud" : false,
				   "m2m:sur" : "/antares-cse/sub-HQ5_ZGw6Ts6SUnrz"
				}
			 }`,
			statusCode: 200,
			response:   ``,
		},
	}
	r := SetUpRouter()
	r.POST("/api/v1/webhook", SetupHandler().SubscribeWebhook)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest("POST", "/api/v1/webhook", bytes.NewBufferString(test.inputJSON))
			if err != nil {
				t.Errorf(err.Error())
			}

			response := httptest.NewRecorder()
			r.ServeHTTP(response, request)
			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf(err.Error())
			}

			assert.Equal(t, response.Code, test.statusCode)
			assert.Equal(t, string(responseData), test.response)
		})
	}
}