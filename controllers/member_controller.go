package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/notflex/models"
	model "github.com/notflex/models"
)

func UpdateMember(w http.ResponseWriter, r *http.Request) {
	db := Connect()

	err := r.ParseForm()
	if err != nil {
		return
	}

	// Get email from logged member
	email := GetEmailFromToken(r)

	var user models.User
	db.Where("email=?", email).First(&user)

	if name := r.Form.Get("name"); name != "" {
		user.Name = name
	}
	if birthdate := r.Form.Get("birthdate"); birthdate != "" {
		user.BirthDate = birthdate
	}
	if gender := r.Form.Get("gender"); gender != "" {
		user.Gender = gender
	}

	result := db.Save(&user)

	var response models.UserResponse
	if result.Error == nil {
		response.Status = 200
		response.Message = "Success Update User Data"
	} else {
		response.Status = 400
		response.Message = "Update Failed"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	db := Connect()

	err := r.ParseForm()
	if err != nil {
		return
	}

	// Get email from logged member
	email := GetEmailFromToken(r)

	// Get current time
	now := time.Now()

	// Get subscribe type from request
	subscribe := r.Form.Get("subscribe")
	cardNumber := r.Form.Get("cardNumber")
	exp := r.Form.Get("exp")
	cvc := r.Form.Get("cvc")

	// Set response
	var response models.SubscribeResponse
	if subscribe == "Premium" || subscribe == "Basic" {
		// Update subscribe
		result := db.Model(&models.User{}).Where("email = ?", email).Updates(map[string]interface{}{"subscribe": subscribe, "sub_date": now})

		var creditCard models.Credit
		creditCard.CardNumber = cardNumber
		creditCard.Exp = exp
		creditCard.Cvc = cvc
		creditCard.UserID = email

		if result.Error == nil {
			// Add new credit card
			err := db.Save(&creditCard)
			if err.Error == nil {
				response.Status = 200
				response.Message = "Subscribe success"
				response.Type = subscribe
			} else {
				response.Status = 400
				response.Message = "Subscribe Failed | " + result.Error.Error()
			}
		} else {
			response.Status = 400
			response.Message = "Subscribe Failed | " + result.Error.Error()
		}
	} else {
		response.Status = 400
		response.Message = "Subscribe Failed | Invalid Subscribe Type (Basic | Premium)"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Unsubscribe(w http.ResponseWriter, r *http.Request) {
	db := Connect()

	email := GetEmailFromToken(r)
	var response models.UnsubscribeResponse
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err == nil {
		if user.Subscribe == "Basic" || user.Subscribe == "Premium" {
			result := db.Model(&model.User{}).Where("email = ?", email).Updates(map[string]interface{}{"subscribe": "", "sub_date": nil})

			if result.Error == nil {
				response.Status = 200
				response.Message = "We are so sad you unsubscribed :("
			} else {
				response.Status = 400
				response.Message = "Unsubscribe Failed!" + result.Error.Error()
			}
		} else {
			response.Status = 400
			response.Message = "You Haven't Subscribe or Have Unsubscribe before"
		}
	} else {
		response.Status = 400
		response.Message = "There is Something Wrong"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
