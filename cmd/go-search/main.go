package main

import (
	"fmt"
	"os"

	"github.com/Naumovets/go-search/internal/site"
)

func main() {

	// parent, err := parser.NewSite("https://tproger.ru/translations/data-structure-time-complexity-in-python")
	site, err := site.NewSite("https://tproger.ru/translations/data-structure-time-complexity-in-python")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	text, sites, err := site.GetText()

	if err != nil {
		fmt.Println("err:", err)
		os.Exit(1)
	}

	fmt.Printf("text: %s\n", text)

	for _, site := range sites {
		url, err := site.CompleteURL()

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(url)
		}
	}

	// child, err := parser.NewChildSite("/tag/java/", *parent)

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// res, err := child.CompleteURL()

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(res)
	// fmt.Print(parent.BasedURL)

	// text, links, err := parser.GetText("https://tilda.ru/ru/")

	// if err != nil {
	// 	os.Exit(1)
	// }

	// // Создаем файл с именем "result.txt"
	// file, err := os.Create("result.txt")
	// if err != nil {
	// 	fmt.Println("Ошибка при создании файла:", err)
	// 	return
	// }
	// defer file.Close()

	// // Записываем результат в файл
	// _, err = file.WriteString(text)
	// if err != nil {
	// 	fmt.Println("Ошибка при записи в файл:", err)
	// 	return
	// }

	// for _, item := range links {
	// 	fmt.Println(item.Url)
	// }
}
