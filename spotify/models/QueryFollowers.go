package models

import (
	"fmt"

	"../Config"
	_ "github.com/go-sql-driver/mysql"
)

//fetch all followers at once
func GetAllFollowers(follower *[]Follower) (err error) {
	if err = Config.DB.Find(follower).Error; err != nil {
		return err
	}
	return nil
}

//insert a follower
func CreateAFollower(follower *Follower) (err error) {
	if err = Config.DB.Create(follower).Error; err != nil {
		return err
	}
	return nil
}

//fetch one follower
func GetAFollower(follower *Follower, name string) (err error) {
	if err := Config.DB.Where("name = ?", name).First(follower).Error; err != nil {
		return err
	}
	return nil
}

//update a follower
func UpdateAFollower(follower *Follower, name string) (err error) {
	fmt.Println(follower)
	Config.DB.Save(follower)
	return nil
}

//delete a follower
func DeleteAFollower(follower *Follower, name string) (err error) {
	Config.DB.Where("name = ?", name).Delete(follower)
	return nil
}
