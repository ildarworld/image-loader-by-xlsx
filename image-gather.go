package main


import (
	"errors"
	"io"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"path"
	"path/filepath"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func downloadFile(URL string, path string) error {
	//Get the response bytes from the url
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
	
	fileName := path + "/" + s[len(s)-1]

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

	xlsx_file_name := path.Join(_path, "ссылки.xlsx")
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
	image_path := path.Join(_path, "Фотографии " + t.Format("2006-01-02 15-04-05"))
    os.Mkdir(image_path, 0755)
	
	for i := 1; i <= len(rows)-1; i++ {
		cell := rows[i][links_column_index]
		fmt.Println(cell)
		downloadFile(cell, image_path)
		links = append(links, cell)
	}

}