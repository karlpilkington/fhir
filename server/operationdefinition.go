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

func OperationDefinitionIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.OperationDefinition
	c := Database.C("operationdefinitions")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var operationdefinitionEntryList []models.OperationDefinitionBundleEntry
	for _, operationdefinition := range result {
		var entry models.OperationDefinitionBundleEntry
		entry.Title = "OperationDefinition " + operationdefinition.Id
		entry.Id = operationdefinition.Id
		entry.Content = operationdefinition
		operationdefinitionEntryList = append(operationdefinitionEntryList, entry)
	}

	var bundle models.OperationDefinitionBundle
	bundle.Type = "Bundle"
	bundle.Title = "OperationDefinition Index"
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Updated = time.Now()
	bundle.TotalResults = len(result)
	bundle.Entry = operationdefinitionEntryList

	log.Println("Setting operationdefinition search context")
	context.Set(r, "OperationDefinition", result)
	context.Set(r, "Resource", "OperationDefinition")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadOperationDefinition(r *http.Request) (*models.OperationDefinition, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("operationdefinitions")
	result := models.OperationDefinition{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting operationdefinition read context")
	context.Set(r, "OperationDefinition", result)
	context.Set(r, "Resource", "OperationDefinition")
	return &result, nil
}

func OperationDefinitionShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadOperationDefinition(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "OperationDefinition"))
}

func OperationDefinitionCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	operationdefinition := &models.OperationDefinition{}
	err := decoder.Decode(operationdefinition)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("operationdefinitions")
	i := bson.NewObjectId()
	operationdefinition.Id = i.Hex()
	err = c.Insert(operationdefinition)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting operationdefinition create context")
	context.Set(r, "OperationDefinition", operationdefinition)
	context.Set(r, "Resource", "OperationDefinition")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/OperationDefinition/"+i.Hex())
}

func OperationDefinitionUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	operationdefinition := &models.OperationDefinition{}
	err := decoder.Decode(operationdefinition)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("operationdefinitions")
	operationdefinition.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, operationdefinition)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting operationdefinition update context")
	context.Set(r, "OperationDefinition", operationdefinition)
	context.Set(r, "Resource", "OperationDefinition")
	context.Set(r, "Action", "update")
}

func OperationDefinitionDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("operationdefinitions")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting operationdefinition delete context")
	context.Set(r, "OperationDefinition", id.Hex())
	context.Set(r, "Resource", "OperationDefinition")
	context.Set(r, "Action", "delete")
}
