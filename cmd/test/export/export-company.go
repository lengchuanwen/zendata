package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/cmd/test/import/comm"
	"github.com/easysoft/zendata/cmd/test/import/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"path/filepath"
	"reflect"
)

func main() {
	filePath := "data/company/v2.xlsx"
	sheetName := "company"

	fileUtils.MkDirIfNeeded(filepath.Dir(filePath))

	db := comm.GetDB()
	db.AutoMigrate(
		&model.DataCompany{},
	)

	pos := make([]model.DataCompany, 0)
	db.Find(&pos)

	f := excelize.NewFile()
	index := f.NewSheet(sheetName)
	f.SetActiveSheet(index)

	sheet1 := f.GetSheetName(0)
	f.DeleteSheet(sheet1)

	var infos []model.TableInfo
	db.Raw("desc " + model.DataCompany{}.TableName()).Scan(&infos)

	excelColNameArr, excelColNameHeader := comm.GetExcelColsByTableDef(infos)
	fieldNames := comm.GetStructFields(model.DataCompany{})

	// gen headers
	for index, name := range excelColNameHeader {
		excelColName := excelColNameArr[index]
		excelColId := fmt.Sprintf("%s%d", excelColName, 1)

		f.SetCellValue(sheetName, excelColId, name)
	}

	// gen rows
	for rowIndex, po := range pos {
		for fieldIndex, fieldName := range fieldNames {
			val := ""

			if fieldName == "Id" {
				val = fmt.Sprintf("%d", reflect.ValueOf(po).FieldByName(fieldName).Uint())
			} else {
				val = reflect.ValueOf(po).FieldByName(fieldName).String()
			}

			excelColName := excelColNameArr[fieldIndex]
			excelColId := fmt.Sprintf("%s%d", excelColName, rowIndex+2)

			f.SetCellValue(sheetName, excelColId, val)
		}
	}

	f.SaveAs(filePath)
}
