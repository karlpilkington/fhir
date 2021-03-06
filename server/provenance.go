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

func ProvenanceIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.Provenance
	c := Database.C("provenances")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var provenanceEntryList []models.ProvenanceBundleEntry
	for _, provenance := range result {
		var entry models.ProvenanceBundleEntry
		entry.Title = "Provenance " + provenance.Id
		entry.Id = provenance.Id
		entry.Content = provenance
		provenanceEntryList = append(provenanceEntryList, entry)
	}

	var bundle models.ProvenanceBundle
	bundle.Type = "Bundle"
	bundle.Title = "Provenance Index"
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Updated = time.Now()
	bundle.TotalResults = len(result)
	bundle.Entry = provenanceEntryList

	log.Println("Setting provenance search context")
	context.Set(r, "Provenance", result)
	context.Set(r, "Resource", "Provenance")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadProvenance(r *http.Request) (*models.Provenance, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("provenances")
	result := models.Provenance{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting provenance read context")
	context.Set(r, "Provenance", result)
	context.Set(r, "Resource", "Provenance")
	return &result, nil
}

func ProvenanceShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadProvenance(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "Provenance"))
}

func ProvenanceCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	provenance := &models.Provenance{}
	err := decoder.Decode(provenance)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("provenances")
	i := bson.NewObjectId()
	provenance.Id = i.Hex()
	err = c.Insert(provenance)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting provenance create context")
	context.Set(r, "Provenance", provenance)
	context.Set(r, "Resource", "Provenance")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/Provenance/"+i.Hex())
}

func ProvenanceUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	provenance := &models.Provenance{}
	err := decoder.Decode(provenance)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("provenances")
	provenance.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, provenance)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting provenance update context")
	context.Set(r, "Provenance", provenance)
	context.Set(r, "Resource", "Provenance")
	context.Set(r, "Action", "update")
}

func ProvenanceDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("provenances")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting provenance delete context")
	context.Set(r, "Provenance", id.Hex())
	context.Set(r, "Resource", "Provenance")
	context.Set(r, "Action", "delete")
}
