package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/intervention-engine/fhir/models"
	"gopkg.in/mgo.v2/bson"
)

func OperationOutcomeIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.OperationOutcome
	c := Database.C("operationoutcomes")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var operationoutcomeEntryList []models.OperationOutcomeBundleEntry
	for _, operationoutcome := range result {
		var entry models.OperationOutcomeBundleEntry
		entry.Title = "OperationOutcome " + operationoutcome.Id
		entry.Id = operationoutcome.Id
		entry.Content = operationoutcome
		operationoutcomeEntryList = append(operationoutcomeEntryList, entry)
	}

	var bundle models.OperationOutcomeBundle
	bundle.Type = "Bundle"
	bundle.Title = "OperationOutcome Index"
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Updated = time.Now()
	bundle.TotalResults = len(result)
	bundle.Entry = operationoutcomeEntryList

	log.Println("Setting operationoutcome search context")
	context.Set(r, "OperationOutcome", result)
	context.Set(r, "Resource", "OperationOutcome")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadOperationOutcome(r *http.Request) (*models.OperationOutcome, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("operationoutcomes")
	result := models.OperationOutcome{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting operationoutcome read context")
	context.Set(r, "OperationOutcome", result)
	context.Set(r, "Resource", "OperationOutcome")
	return &result, nil
}

func OperationOutcomeShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadOperationOutcome(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "OperationOutcome"))
}

func OperationOutcomeCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	operationoutcome := &models.OperationOutcome{}
	err := decoder.Decode(operationoutcome)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("operationoutcomes")
	i := bson.NewObjectId()
	operationoutcome.Id = i.Hex()
	err = c.Insert(operationoutcome)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting operationoutcome create context")
	context.Set(r, "OperationOutcome", operationoutcome)
	context.Set(r, "Resource", "OperationOutcome")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/OperationOutcome/"+i.Hex())
}

func OperationOutcomeUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	operationoutcome := &models.OperationOutcome{}
	err := decoder.Decode(operationoutcome)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("operationoutcomes")
	operationoutcome.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, operationoutcome)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting operationoutcome update context")
	context.Set(r, "OperationOutcome", operationoutcome)
	context.Set(r, "Resource", "OperationOutcome")
	context.Set(r, "Action", "update")
}

func OperationOutcomeDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("operationoutcomes")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting operationoutcome delete context")
	context.Set(r, "OperationOutcome", id.Hex())
	context.Set(r, "Resource", "OperationOutcome")
	context.Set(r, "Action", "delete")
}
