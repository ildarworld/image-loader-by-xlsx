package main

import (
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
	if strings.HasPrefix(URL, "http") != true {
		URL = "https://" + URL
	}

	client := &http.Client{}

	// Create the request
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}

	// Set headers to emulate a browser
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
	)

	// Perform the request
	fmt.Println("File download started " + URL)
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("received non 200 response code: %s", response.Body)
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
	// executable, _ := os.Executable()
	// _path := filepath.Dir(executable)
	_path, _ := os.Getwd()
	// fmt.Println("Current Directory:", _path)

	pause := 0

	fmt.Print("Введите длительность паузы в секундах: ")
	_, err := fmt.Scanf("%d", &pause)
	fmt.Println("Пауза устанволена:")
	fmt.Println(pause)

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

	count := 1
	PAUSE_PER := 100

	for i := 1; i <= len(rows)-1; i++ {
		cell_list := strings.Split(rows[i][links_column_index], ";")
		// Dowmload all images even links are separetaed by comma
		for j := 0; j < len(cell_list); j++ {
			err = downloadFile(cell_list[j], image_path)
			if err != nil {
				errors = append(errors, cell_list[j]+"\t"+err.Error())
				fmt.Println(err)
			}
			fmt.Println(fmt.Sprintf("Обработано фотографий %d", count))
			links = append(links, cell_list[j])
			if count > PAUSE_PER {
				PAUSE_PER += 100
				fmt.Printf("Пауза на %d секунд\n", pause)
				time.Sleep(time.Duration(pause) * time.Second)
			}
			count++
		}
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
