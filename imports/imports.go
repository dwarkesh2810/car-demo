package imports

import (
	"car_demo/helper"
	"car_demo/models"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

func ReadCSVFile(filePath string) ([]map[string]interface{}, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1

	var allRows []map[string]interface{}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	columnHeaders := records[0]

	for _, dataRow := range records[1:] {
		rowData := make(map[string]interface{})

		for index, value := range dataRow {
			rowData[columnHeaders[index]] = value
		}

		allRows = append(allRows, rowData)
	}

	return allRows, nil
}

func ImportFile(filePath string) {
	o := orm.NewOrm()
	allRows, err := ReadCSVFile(filePath)
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
	}
	a, _ := helper.HashData("123456")
	tx, _ := o.Begin()

	for _, j := range allRows {
		user := models.Users{
			FirstName: j["FirstName"].(string),
			LastName:  j["LastName"].(string),
			Mobile:    j["Mobile"].(string),
			Email:     j["Email"].(string),
			Status:    0,
			Role:      "user",
			Password:  a,
			CreatedAt: time.Now().UnixMilli(),
		}
		_, err := o.Insert(&user)
		if err != nil {
			tx.Rollback() // Rollback the transaction if any error occurs
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		return
	}
}
