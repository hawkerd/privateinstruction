package handlers

// // create a class
// func CreateClass(w http.ResponseWriter, r *http.Request) {
// 	// make sure the user is authenticated
// 	userID, ok := isUserAuthenticated(r)
// 	if !ok {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	// decode the request body into the class struct
// 	var class models.Class
// 	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
// 		http.Error(w, "err.Error()", http.StatusBadRequest)
// 		return
// 	}

// 	// validate the class input
// 	if class.Name == "" {
// 		class.Name = "Unnamed Class"
// 	}

// 	// create the new class
// 	if err := db.Create(&class).Error; err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// create a class member entry for the user
// 	classMember := models.ClassMember{
// 		ClassID: class.ID,
// 		UserID:  userID,
// 		Role:    "instructor",
// 	}

// 	if err := db.Create(&classMember).Error; err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// respond with success
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(class)
// }

// //func GetClasses(w http.ResponseWriter, r *http.Request) {
// //	// make sure the user is authenticated
// //	userID, ok := isUserAuthenticated(r)
// //	if !ok {
// //		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// //		return
// //	}
// //
// //	var classes []models.Class
// //	if err := db.Model(&models.ClassMember{}).
// //		Joins("JOIN classes ON class_members.class_id = classes.id").
// //		Where("class_members.user_id = ?", userID).
// //		Select("classes.*").
// //		Scan(&classes).Error; err != nil {
// //		http.Error(w, err.Error(), http.StatusInternalServerError)
// //		return
// //	}
// //
// //	w.WriteHeader(http.StatusOK)
// //	json.NewEncoder(w).Encode(classes)
// //}
