package cli

import (
	"fmt"

	"syspodhot/internal/app/usecase"
	"syspodhot/internal/pkg/password"
)

// UserController represents a controller for all auth routes.
type UserController struct {
	userUC *usecase.UserUC
}

// NewUserController returns a new instance of [UserController].
func NewUserController(userUC *usecase.UserUC) *UserController {
	return &UserController{
		userUC: userUC,
	}
}

// Create creates a new user.
func (c *UserController) Create() error {
	var (
		inputBody        = &userBody{}
		username, passwd string
	)
	// ask for username
	fmt.Print("Введите логин: ")
	fmt.Scan(&username)
	if err := inputBody.ParseUsername(username); err != nil {
		return err
	}
	// ask for password
	fmt.Println(`Пароль должен содержать:
	1) Как минимум 8 символов;
	2) Как минимум 1 цифру;
	3) Как минимум 1 заглавную букву;
	4) Как минимум 1 строчную букву;
	5) Как минимум 1 спецсимвол.`)
	fmt.Print("Введите пароль: ")
	passwd, err := password.ScanPassword()
	if err != nil {
		return fmt.Errorf("scan password: %w", err)
	}
	if err := inputBody.ParsePassword(passwd); err != nil {
		return err
	}

	// create a new admin
	userObj, err := c.userUC.Create(inputBody.Username, inputBody.Password)
	if err != nil {
		return err
	}
	fmt.Printf("Юзер %q успешно создан!\n", userObj.Username)
	return nil
}

// Delete deletes a user.
func (c *UserController) Delete() error {
	var (
		inputBody = &userBody{}
		username  string
	)
	// ask for ID
	fmt.Print("Введите логин юзера: ")
	fmt.Scan(&username)
	if err := inputBody.ParseUsername(username); err != nil {
		return err
	}
	// delete user
	if err := c.userUC.DeleteByUsername(inputBody.Username); err != nil {
		return err
	}
	fmt.Printf("Юзер %s успешно удалён!\n", inputBody.Username)
	return nil
}

// GetAll returna a slice with all users.
func (c *UserController) GetAll() error {
	// get all
	users, err := c.userUC.GetAll()
	if err != nil {
		return err
	}
	// no one user
	usersAmount := len(users)
	if usersAmount == 0 {
		fmt.Println("Никого не найдено")
		return nil
	}
	// print out all users
	fmt.Println("Найдено записей:", usersAmount)
	for idx, user := range users {
		fmt.Printf("%d: ID - %s | Логин - %s\n", idx+1, user.ID, user.Username)
	}
	return nil
}
