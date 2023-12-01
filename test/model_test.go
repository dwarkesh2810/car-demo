package test

import (
	"car_demo/conf"
	"car_demo/models"
	"car_demo/request"
	"log"
	"testing"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/lib/pq"
)

// type TUsers struct {
// 	Id        int64  `orm:"auto"`
// 	FirstName string `orm:"size(128)" json:"first_name" form:"first_name"`
// 	LastName  string `orm:"size(128)" json:"last_name" form:"last_name"`
// 	Email     string `orm:"size(128)" json:"email" form:"email"`
// 	Mobile    string `json:"mobile" form:"mobile"`
// 	Password  string `orm:"size(128)" json:"password" form:"password"`
// 	Status    int    `orm:"" json:"status" form:"status"`
// 	Role      string `orm:"size(20)" json:"role" form:"role"`
// 	Otp       string `orm:"size(20)" json:"otp" form:"otp"`
// 	CreatedAt int64
// 	UpdatedAt int64
// 	DeletedAt int64
// }

func Init() {
	conf.LoadEnv("..")
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "user=postgres password=1234 dbname=postgres sslmode=disable")
	orm.RunSyncdb("default", false, true)
}

func TestIsExistingUser(t *testing.T) {
	Init()
	t.Run("Existing User", func(t *testing.T) {
		email := "dwarkesh0007@gmail.com"
		mobile := "12345648"

		result := models.IsExistingUser(email, mobile)
		expected := true
		if result != expected {
			t.Fatal("expected true but got false ")
		}
	})

	t.Run("Non Existing User", func(t *testing.T) {
		email := "dwarkeshbpatel@gmail.com"
		mobile := "123456481212"

		result := models.IsExistingUser(email, mobile)
		expected := false
		if result != expected {
			t.Fatal("expected false but got true ")
		}
	})
}

func TestAddUser(t *testing.T) {
	Init()
	ClearData("users")
	t.Run("Add User", func(t *testing.T) {

		user := UserCreateData()
		id, _ := models.AddUsers(user)

		var fields []string
		var sortby []string
		var order []string
		var query = make(map[string]string)
		var limit int64 = 1
		var offset int64

		users, _ := models.GetAllUsers(query, fields, sortby, order, offset, limit)
		last_id := users[0].(models.Users).Id

		if id != last_id {
			t.Fatal("failed to Add user")
		}

		// models.GetAllUsers(order []string {"desc"})

	})
}

func TestFindUserByEmail(t *testing.T) {
	Init()
	ClearData("users")
	t.Run("Find User By Email", func(t *testing.T) {

		user := UserCreateData()

		models.AddUsers(user)

		users, _ := models.FindUserByEmail(user.Email)

		result := user.Mobile == users.Mobile
		if !result {
			t.Fatalf("function not working as expected ")
		}

	})
}

func TestFindUserByMobile(t *testing.T) {
	Init()
	ClearData("users")
	t.Run("Find User By Mobile", func(t *testing.T) {

		user := UserCreateData()

		models.AddUsers(user)

		users, _ := models.FindUserByMobile(user.Mobile)

		// result := reflect.DeepEqual(user, users)
		result := user.Email == users.Email

		if !result {
			t.Fatalf("function not working as expected ")
		}

	})
}

func TestUpdateUsersById(t *testing.T) {
	Init()
	ClearData("users")
	t.Run("Update User By Id", func(t *testing.T) {

		user := UserCreateData()

		id, _ := models.AddUsers(user)

		uUser := UserUpdateData(id)

		models.UpdateUsersById(uUser)
		a, _ := models.GetUsersById(id)
		log.Print(a)
		log.Print(uUser)
		if a.Email != uUser.Email {
			t.Fatalf("function not working as expected ")
		}

	})
}

func TestSendOTPMOdel(t *testing.T) {

}

func UserUpdateData(id int64) *request.UserUpdateRequest {
	return &request.UserUpdateRequest{
		FirstName: "Dax",
		LastName:  "dexter",
		Email:     "dwarkesh@mail.com",
		Mobile:    "1232131232",
		Role:      "user",
		Id:        id,
		Password:  "123456",
	}

}
