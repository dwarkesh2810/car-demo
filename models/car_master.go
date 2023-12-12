package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/beego/beego/v2/client/orm"
)

type Car_master struct {
	Id        int64   `orm:"auto"`
	UserId    int64   `json:"user_id"`
	CarName   string  `orm:"size(50)" form:"car_name" json:"car_name"`
	CarImage  string  `orm:"size(200)" json:"car_image" `
	Make      string  `orm:"size(20)" form:"make" json:"make"`
	Model     string  `orm:"size(20)" form:"model" json:"model"`
	CarType   CarType `orm:"size(15)" form:"car_type" json:"car_type"`
	CreatedAt int
	UpdatedAt int
}

type CarType string

const (
	Sedan     CarType = "sedan"
	Hatchback CarType = "hatchback"
	SUV       CarType = "SUV"
)

func init() {
	orm.RegisterModel(new(Car_master))
}

// AddCar_master insert a new Car_master into database and returns
// last inserted Id on success.
func AddCar_master(m *Car_master) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCar_masterById retrieves Car_master by Id. Returns error if
// Id doesn't exist
func GetCar_masterById(id int64) (v *Car_master, err error) {
	o := orm.NewOrm()
	v = &Car_master{Id: id}
	if err = o.QueryTable(new(Car_master)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCar_master retrieves all Car_master matches certain condition. Returns empty list if
// no records exist
func GetAllCar_master(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Car_master))
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

	var l []Car_master
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

// UpdateCar_master updates Car_master by Id and returns error if
// the record to be updated doesn't exist
func UpdateCar_masterById(m *Car_master) (err error) {
	o := orm.NewOrm()
	v := Car_master{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCar_master deletes Car_master by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCar_master(id int64) (err error) {
	o := orm.NewOrm()
	v := Car_master{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Car_master{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func IsValidCarType(car CarType) bool {
	switch car {
	case Sedan, Hatchback, SUV:
		return true
	default:
		return false
	}
}
