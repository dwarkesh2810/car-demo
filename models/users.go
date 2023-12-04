package models

import (
	"car_demo/helper"
	"car_demo/request"

	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Users struct {
	Id        int64  `orm:"auto"`
	FirstName string `orm:"size(128)" json:"first_name" form:"first_name"`
	LastName  string `orm:"size(128)" json:"last_name" form:"last_name"`
	Email     string `orm:"size(128)" json:"email" form:"email"`
	Mobile    string `json:"mobile" form:"mobile"`
	Password  string `orm:"size(128)" json:"password" form:"password"`
	Status    int    `orm:"" json:"status" form:"status"`
	Role      string `orm:"size(20)" json:"role" form:"role"`
	Otp       string `orm:"size(20)" json:"otp" form:"otp"`
	CreatedAt int64
	UpdatedAt int64
}

func init() {
	orm.RegisterModel(new(Users))
}

// AddUsers insert a new Users into database and returns
// last inserted Id on success.
func AddUsers(m *request.CreateUserRequest) (int64, error) {
	o := orm.NewOrm()

	up, err := helper.HashData(m.Password)
	if err != nil {

		return 0, err
	}

	user := Users{
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Mobile:    m.Mobile,
		Email:     m.Email,
		Status:    0,
		Role:      m.Role,
		Password:  up,
		CreatedAt: time.Now().UnixMilli(),
	}

	id, err := o.Insert(&user)
	if err != nil {
		return 0, errors.New("Unexpected database error" + err.Error())
	}
	return id, nil
}

// GetUsersById retrieves Users by Id. Returns error if
// Id doesn't exist
func GetUsersById(id int64) (v *Users, err error) {
	o := orm.NewOrm()
	v = &Users{Id: id}
	if err = o.QueryTable(new(Users)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUsers retrieves all Users matches certain condition. Returns empty list if
// no records exist
func GetAllUsers(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Users))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Users
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateUsers updates Users by Id and returns error if
// the record to be updated doesn't exist
func UpdateUsersById(m *request.UserUpdateRequest) (err error) {
	o := orm.NewOrm()
	v := Users{Id: m.Id}

	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {

		v.FirstName = m.FirstName
		v.LastName = m.LastName
		v.Email = m.Email
		v.Mobile = m.Mobile
		v.Role = m.Role
		var num int64
		if num, err = o.Update(&v); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUsers deletes Users by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUsers(id int64) (err error) {
	o := orm.NewOrm()
	v := Users{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Users{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func FindUserByEmail(email string) (*Users, error) {
	o := orm.NewOrm()
	v := &Users{Email: email}

	if err := o.QueryTable(new(Users)).Filter("Email", email).RelatedSel().One(v); err != nil {
		return nil, err
	}
	return v, nil
}

func FindUserByMobile(mobile string) (*Users, error) {
	o := orm.NewOrm()
	v := &Users{Mobile: mobile}

	if err := o.QueryTable(new(Users)).Filter("Mobile", mobile).RelatedSel().One(v); err != nil {
		return nil, err
	}
	return v, nil
}

func IsExistingUser(email string, mobile string) bool {
	v, _ := FindUserByEmail(email)
	if v != nil {
		return true
	}
	u, _ := FindUserByMobile(mobile)
	return u != nil
}

func VerifyUserAndUpdate(email string, otp string) (*Users, error) {
	o := orm.NewOrm()

	tx, err := o.Begin()
	if err != nil {
		return nil, err
	}
	var data Users
	if err = tx.QueryTable("users").Filter("email", email).One(&data); err != nil {
		// Rollback the transaction if there's an error
		tx.Rollback()
		return nil, err
	}

	verify := (otp == data.Otp)

	if !verify {
		tx.Rollback()
		return nil, errors.New("Wrong OTP")
	}

	data.Status = 1

	if _, err := tx.Update(&data); err != nil {
		// Rollback the transaction if there's an error
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &data, nil
}

func SendOtpToUser(email string, otp string) (bool, error) {
	o := orm.NewOrm()

	_, err := o.Raw("CALL updateOtps(?, ?)", email, otp).Exec()
	// tx, err := o.Begin()
	// if err != nil {
	// 	return nil, err
	// }
	// var data Users
	// if err = tx.QueryTable("users").Filter("email", email).One(&data); err != nil {
	// 	// Rollback the transaction if there's an error
	// 	tx.Rollback()
	// 	return nil, err
	// }

	// data.Otp = otp
	// data.UpdatedAt = time.Now().UnixMilli()

	// if _, err := tx.Update(&data); err != nil {
	// 	// Rollback the transaction if there's an error
	// 	tx.Rollback()
	// 	return nil, err
	// }

	// if err := tx.Commit(); err != nil {
	// 	return nil, err
	// }
	if err != nil {
		return false, err
	}

	return true, nil
}

func UserPasswordUpdate(email string, password string, otp string) (*Users, error) {
	o := orm.NewOrm()

	tx, err := o.Begin()
	if err != nil {
		return nil, err
	}
	var data Users
	if err = tx.QueryTable("users").Filter("email", email).One(&data); err != nil {
		// Rollback the transaction if there's an error
		tx.Rollback()
		return nil, err
	}

	verify := (otp == data.Otp)

	if !verify {
		tx.Rollback()
		return nil, errors.New("Wrong OTP")
	}

	data.Password = password
	data.UpdatedAt = time.Now().Unix()

	if _, err := tx.Update(&data); err != nil {
		// Rollback the transaction if there's an error
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &data, nil
}

func UserPasswordReset(email, currentPassword, NewPassword string) (*Users, error) {
	o := orm.NewOrm()

	tx, err := o.Begin()
	if err != nil {
		return nil, err
	}
	var data Users
	if err = tx.QueryTable("users").Filter("email", email).One(&data); err != nil {
		// Rollback the transaction if there's an error
		tx.Rollback()
		return nil, err
	}

	verify, _ := helper.VerifyHashedData(data.Password, currentPassword)

	if !verify {
		tx.Rollback()
		return nil, errors.New("Wrong Password")
	}

	data.Password = NewPassword
	data.UpdatedAt = time.Now().Unix()

	if _, err := tx.Update(&data); err != nil {
		// Rollback the transaction if there's an error
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &data, nil
}

func LastInsertedUser() (*Users, error) {
	o := orm.NewOrm()

	var lastRecord Users
	err := o.QueryTable(new(Users)).OrderBy("-id").Limit(1).One(&lastRecord)
	if err != nil {
		return nil, err
	}
	return &lastRecord, nil
}
