package data

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-playground/validator"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		Username   string `json:"username" validate:"required"`
		Name       string `json:"name" validate:"required"`
		ProfileURI string `json:"profile_uri"`
		CoverURI   string `json:"cover_uri"`
		Password   string `json:"password" validate:"required"`
		Email      string `json:"email" validate:"required,email"`
		UserType   int    `json:"user_type" validate:"gte=0,lte=1"`
	}

	UserRepo struct {
		DB neo4j.Driver
	}
)

func NewUserRepo() *UserRepo {
	driver, err := neo4j.NewDriver(os.Getenv("DB"), neo4j.BasicAuth(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), ""))
	if err != nil {
		panic(err)
	}
	init := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer init.Close()
	_, err = init.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(
			"CREATE CONSTRAINT uniqueUsername IF NOT EXISTS ON (p:User) ASSERT p.username IS UNIQUE",
			nil,
		)
		return nil, err
	})
	if err != nil {
		panic(err)
	}
	_, err = init.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(
			"CREATE CONSTRAINT uniqueEmail IF NOT EXISTS ON (p:User) ASSERT p.email IS UNIQUE",
			nil,
		)
		return nil, err
	})
	if err != nil {
		panic(err)
	}
	return &UserRepo{driver}
}

func (ur *UserRepo) CreateUser(u *User) error {
	u.Password, _ = hashPassword(u.Password)
	session := ur.DB.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"CREATE (p:User{username:$username,name:$name, profileUri: $profileUri, coverUri: $coverUri, password:$password,email:$email,type:$type}) RETURN p",
			map[string]interface{}{
				"username":   u.Username,
				"name":       u.Name,
				"profileUri": u.ProfileURI,
				"coverUri":   u.CoverURI,
				"password":   u.Password,
				"email":      u.Email,
				"type":       u.UserType,
			},
		)
		if err != nil {
			return nil, err
		}
		if result.Next() {
			return result.Record().Values[0], nil
		}
		return nil, result.Err()
	})
	return err
}

func (ur *UserRepo) GetUser(username string) (*User, error) {
	session := ur.DB.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	query := "MATCH (a:User) WHERE a.username = $username RETURN a"
	result, err := session.Run(query, map[string]interface{}{
		"username": username,
	})
	if err != nil {
		return nil, err
	}
	user := User{}
	if result.Next() {
		record := result.Record()
		if value, ok := record.Get("a"); ok {
			node := value.(neo4j.Node)
			props := node.Props
			if err := mapstructure.Decode(props, &user); err != nil {
				return nil, err
			}
		}
	}
	if user.Username == "" {
		return nil, fmt.Errorf("no user found")
	}
	return &user, nil

}

func (ur *UserRepo) DeleteUser(username string) error {
	session := ur.DB.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	query := "MATCH (n:User) WHERE n.username = $username DETACH DELETE n"
	_, err := session.Run(query, map[string]interface{}{
		"username": username,
	})
	return err

}

func (ur *UserRepo) UpdateUser(username string, updateInfo map[string]string) error {
	session := ur.DB.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	user, err := ur.GetUser(username)
	if err != nil {
		return err
	}
	userInfoBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	userMap := map[string]interface{}{}
	err = json.Unmarshal(userInfoBytes, &userMap)
	if err != nil {
		return err
	}

	for key, value := range updateInfo {
		if _, ok := userMap[key]; ok {
			switch key {
			case "username":
				user.Username = value
			case "name":
				user.Name = value
			case "profile_uri":
				user.ProfileURI = value
			case "cover_uri":
				user.CoverURI = value
			case "password":
				hashP, err := hashPassword(value)
				if err != nil {
					return err
				}
				user.Password = hashP
			case "email":
				user.Email = value
			default:
				return fmt.Errorf("provided field %s is not a valid property", key)
			}
		}
	}
	if err := user.Validate(); err != nil {
		return err
	}
	query := `MATCH (a:User{username: $username})
			  SET a.username = $newUserName
			  SET a.profileUri = $profileUri
			  SET a.coverUri = $coverUri
			  SET a.password = $password
			  SET a.email = $email`
	_, err = session.Run(query, map[string]interface{}{
		"username":    username,
		"newUserName": user.Username,
		"profileUri":  user.ProfileURI,
		"coverUri":    user.CoverURI,
		"password":    user.Password,
		"email":       user.Email,
	})

	return err

}

func (u *User) CheckPassword(provided string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(provided))
}

func (u *User) Validate() error {
	validator := validator.New()
	return validator.Struct(u)
}

func hashPassword(plainText string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), 14)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
