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

func RiskAssessmentIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.RiskAssessment
	c := Database.C("riskassessments")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var riskassessmentEntryList []models.RiskAssessmentBundleEntry
	for _, riskassessment := range result {
		var entry models.RiskAssessmentBundleEntry
		entry.Title = "RiskAssessment " + riskassessment.Id
		entry.Id = riskassessment.Id
		entry.Content = riskassessment
		riskassessmentEntryList = append(riskassessmentEntryList, entry)
	}

	var bundle models.RiskAssessmentBundle
	bundle.Type = "Bundle"
	bundle.Title = "RiskAssessment Index"
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Updated = time.Now()
	bundle.TotalResults = len(result)
	bundle.Entry = riskassessmentEntryList

	log.Println("Setting riskassessment search context")
	context.Set(r, "RiskAssessment", result)
	context.Set(r, "Resource", "RiskAssessment")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadRiskAssessment(r *http.Request) (*models.RiskAssessment, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("riskassessments")
	result := models.RiskAssessment{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting riskassessment read context")
	context.Set(r, "RiskAssessment", result)
	context.Set(r, "Resource", "RiskAssessment")
	return &result, nil
}

func RiskAssessmentShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadRiskAssessment(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "RiskAssessment"))
}

func RiskAssessmentCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	riskassessment := &models.RiskAssessment{}
	err := decoder.Decode(riskassessment)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("riskassessments")
	i := bson.NewObjectId()
	riskassessment.Id = i.Hex()
	err = c.Insert(riskassessment)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting riskassessment create context")
	context.Set(r, "RiskAssessment", riskassessment)
	context.Set(r, "Resource", "RiskAssessment")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/RiskAssessment/"+i.Hex())
}

func RiskAssessmentUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	riskassessment := &models.RiskAssessment{}
	err := decoder.Decode(riskassessment)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("riskassessments")
	riskassessment.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, riskassessment)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting riskassessment update context")
	context.Set(r, "RiskAssessment", riskassessment)
	context.Set(r, "Resource", "RiskAssessment")
	context.Set(r, "Action", "update")
}

func RiskAssessmentDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("riskassessments")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting riskassessment delete context")
	context.Set(r, "RiskAssessment", id.Hex())
	context.Set(r, "Resource", "RiskAssessment")
	context.Set(r, "Action", "delete")
}
