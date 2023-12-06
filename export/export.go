package export

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/beego/beego/v2/client/orm"
	"github.com/jung-kurt/gofpdf"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func ExportToCSV(model interface{}, columnsToInclude []string) error {
	o := orm.NewOrm()

	modelType := reflect.TypeOf(model)

	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	rowsSlice := reflect.New(reflect.SliceOf(modelType)).Interface()

	_, err := o.QueryTable(model).All(rowsSlice)
	if err != nil {
		return err
	}

	rows := reflect.ValueOf(rowsSlice).Elem()

	file, err := os.Create("output.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	var columnNames []string
	var columnIndexMap = make(map[string]int)

	// Get column headers and map column indices
	for i := 0; i < modelType.NumField(); i++ {
		fieldName := modelType.Field(i).Name

		if len(columnsToInclude) == 0 || contains(columnsToInclude, fieldName) {
			columnNames = append(columnNames, fieldName)
			columnIndexMap[fieldName] = i
		}
	}

	if err := csvWriter.Write(columnNames); err != nil {
		return err
	}

	for i := 0; i < rows.Len(); i++ {
		var data []string
		for _, colName := range columnNames {
			colIndex, ok := columnIndexMap[colName]
			if !ok {
				continue
			}
			data = append(data, fmt.Sprintf("%v", rows.Index(i).Field(colIndex).Interface()))
		}
		if err := csvWriter.Write(data); err != nil {
			return err
		}
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ExportToExcel(model interface{}, columnsToInclude []string) error {
	o := orm.NewOrm()

	modelType := reflect.TypeOf(model)

	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	rowsSlice := reflect.New(reflect.SliceOf(modelType)).Interface()
	_, err := o.QueryTable(model).All(rowsSlice)
	if err != nil {
		return err
	}

	rows := reflect.ValueOf(rowsSlice).Elem()

	f := excelize.NewFile()
	sheetName := "Sheet1"

	index := f.NewSheet(sheetName)
	if index == 0 {
		f.SetSheetName("Sheet1", sheetName)
	}

	var columnNames []string
	var columnIndexMap = make(map[string]int)

	// Get column headers and map column indices
	for i := 0; i < modelType.NumField(); i++ {
		fieldName := modelType.Field(i).Name
		if len(columnsToInclude) == 0 || contains(columnsToInclude, fieldName) {
			columnNames = append(columnNames, fieldName)
			columnIndexMap[fieldName] = i
		}
	}

	// Write column headers to Excel sheet
	for i, columnName := range columnNames {
		colChar := string('A' + i)
		cell := colChar + "1"
		f.SetCellValue(sheetName, cell, columnName)
	}

	// Write data rows to Excel sheet
	for i := 0; i < rows.Len(); i++ {
		for j, colName := range columnNames {
			colIndex, ok := columnIndexMap[colName]
			if !ok {
				continue
			}
			colChar := string('A' + j)
			cell := colChar + strconv.Itoa(i+2) // i+2 to start from row 2 (Excel rows are 1-indexed)
			f.SetCellValue(sheetName, cell, rows.Index(i).Field(colIndex).Interface())
		}
	}

	err = f.SaveAs("output.xlsx")
	if err != nil {
		return err
	}

	return nil
}

func CsvToPdf(filePath string) {

	csvFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 8) // Set font: Arial, regular, size 8

	// Calculate the maximum widths for each column
	var maxColWidths []float64
	for _, row := range rows {
		for i, col := range row {
			width := GetAdjustedWidth(pdf, col)
			if i >= len(maxColWidths) {
				maxColWidths = append(maxColWidths, width)
			} else if width > maxColWidths[i] {
				maxColWidths[i] = width
			}
		}
	}

	// Print the data and format the cells based on the maximum widths
	for _, row := range rows {
		for i, col := range row {
			pdf.CellFormat(maxColWidths[i], 8.0, col, "1", 0, "", false, 0, "")
		}
		pdf.Ln(-1)
	}

	outputPath := "output.pdf" // Replace with desired output file path
	err = pdf.OutputFileAndClose(outputPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PDF created:", outputPath)
}

func GetAdjustedWidth(pdf *gofpdf.Fpdf, text string) float64 {
	textWidth := pdf.GetStringWidth(text) + 10 // Adding a little extra padding
	return textWidth
}

func DbToPdf(model interface{}, columnNames []string) error {
	// Connect to PostgreSQL using Beego
	o := orm.NewOrm()

	var rows []orm.Params
	_, err := o.QueryTable(model).Values(&rows, columnNames...)
	if err != nil {
		return err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 8) // Set font: Arial, regular, size 8

	// Calculate the maximum widths for each column
	var maxColWidths []float64
	for _, row := range rows {
		for i, col := range columnNames {
			width := GetAdjustedWidth(pdf, fmt.Sprintf("%v", row[col]))
			if i >= len(maxColWidths) {
				maxColWidths = append(maxColWidths, width)
			} else if width > maxColWidths[i] {
				maxColWidths[i] = width
			}
		}
	}

	// Print the data and format the cells based on the maximum widths
	for _, row := range rows {
		for i, col := range columnNames {
			pdf.CellFormat(maxColWidths[i], 8.0, fmt.Sprintf("%v", row[col]), "1", 0, "", false, 0, "")
		}
		pdf.Ln(-1)
	}

	outputPath := "output.pdf" // Replace with desired output file path
	err = pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return err
	}

	fmt.Println("PDF created:", outputPath)
	return nil
}
