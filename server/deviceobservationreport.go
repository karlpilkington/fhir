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

func DeviceObservationReportIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.DeviceObservationReport
	c := Database.C("deviceobservationreports")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var deviceobservationreportEntryList []models.DeviceObservationReportBundleEntry
	for _, deviceobservationreport := range result {
		var entry models.DeviceObservationReportBundleEntry
		entry.Title = "DeviceObservationReport " + deviceobservationreport.Id
		entry.Id = deviceobservationreport.Id
		entry.Content = deviceobservationreport
		deviceobservationreportEntryList = append(deviceobservationreportEntryList, entry)
	}

	var bundle models.DeviceObservationReportBundle
	bundle.Type = "Bundle"
	bundle.Title = "DeviceObservationReport Index"
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Updated = time.Now()
	bundle.TotalResults = len(result)
	bundle.Entry = deviceobservationreportEntryList

	log.Println("Setting deviceobservationreport search context")
	context.Set(r, "DeviceObservationReport", result)
	context.Set(r, "Resource", "DeviceObservationReport")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadDeviceObservationReport(r *http.Request) (*models.DeviceObservationReport, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("deviceobservationreports")
	result := models.DeviceObservationReport{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting deviceobservationreport read context")
	context.Set(r, "DeviceObservationReport", result)
	context.Set(r, "Resource", "DeviceObservationReport")
	return &result, nil
}

func DeviceObservationReportShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadDeviceObservationReport(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "DeviceObservationReport"))
}

func DeviceObservationReportCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	deviceobservationreport := &models.DeviceObservationReport{}
	err := decoder.Decode(deviceobservationreport)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("deviceobservationreports")
	i := bson.NewObjectId()
	deviceobservationreport.Id = i.Hex()
	err = c.Insert(deviceobservationreport)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting deviceobservationreport create context")
	context.Set(r, "DeviceObservationReport", deviceobservationreport)
	context.Set(r, "Resource", "DeviceObservationReport")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/DeviceObservationReport/"+i.Hex())
}

func DeviceObservationReportUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	deviceobservationreport := &models.DeviceObservationReport{}
	err := decoder.Decode(deviceobservationreport)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("deviceobservationreports")
	deviceobservationreport.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, deviceobservationreport)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting deviceobservationreport update context")
	context.Set(r, "DeviceObservationReport", deviceobservationreport)
	context.Set(r, "Resource", "DeviceObservationReport")
	context.Set(r, "Action", "update")
}

func DeviceObservationReportDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("deviceobservationreports")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting deviceobservationreport delete context")
	context.Set(r, "DeviceObservationReport", id.Hex())
	context.Set(r, "Resource", "DeviceObservationReport")
	context.Set(r, "Action", "delete")
}
