package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func downloadFile(URL string, save_path string) error {
	//Get the response bytes from the url
	fmt.Println("File download started " + URL)
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	s := strings.Split(URL, "/")
	fmt.Println(s)

	fileName := path.Join(save_path, s[len(s)-1])
	fmt.Println(fileName)

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	fmt.Println("File saved into " + fileName)
	return nil
}

func index(slice []string, item string) int {
	for i, _ := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

func main() {
	const URL_HEADER = "Ссылки"
	exectable, _ := os.Executable()
	_path := filepath.Dir(exectable)

	xlsx_file_path := _path

	if filepath.Base(_path) == "MacOS" {
		xlsx_file_path = filepath.Dir(filepath.Dir(filepath.Dir(_path)))
	}

	xlsx_file_name := path.Join(xlsx_file_path, "ссылки.xlsx")

	f, err := excelize.OpenFile(xlsx_file_name)
	if err != nil {
		fmt.Println(err)
		return
	}

	var headers []string
	firstSheet := f.WorkBook.Sheets.Sheet[0].Name
	rows, err := f.GetRows(firstSheet)
	for _, colCell := range rows[0] {
		headers = append(headers, colCell)
	}
	links_column_index := index(headers, URL_HEADER)

	var links []string

	t := time.Now()
	fmt.Println(xlsx_file_path)
	image_path := path.Join(xlsx_file_path, "Фотографии "+t.Format("2006-01-02 15-04-05"))
	os.Mkdir(image_path, 0755)
	fmt.Println(image_path)

	var errors []string

	for i := 1; i <= len(rows)-1; i++ {
		cell := rows[i][links_column_index]
		err = downloadFile(cell, image_path)
		if err != nil {
			errors = append(errors, cell+"\t"+err.Error())
			fmt.Println(err)
		}
		links = append(links, cell)
	}
	fmt.Println("")
	if len(errors) > 0 {
		fmt.Println("Ошибки ниже:")
		for _, err := range errors {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Ошибок нет")
	}
}
