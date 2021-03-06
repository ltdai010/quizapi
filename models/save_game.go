package models

import (
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const SAVE_GAME = "saveGame"

type SaveGame struct {
	UserID    		   string
	QuizID    		   string
	ListAnsweredQuest  []int
}

func AddSaveGame(saveGame SaveGame) string {
	//check if already exist
	list := client.Collection(SAVE_GAME).Where("UserID", "==", saveGame.UserID).Where("QuizID", "==", saveGame.QuizID).Documents(ctx)
	for {
		doc, err := list.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err.Error()
		}
		_, err = client.Collection(SAVE_GAME).Doc(doc.Ref.ID).Delete(ctx)
		if err != nil {
			return err.Error()
		}
	}
	//add save game
	ref, _, err := client.Collection(SAVE_GAME).Add(ctx, saveGame)
	if err != nil {
		return err.Error()
	}
	return ref.ID
}

func GetSaveGame(id string) (*SaveGame, error) {
	ref := client.Collection(SAVE_GAME).Doc(id)
	doc, err := ref.Get(ctx)
	if err != nil {
		return nil, err
	}
	sg := &SaveGame{}
	err = doc.DataTo(sg)
	return sg, err
}

func GetSaveGameByUserQuiz(userID, quizID string) (SaveGame, error) {
	list := client.Collection(SAVE_GAME).Where("UserID", "==", userID).Where("QuizID", "==", quizID).Documents(ctx)
	doc, err := list.Next()
	if err == iterator.Done {
		return SaveGame{}, err
	}
	if err != nil {
		return SaveGame{}, err
	}
	saveGame := SaveGame{}
	err = doc.DataTo(&saveGame)
	return saveGame, err
}

func DeleteSaveGame(id string) error {
	_, err := client.Collection(SAVE_GAME).Doc(id).Delete(ctx)
	return err
}

func GetAllSaveGameByUser(userID string) (map[string]*SaveGame, error) {
	iter := client.Collection(SAVE_GAME).Where("UserID", "==", userID).Documents(ctx)
	mapS := make(map[string]*SaveGame)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		s := &SaveGame{}
		err = doc.DataTo(s)
		if err != nil {
			return nil, err
		}
		mapS[doc.Ref.ID] = s
	}
	return mapS, nil
}

func UpdateSaveGame(id string, list []int) error {
	_, err := client.Collection(SAVE_GAME).Doc(id).Set(ctx, map[string]interface{}{
		"ListAnsweredQuest" : list,
	}, firestore.MergeAll)
	return err
}