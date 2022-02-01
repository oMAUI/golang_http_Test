package User

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type (
	User struct {
		CreatedAt   time.Time `json:"created_at"`
		DisplayName string    `json:"display_name"`
		Email       string    `json:"email"`
	}

	UserList struct {
		Increment int             `json:"increment"`
		List      map[string]User `json:"list"`
	}
)

func GetUsersInData(Data string)  (UserList, error) {
	data, errGetData := ioutil.ReadFile(Data)
	if errGetData != nil {
		return UserList{}, fmt.Errorf("failed to get data: %w", errGetData)
	}

	var users UserList
	if errGetUsers := json.Unmarshal(data, &users); errGetUsers != nil {
		return UserList{}, fmt.Errorf("failed to unmarshal: %w", errGetUsers)
	}

	return users, nil
}

func UpdateUser(UserFromBody User, User *User){
	if UserFromBody.DisplayName != "" {
		User.DisplayName = UserFromBody.DisplayName
	}
	if UserFromBody.Email != "" {
		User.Email = UserFromBody.Email
	}
}