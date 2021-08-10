package data

import (
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type (
	FollowInformation struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
)

func (ur *UserRepo) GetFollowingList(username string) ([]*FollowInformation, error) {
	session := ur.DB.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	query := "MATCH (a:User{username:$username})-[r:Following]->(n:User) RETURN n"
	result, err := session.Run(query, map[string]interface{}{
		"username": username,
	})
	if err != nil {
		return []*FollowInformation{}, err
	}
	users := []*FollowInformation{}
	for result.Next() {
		record := result.Record()
		if value, ok := record.Get("n"); ok {
			node := value.(neo4j.Node)
			props := node.Props
			user := FollowInformation{}
			if err := mapstructure.Decode(props, &user); err != nil {
				return []*FollowInformation{}, nil
			}
			users = append(users, &user)

		}
	}
	return users, nil
}

func (ur *UserRepo) GetFollowersList(username string) ([]*FollowInformation, error) {
	session := ur.DB.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	query := "MATCH (a:User{username:$username})<-[r:Following]-(n:User) RETURN n"
	result, err := session.Run(query, map[string]interface{}{
		"username": username,
	})
	if err != nil {
		return []*FollowInformation{}, err
	}
	users := []*FollowInformation{}
	for result.Next() {
		record := result.Record()
		if value, ok := record.Get("n"); ok {
			node := value.(neo4j.Node)
			props := node.Props
			user := FollowInformation{}
			if err := mapstructure.Decode(props, &user); err != nil {
				return []*FollowInformation{}, err
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