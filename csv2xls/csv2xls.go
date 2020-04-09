package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"encoding/csv"
	"io/ioutil"
	"log"
	"strings"
	"io"
	"os"
)

/**
 *  csv转excel工具
 *	功能: 执行后会把当前目录下所有的*.csv文件转换成*.csv.xls文件
*/

func transFile(csvPath string) {
	//csvPath := "111.csv"
	csvfile,err := os.Open(csvPath)
	if err != nil {
		log.Fatal("Error when open csv file:", err)
	}
	defer csvfile.Close()

	//decoder := mahonia.NewDecoder("gbk")
	//r := csv.NewReader(decoder.NewReader(csvfile))
	r := csv.NewReader(csvfile)

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	file = xlsx.NewFile()
	sheet,_ = file.AddSheet("Sheet1")

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(record)
		row = sheet.AddRow()
		for _,v := range record {
			cell = row.AddCell()
			if strings.Index(v, "\t") != -1 {
				strings.Replace(v, "\t", "", -1)
			}
			//cell.Type(xlsx.CellTypeString)
			cell.Value = v
			//cell.SetString(v)
			//fmt.Println(fmt.Sprintf("%d=>%s",k,v))
		}
	}

	err = file.Save(fmt.Sprintf("%s.xlsx", csvPath))
	if err != nil {
		fmt.Printf(err.Error())
	} else {
		fmt.Println(fmt.Sprintf("%s.xlsx 转换成功", csvPath))
	}
}

func main() {

	dirPth := "."
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		log.Fatal("Error when ReadDir:", err)
	}
	suffix := "CSV"
	//PthSep := string(os.PathSeparator)
	files := make([]string, 0, 10)
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, fi.Name())
		}
	}
	//fmt.Println(files)

	for _,csvPath := range files {
		transFile(csvPath)
	}

	var in string
	fmt.Printf("按回车退出: ")
	fmt.Scanln(&in)
	fmt.Printf("%s\n", in)
}