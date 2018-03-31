package userCenter

import (
	"github.com/jinzhu/gorm"
	"github.com/firerainos/firerain-web-go/core"
)

type User struct {
	gorm.Model
	Username string
	Password string
	Email string
	Group []Group `gorm:"many2many:user_group"`
}

func AddUser(username,password,email string,group []string) error {
	g,err := GetGroupByNames(group)
	if err != nil {
		return err
	}

	db,err := core.GetSqlConn()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Create(&User{Username:username,Password:password,Email:email,Group:g}).Error
}

func GetUser() ([]User,error) {
	var users []User

	db,err := core.GetSqlConn()
	if err != nil {
		return users,err
	}
	defer db.Close()

	db.Preload("Group").Find(&users)

	return users, nil
}

func GetUserByName(name string) (User,error) {
	var user User

	db,err := core.GetSqlConn()
	if err != nil {
		return user,err
	}
	defer db.Close()

	if db.Where("name = ?",name).Preload("Group").First(&user).RecordNotFound() {
		return user, db.Error
	}

	return user, nil
}

func GetUserById(id int) (User,error) {
	var user User

	db,err := core.GetSqlConn()
	if err != nil {
		return user,err
	}
	defer db.Close()

	if db.Where("id = ?",id).Preload("Group").First(&user).RecordNotFound() {
		return user, db.Error
	}

	return user, nil
}

func (user User) Delete() error {
	db,err := core.GetSqlConn()
	if err != nil {
		return err
	}
	defer db.Close()

	db.Delete(user)

	return nil
}

func (user User) AddGroup(group string) error {
	g,err := GetGroupByName(group)
	if err != nil {
		return err
	}

	user.Group = append(user.Group, g)

	db,err := core.GetSqlConn()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Save(user).Error
}

func (user User) DeleteGroup(group string) error {
	var groups []Group
	for _,g := range user.Group {
		if g.Name != group {
			groups = append(groups, g)
		}
	}

	user.Group = groups

	db,err := core.GetSqlConn()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Save(user).Error
}