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

func PractitionerIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.Practitioner
	c := Database.C("practitioners")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var practitionerEntryList []models.PractitionerBundleEntry
	for _, practitioner := range result {
		var entry models.PractitionerBundleEntry
		entry.Title = "Practitioner " + practitioner.Id
		entry.Id = practitioner.Id
		entry.Content = practitioner
		practitionerEntryList = append(practitionerEntryList, entry)
	}

	var bundle models.PractitionerBundle
	bundle.Type = "Bundle"
	bundle.Title = "Practitioner Index"
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Updated = time.Now()
	bundle.TotalResults = len(result)
	bundle.Entry = practitionerEntryList

	log.Println("Setting practitioner search context")
	context.Set(r, "Practitioner", result)
	context.Set(r, "Resource", "Practitioner")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadPractitioner(r *http.Request) (*models.Practitioner, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("practitioners")
	result := models.Practitioner{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting practitioner read context")
	context.Set(r, "Practitioner", result)
	context.Set(r, "Resource", "Practitioner")
	return &result, nil
}

func PractitionerShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadPractitioner(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "Practitioner"))
}

func PractitionerCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	practitioner := &models.Practitioner{}
	err := decoder.Decode(practitioner)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("practitioners")
	i := bson.NewObjectId()
	practitioner.Id = i.Hex()
	err = c.Insert(practitioner)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting practitioner create context")
	context.Set(r, "Practitioner", practitioner)
	context.Set(r, "Resource", "Practitioner")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/Practitioner/"+i.Hex())
}

func PractitionerUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	practitioner := &models.Practitioner{}
	err := decoder.Decode(practitioner)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("practitioners")
	practitioner.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, practitioner)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting practitioner update context")
	context.Set(r, "Practitioner", practitioner)
	context.Set(r, "Resource", "Practitioner")
	context.Set(r, "Action", "update")
}

func PractitionerDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("practitioners")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting practitioner delete context")
	context.Set(r, "Practitioner", id.Hex())
	context.Set(r, "Resource", "Practitioner")
	context.Set(r, "Action", "delete")
}
