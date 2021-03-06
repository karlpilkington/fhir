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

func EncounterIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.Encounter
	c := Database.C("encounters")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var encounterEntryList []models.EncounterBundleEntry
	for _, encounter := range result {
		var entry models.EncounterBundleEntry
		entry.Title = "Encounter " + encounter.Id
		entry.Id = encounter.Id
		entry.Content = encounter
		encounterEntryList = append(encounterEntryList, entry)
	}

	var bundle models.EncounterBundle
	bundle.Type = "Bundle"
	bundle.Title = "Encounter Index"
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Updated = time.Now()
	bundle.TotalResults = len(result)
	bundle.Entry = encounterEntryList

	log.Println("Setting encounter search context")
	context.Set(r, "Encounter", result)
	context.Set(r, "Resource", "Encounter")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadEncounter(r *http.Request) (*models.Encounter, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("encounters")
	result := models.Encounter{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting encounter read context")
	context.Set(r, "Encounter", result)
	context.Set(r, "Resource", "Encounter")
	return &result, nil
}

func EncounterShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadEncounter(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "Encounter"))
}

func EncounterCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	encounter := &models.Encounter{}
	err := decoder.Decode(encounter)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("encounters")
	i := bson.NewObjectId()
	encounter.Id = i.Hex()
	err = c.Insert(encounter)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting encounter create context")
	context.Set(r, "Encounter", encounter)
	context.Set(r, "Resource", "Encounter")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/Encounter/"+i.Hex())
}

func EncounterUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	encounter := &models.Encounter{}
	err := decoder.Decode(encounter)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("encounters")
	encounter.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, encounter)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting encounter update context")
	context.Set(r, "Encounter", encounter)
	context.Set(r, "Resource", "Encounter")
	context.Set(r, "Action", "update")
}

func EncounterDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("encounters")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting encounter delete context")
	context.Set(r, "Encounter", id.Hex())
	context.Set(r, "Resource", "Encounter")
	context.Set(r, "Action", "delete")
}
