package controllers

import (
	"net/http"

	"../Models"
)

//List all followers
func GetFollowers(c *gin.Context) {
	var follower []Models.Follower
	err := Models.GetAllFollowers(&follower)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, follower)
	}
}

//Create a follower
func CreateAFollower(c *gin.Context) {
	var follower Models.Follower
	c.BindJSON(&follower)
	err := Models.CreateAFollower(&follower)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, follower)
	}
}

//Get a particular
func GetAFollower(c *gin.Context) {
	name := c.Params.ByName("name")
	var follower Models.Follower
	err := Models.GetAFollower(&follower, name)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, follower)
	}
}

// Update an existing follower
func UpdateAFollower(c *gin.Context) {
	var follower Models.Follower
	name := c.Params.ByName("name")
	err := Models.GetAFollower(&follower, name)
	if err != nil {
		c.JSON(http.StatusNotFound, follower)
	}
	c.BindJSON(&follower)
	err = Models.UpdateAFollower(&follower, name)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, follower)
	}
}

// Delete a follower
func DeleteAFollower(c *gin.Context) {
	var follower Models.Follower
	name := c.Params.ByName("name")
	err := Models.DeleteAFollower(&follower, name)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, gin.H{"name:" + name: "deleted"})
	}
}
