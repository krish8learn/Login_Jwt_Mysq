package controllers

import(
	"../models"
	"../database"
	"github.com/gofiber/fiber"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

const(
	secretkey = "secret"
)

func Regisiter(c *fiber.Ctx) error{
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil{
		return  err
	}

	password ,err := bcrypt.GenerateFromPassword([]byte(data["password"]),14)
	if err != nil{
		return  err
	}
	user := models.User{
		Name: data["name"],
		Email: data["email"],
		Password: password,
	}

	database.DB.Create(&user)
	return c.JSON(user)
}

func Login(c *fiber.Ctx)error{
	var data map[string]string

	err := c.BodyParser(&data)
	if err != nil{
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	//wrong email
	if user.Id == 0{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message":"user not found",
		})
	}

	//wrong password
	if err := bcrypt.CompareHashAndPassword(user.Password,[]byte(data["password"])); err!= nil{
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message":"Wrong Password",
		})
	}

	//right email and password

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour*24).Unix(),
	})

	token, err := claims.SignedString([]byte(secretkey))
	if err != nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message":"could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:"jwt",
		Value: token,
		Expires:time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func User(c *fiber.Ctx) error{
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie,&jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error){
		return []byte(secretkey), nil
	})

	if err != nil{
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message":"unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)
	var user models.User

	database.DB.Where("id=?",claims.Issuer).First(&user)
	return c.JSON(user)
}

func LogOut(c *fiber.Ctx) error{
	//to log out we need to cancel the cookie which is done by removing the current by replacing with new one

	cookie := fiber.Cookie{
		Name:"jwt",
		Value:"",
		Expires: time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message":"Log out successfull",
	})
}