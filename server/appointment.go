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

func AppointmentIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.Appointment
	c := Database.C("appointments")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var appointmentEntryList []models.AppointmentBundleEntry
	for _, appointment := range result {
		var entry models.AppointmentBundleEntry
		entry.Title = "Appointment " + appointment.Id
		entry.Id = appointment.Id
		entry.Content = appointment
		appointmentEntryList = append(appointmentEntryList, entry)
	}

	var bundle models.AppointmentBundle
	bundle.Type = "Bundle"
	bundle.Title = "Appointment Index"
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Updated = time.Now()
	bundle.TotalResults = len(result)
	bundle.Entry = appointmentEntryList

	log.Println("Setting appointment search context")
	context.Set(r, "Appointment", result)
	context.Set(r, "Resource", "Appointment")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadAppointment(r *http.Request) (*models.Appointment, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("appointments")
	result := models.Appointment{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting appointment read context")
	context.Set(r, "Appointment", result)
	context.Set(r, "Resource", "Appointment")
	return &result, nil
}

func AppointmentShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadAppointment(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "Appointment"))
}

func AppointmentCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	appointment := &models.Appointment{}
	err := decoder.Decode(appointment)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("appointments")
	i := bson.NewObjectId()
	appointment.Id = i.Hex()
	err = c.Insert(appointment)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting appointment create context")
	context.Set(r, "Appointment", appointment)
	context.Set(r, "Resource", "Appointment")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/Appointment/"+i.Hex())
}

func AppointmentUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	appointment := &models.Appointment{}
	err := decoder.Decode(appointment)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("appointments")
	appointment.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, appointment)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting appointment update context")
	context.Set(r, "Appointment", appointment)
	context.Set(r, "Resource", "Appointment")
	context.Set(r, "Action", "update")
}

func AppointmentDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("appointments")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting appointment delete context")
	context.Set(r, "Appointment", id.Hex())
	context.Set(r, "Resource", "Appointment")
	context.Set(r, "Action", "delete")
}
