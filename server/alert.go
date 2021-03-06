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

func AlertIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.Alert
	c := Database.C("alerts")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var alertEntryList []models.AlertBundleEntry
	for _, alert := range result {
		var entry models.AlertBundleEntry
		entry.Title = "Alert " + alert.Id
		entry.Id = alert.Id
		entry.Content = alert
		alertEntryList = append(alertEntryList, entry)
	}

	var bundle models.AlertBundle
	bundle.Type = "Bundle"
	bundle.Title = "Alert Index"
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Updated = time.Now()
	bundle.TotalResults = len(result)
	bundle.Entry = alertEntryList

	log.Println("Setting alert search context")
	context.Set(r, "Alert", result)
	context.Set(r, "Resource", "Alert")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadAlert(r *http.Request) (*models.Alert, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("alerts")
	result := models.Alert{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting alert read context")
	context.Set(r, "Alert", result)
	context.Set(r, "Resource", "Alert")
	return &result, nil
}

func AlertShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadAlert(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "Alert"))
}

func AlertCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	alert := &models.Alert{}
	err := decoder.Decode(alert)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("alerts")
	i := bson.NewObjectId()
	alert.Id = i.Hex()
	err = c.Insert(alert)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting alert create context")
	context.Set(r, "Alert", alert)
	context.Set(r, "Resource", "Alert")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/Alert/"+i.Hex())
}

func AlertUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	alert := &models.Alert{}
	err := decoder.Decode(alert)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("alerts")
	alert.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, alert)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting alert update context")
	context.Set(r, "Alert", alert)
	context.Set(r, "Resource", "Alert")
	context.Set(r, "Action", "update")
}

func AlertDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("alerts")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting alert delete context")
	context.Set(r, "Alert", id.Hex())
	context.Set(r, "Resource", "Alert")
	context.Set(r, "Action", "delete")
}
