package data

import (
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
			"CREATE (p:User{username:$username,name:$name,password:$password,email:$email,type:$type}) RETURN p",
			map[string]interface{}{"username": u.Username, "name": u.Name, "password": u.Password, "email": u.Email, "type": u.UserType})
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
				return &User{}, err
			}
		}
	}
	return &user, nil

}

func (ur *UserRepo) GetFollowerList(username string) ([]*User, error) {
	session := ur.DB.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	query := "MATCH (a:User{username:$username})-[r:Following]->(n:User) RETURN n"
	result, err := session.Run(query, map[string]interface{}{
		"username": username,
	})
	if err != nil {
		return nil, err
	}
	users := []*User{}
	for result.Next() {
		record := result.Record()
		if value, ok := record.Get("n"); ok {
			node := value.(neo4j.Node)
			props := node.Props
			user := User{}
			if err := mapstructure.Decode(props, &user); err != nil {
				return []*User{}, nil
			}
			users = append(users, &user)

		}
	}
	return users, nil
}

func (ur *UserRepo) FollowUser(user, following string) error {
	session := ur.DB.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(
			"MATCH (a:User),(b:User) WHERE a.username = $user AND b.username = $following CREATE (a)-[r:Following]->(b) RETURN type(r)",
			map[string]interface{}{"user": user, "following": following},
		)
		return nil, err
	})
	return err
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

func (ur *UserRepo) UnFollowUser(user, following string) error {
	session := ur.DB.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	query := "MATCH (a:User{username:$username})-[r:Following]->(b:User{username:$following}) DELETE r"
	_, err := session.Run(query, map[string]interface{}{
		"username":  user,
		"following": following,
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
