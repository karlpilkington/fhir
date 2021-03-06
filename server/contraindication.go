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

func ContraindicationIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.Contraindication
	c := Database.C("contraindications")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var contraindicationEntryList []models.ContraindicationBundleEntry
	for _, contraindication := range result {
		var entry models.ContraindicationBundleEntry
		entry.Title = "Contraindication " + contraindication.Id
		entry.Id = contraindication.Id
		entry.Content = contraindication
		contraindicationEntryList = append(contraindicationEntryList, entry)
	}

	var bundle models.ContraindicationBundle
	bundle.Type = "Bundle"
	bundle.Title = "Contraindication Index"
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Updated = time.Now()
	bundle.TotalResults = len(result)
	bundle.Entry = contraindicationEntryList

	log.Println("Setting contraindication search context")
	context.Set(r, "Contraindication", result)
	context.Set(r, "Resource", "Contraindication")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadContraindication(r *http.Request) (*models.Contraindication, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("contraindications")
	result := models.Contraindication{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting contraindication read context")
	context.Set(r, "Contraindication", result)
	context.Set(r, "Resource", "Contraindication")
	return &result, nil
}

func ContraindicationShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadContraindication(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "Contraindication"))
}

func ContraindicationCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	contraindication := &models.Contraindication{}
	err := decoder.Decode(contraindication)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("contraindications")
	i := bson.NewObjectId()
	contraindication.Id = i.Hex()
	err = c.Insert(contraindication)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting contraindication create context")
	context.Set(r, "Contraindication", contraindication)
	context.Set(r, "Resource", "Contraindication")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/Contraindication/"+i.Hex())
}

func ContraindicationUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	contraindication := &models.Contraindication{}
	err := decoder.Decode(contraindication)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("contraindications")
	contraindication.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, contraindication)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting contraindication update context")
	context.Set(r, "Contraindication", contraindication)
	context.Set(r, "Resource", "Contraindication")
	context.Set(r, "Action", "update")
}

func ContraindicationDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("contraindications")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting contraindication delete context")
	context.Set(r, "Contraindication", id.Hex())
	context.Set(r, "Resource", "Contraindication")
	context.Set(r, "Action", "delete")
}
