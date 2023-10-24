package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go-mongo-project/models"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	client *mongo.Client
}

func NewUserController(client *mongo.Client, err error) *UserController {
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
	}
	return &UserController{client}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userCollection := uc.client.Database("go-mongo-project").Collection("users")
	ctx := context.TODO()

	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	uj, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userCollection := uc.client.Database("go-mongo-project").Collection("users")
	ctx := context.TODO()

	_, err := userCollection.InsertOne(ctx, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uj, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userCollection := uc.client.Database("go-mongo-project").Collection("users")
	ctx := context.TODO()

	_, err = userCollection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted user ", oid, "\n")
}
