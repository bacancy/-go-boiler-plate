package user

import (
	"bacancy/go-boiler-plate/app/common"
	"bacancy/go-boiler-plate/app/models/user"
	"bacancy/go-boiler-plate/app/security"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	name := c.PostForm("name")
	lastName := c.PostForm("last_name")

	if err := Validate(password, email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "description": err.Error()})
		return
	}

	salt, err := security.GetSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "description": err.Error()})
		return
	}

	hash, err := security.Hash(password, salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "description": err.Error()})
		return
	}

	u, err := user.Create(email, hash, string(salt[:]), name, lastName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "description": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "true", "description": u})
}

func EditProfile(c *gin.Context) {

	id, _ := c.MustGet("id").(uint)
	phone := c.PostForm("phone")
	name := c.PostForm("name")
	lastName := c.PostForm("last_name")

	_, found, err := user.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "description": err.Error()})
		return
	}
	if found == false {
		c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "description": "User not found"})
		return
	}

	err = user.ChangeProfileData(id, name, phone, lastName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "description": err.Error()})
		return
	}

	//Get user data (incluiding profile picture) and send it
	userData, _, _ := user.GetUserById(id)

	c.JSON(http.StatusOK, gin.H{"success": "true", "description": userData})

}

func GetUserProfile(c *gin.Context) {

	id := c.MustGet("id").(uint)

	userData, found, err := user.GetUserProfile(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "description": err.Error()})
		return
	}
	if found == false {
		c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "description": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "true", "description": userData})

}

func GetUserById(c *gin.Context) {

	id := c.Param("id")

	userID, parseErr := strconv.ParseUint(id, 10, 32)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": "false", "description": "Invalid user ID"})
		return
	}

	r := common.GetDatabase()

	userData := []user.User{}

	if err := r.Preload("Picture").Preload("Newsletter").Where("id = ?", userID).First(&userData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": "false", "description": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "true", "description": userData})

}
