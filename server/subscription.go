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

func SubscriptionIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.Subscription
	c := Database.C("subscriptions")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var subscriptionEntryList []models.SubscriptionBundleEntry
	for _, subscription := range result {
		var entry models.SubscriptionBundleEntry
		entry.Title = "Subscription " + subscription.Id
		entry.Id = subscription.Id
		entry.Content = subscription
		subscriptionEntryList = append(subscriptionEntryList, entry)
	}

	var bundle models.SubscriptionBundle
	bundle.Type = "Bundle"
	bundle.Title = "Subscription Index"
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Updated = time.Now()
	bundle.TotalResults = len(result)
	bundle.Entry = subscriptionEntryList

	log.Println("Setting subscription search context")
	context.Set(r, "Subscription", result)
	context.Set(r, "Resource", "Subscription")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadSubscription(r *http.Request) (*models.Subscription, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("subscriptions")
	result := models.Subscription{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting subscription read context")
	context.Set(r, "Subscription", result)
	context.Set(r, "Resource", "Subscription")
	return &result, nil
}

func SubscriptionShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadSubscription(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "Subscription"))
}

func SubscriptionCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	subscription := &models.Subscription{}
	err := decoder.Decode(subscription)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("subscriptions")
	i := bson.NewObjectId()
	subscription.Id = i.Hex()
	err = c.Insert(subscription)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting subscription create context")
	context.Set(r, "Subscription", subscription)
	context.Set(r, "Resource", "Subscription")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/Subscription/"+i.Hex())
}

func SubscriptionUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	subscription := &models.Subscription{}
	err := decoder.Decode(subscription)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("subscriptions")
	subscription.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, subscription)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting subscription update context")
	context.Set(r, "Subscription", subscription)
	context.Set(r, "Resource", "Subscription")
	context.Set(r, "Action", "update")
}

func SubscriptionDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("subscriptions")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting subscription delete context")
	context.Set(r, "Subscription", id.Hex())
	context.Set(r, "Resource", "Subscription")
	context.Set(r, "Action", "delete")
}
