package main

import (
	"errors"
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type server struct {
	ID           string `bson:"id"`
	Prefix       string `bson:"prefix"`
	AnimeChannel string `bson:"aChannel"`
	Mode         string `bson:"mode"`
	GuildName    string `bson:"gName"`
}

// updateAnimeChannel will update this servers aChannel field in the database
// with the newChannelID value
// ID of the server object has to be set
// updates the calling server with the updated server returned from the db
// Returns an error if there was one
func (s *server) updateAnimeChannel(newChannelID string) error {
	updateQuery := bson.M{
		"$set": bson.M{
			"id":       s.ID,
			"aChannel": newChannelID,
		},
	}

	change := mgo.Change{
		Update:    updateQuery,
		Upsert:    true,
		Remove:    false,
		ReturnNew: true,
	}

	_, err := DBserverList.Find(bson.M{"id": s.ID}).Apply(change, s)
	if err != nil {
		return err
	}

	return nil
}

// updateGuildName will update this servers gName field in the database
// with the guildName value
// ID of the server object has to be set
// updates the calling server with the updated server returned from the db
// Returns an error if there was one
func (s *server) updateGuildName(guildName string) error {
	updateQuery := bson.M{
		"$set": bson.M{
			"id":    s.ID,
			"gName": guildName,
		},
	}

	change := mgo.Change{
		Update:    updateQuery,
		Upsert:    true,
		Remove:    false,
		ReturnNew: true,
	}

	_, err := DBserverList.Find(bson.M{"id": s.ID}).Apply(change, s)
	if err != nil {
		return err
	}

	return nil
}

// updatePrefix will update this servers prefix field in the database
// with the prefix value
// ID of the server object has to be set
// updates the calling server with the updated server object returned from the db
// Returns an error if there was one
func (s *server) updatePrefix(prefix string) error {
	updateQuery := bson.M{
		"$set": bson.M{
			"id":     s.ID,
			"prefix": prefix,
		},
	}

	change := mgo.Change{
		Update:    updateQuery,
		Upsert:    true,
		Remove:    false,
		ReturnNew: true,
	}

	_, err := DBserverList.Find(bson.M{"id": s.ID}).Apply(change, s)
	if err != nil {
		return err
	}

	return nil
}

func (s *server) updateGuildMode(newMode string) error {
	updateQuery := bson.M{
		"$set": bson.M{
			"id":   s.ID,
			"mode": newMode,
		},
	}

	change := mgo.Change{
		Update:    updateQuery,
		Upsert:    true,
		Remove:    false,
		ReturnNew: true,
	}

	_, err := DBserverList.Find(bson.M{"id": s.ID}).Apply(change, s)
	if err != nil {
		return err
	}

	return nil
}

// delete this server from the db
func (s *server) delete() error {
	err := DBserverList.Remove(bson.M{"id": s.ID})
	if err != nil {
		return fmt.Errorf("Error trying to delete a guild from the database: %s - %s", s.ID, err)
	}
	return nil
}

// fetch fetches this server document from the db by its ID
func (s *server) fetch() error {
	err := DBserverList.Find(bson.M{"id": s.ID}).One(s)
	if err != nil {
		return errors.New("Error trying to database find: - " + err.Error())
	}

	if s.Prefix == "" {
		s.Prefix = "!"
	}

	return nil
}
