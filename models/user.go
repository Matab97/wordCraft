package models

import (
	"context"
	"errors"
	"goCourseProject/db"
	"goCourseProject/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `bson:"email" binding:"required"`
	Password string             `bson:"password" binding:"required"`
}

func (u *User) Save() error {
	var err error
	u.Password, err = utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	_, err = db.DB.Collection("users").InsertOne(context.Background(), u)
	return err
}

func (u *User) Authenticate() (primitive.ObjectID, error) {
	var storedUser User
	err := db.DB.Collection("users").FindOne(context.Background(), bson.M{"email": u.Email}).Decode(&storedUser)
	if err != nil {
		return primitive.NilObjectID, errors.New("user not found")
	}

	pwdIsValid := utils.CheckPassword(u.Password, storedUser.Password)
	if !pwdIsValid {
		return primitive.NilObjectID, errors.New("credentials do not match")
	}
	return storedUser.ID, nil
}

func GetAllUsers() ([]User, error) {
	var users []User
	cursor, err := db.DB.Collection("users").Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
