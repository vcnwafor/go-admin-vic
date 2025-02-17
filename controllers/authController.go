package controllers

import (
	"strconv"
	"time"

	"github.com/Rajil1213/go-admin/database"
	"github.com/Rajil1213/go-admin/models"
	"github.com/Rajil1213/go-admin/util"
	"github.com/gofiber/fiber/v2"
)

// register body
type RegisterBodyData struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func Register(c *fiber.Ctx) error {
	var data RegisterBodyData

	if err := c.BodyParser(&data); err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if data.Password != data.PasswordConfirm {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Passwords do not match",
		})
	}

	user := models.User{
		FirstName: data.FirstName,
		Lastname:  data.LastName,
		Email:     data.Email,
		RoleId:    3,
	}

	user.SetPassword(data.Password)

	// insert to database
	database.DB.Create(&user)

	// return inserted value
	return c.JSON(user)
}

// login body
type LoginBodyData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	var data LoginBodyData

	if err := c.BodyParser(&data); err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var user models.User
	database.DB.Where("email = ?", data.Email).First(&user)

	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "Incorrect email address",
		})
	}

	if err := user.CheckPassword(data.Password); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}

	token, err := util.GenerateJwt(strconv.Itoa(int(user.Id)))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// store token in a cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true, // cannot be accessed by frontend
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Success",
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	id, _ := util.ParseJwt(cookie)

	var user models.User
	database.DB.Where("id = ?", id).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	// remove cookie by setting expiration in the past
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

// update body data
type UpdateBodyData struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func UpdateInfo(c *fiber.Ctx) error {
	var data UpdateBodyData

	if err := c.BodyParser(&data); err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	cookie := c.Cookies("jwt")
	id, _ := util.ParseJwt(cookie)
	userId, _ := strconv.Atoi(id)

	user := models.User{
		Id:        uint(userId),
		FirstName: data.FirstName,
		Lastname:  data.LastName,
		Email:     data.Email,
	}

	database.DB.Model(&user).Updates(&user)

	return c.JSON(user)
}

// update password body
type UpdatePasswordBody struct {
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func UpdatePassword(c *fiber.Ctx) error {
	var data UpdatePasswordBody

	if err := c.BodyParser(&data); err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if data.Password != data.PasswordConfirm {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Passwords do not match",
		})
	}

	cookie := c.Cookies("jwt")
	id, _ := util.ParseJwt(cookie)
	userId, _ := strconv.Atoi(id)

	user := models.User{
		Id: uint(userId),
	}
	user.SetPassword(data.Password)

	database.DB.Model(&user).Updates(user)

	return c.JSON(user)
}

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
