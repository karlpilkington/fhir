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

func NamespaceIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.Namespace
	c := Database.C("namespaces")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var namespaceEntryList []models.NamespaceBundleEntry
	for _, namespace := range result {
		var entry models.NamespaceBundleEntry
		entry.Title = "Namespace " + namespace.Id
		entry.Id = namespace.Id
		entry.Content = namespace
		namespaceEntryList = append(namespaceEntryList, entry)
	}

	var bundle models.NamespaceBundle
	bundle.Type = "Bundle"
	bundle.Title = "Namespace Index"
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Updated = time.Now()
	bundle.TotalResults = len(result)
	bundle.Entry = namespaceEntryList

	log.Println("Setting namespace search context")
	context.Set(r, "Namespace", result)
	context.Set(r, "Resource", "Namespace")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadNamespace(r *http.Request) (*models.Namespace, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("namespaces")
	result := models.Namespace{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting namespace read context")
	context.Set(r, "Namespace", result)
	context.Set(r, "Resource", "Namespace")
	return &result, nil
}

func NamespaceShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadNamespace(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "Namespace"))
}

func NamespaceCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	namespace := &models.Namespace{}
	err := decoder.Decode(namespace)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("namespaces")
	i := bson.NewObjectId()
	namespace.Id = i.Hex()
	err = c.Insert(namespace)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting namespace create context")
	context.Set(r, "Namespace", namespace)
	context.Set(r, "Resource", "Namespace")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/Namespace/"+i.Hex())
}

func NamespaceUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	namespace := &models.Namespace{}
	err := decoder.Decode(namespace)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("namespaces")
	namespace.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, namespace)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting namespace update context")
	context.Set(r, "Namespace", namespace)
	context.Set(r, "Resource", "Namespace")
	context.Set(r, "Action", "update")
}

func NamespaceDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("namespaces")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting namespace delete context")
	context.Set(r, "Namespace", id.Hex())
	context.Set(r, "Resource", "Namespace")
	context.Set(r, "Action", "delete")
}
