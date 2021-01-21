package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"

	c "github.com/nikhilpatne/provider-management/config"
	d "github.com/nikhilpatne/provider-management/dao"
	. "github.com/nikhilpatne/provider-management/models"
	"github.com/nikhilpatne/provider-management/providers/aws"
)

var config = c.Config{}
var dao = d.ProviderDAO{}

// GetAllBooks ...
func GetAllProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, providers)
}

// FindBook ...
func FindProvider(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	provider, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Provider ID")
		return
	}
	respondWithJson(w, http.StatusOK, provider)
}

// AddnewBook ...
func AddnewProvider(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var provider Provider
	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	connection := aws.CheckConnection(provider.Connection["access_key"], provider.Connection["secret_access_key"])
	if !connection {
		respondWithError(w, http.StatusUnauthorized, "Authentication failed, please check your credentials")
		return
	}
	provider.ID = bson.NewObjectId()
	if err := dao.Insert(provider); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, provider)
}

// PUT update an existing movie
func UpdateProvider(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var provider Provider
	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := dao.Update(provider); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"message": "record updated successfully"})
}

// DELETE an existing movie
func DeleteProvider(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var provider Provider
	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(provider); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"message": "book deleted successfully"})
}

//  VmDeployment ....
func VmDeployment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var awsProvider AWSProvider
	err := json.NewDecoder(r.Body).Decode(&awsProvider)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	params := mux.Vars(r)

	provider, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Provider ID")
		return
	}
	isInstanceLaunched, err := aws.LaunchInstance(provider.Connection["access_key"], provider.Connection["secret_access_key"], awsProvider.Image, awsProvider.Size,awsProvider.VMname)
	if isInstanceLaunched {
		respondWithJson(w, http.StatusCreated, map[string]string{"message": "virtual machine deployed successfully"})
	} else {
		respondWithError(w, http.StatusBadRequest, "Something went wrong"+err.Error())
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

// Define HTTP request routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/providers", GetAllProviders).Methods("GET")
	r.HandleFunc("/api/provider", AddnewProvider).Methods("POST")
	r.HandleFunc("/api/provider", UpdateProvider).Methods("PUT")
	r.HandleFunc("/api/provider", DeleteProvider).Methods("DELETE")
	r.HandleFunc("/api/provider/{id}", FindProvider).Methods("GET")

	r.HandleFunc("/api/vmdeployment/{id}", VmDeployment).Methods("POST")
	if err := http.ListenAndServe(":9090", r); err != nil {
		log.Fatal(err)
	}
}
