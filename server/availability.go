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

func AvailabilityIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.Availability
	c := Database.C("availabilitys")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var availabilityEntryList []models.AvailabilityBundleEntry
	for _, availability := range result {
		var entry models.AvailabilityBundleEntry
		entry.Title = "Availability " + availability.Id
		entry.Id = availability.Id
		entry.Content = availability
		availabilityEntryList = append(availabilityEntryList, entry)
	}

	var bundle models.AvailabilityBundle
	bundle.Type = "Bundle"
	bundle.Title = "Availability Index"
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Updated = time.Now()
	bundle.TotalResults = len(result)
	bundle.Entry = availabilityEntryList

	log.Println("Setting availability search context")
	context.Set(r, "Availability", result)
	context.Set(r, "Resource", "Availability")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadAvailability(r *http.Request) (*models.Availability, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("availabilitys")
	result := models.Availability{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting availability read context")
	context.Set(r, "Availability", result)
	context.Set(r, "Resource", "Availability")
	return &result, nil
}

func AvailabilityShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadAvailability(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "Availability"))
}

func AvailabilityCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	availability := &models.Availability{}
	err := decoder.Decode(availability)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("availabilitys")
	i := bson.NewObjectId()
	availability.Id = i.Hex()
	err = c.Insert(availability)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting availability create context")
	context.Set(r, "Availability", availability)
	context.Set(r, "Resource", "Availability")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/Availability/"+i.Hex())
}

func AvailabilityUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	availability := &models.Availability{}
	err := decoder.Decode(availability)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("availabilitys")
	availability.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, availability)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting availability update context")
	context.Set(r, "Availability", availability)
	context.Set(r, "Resource", "Availability")
	context.Set(r, "Action", "update")
}

func AvailabilityDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("availabilitys")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting availability delete context")
	context.Set(r, "Availability", id.Hex())
	context.Set(r, "Resource", "Availability")
	context.Set(r, "Action", "delete")
}
