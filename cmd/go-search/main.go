package main

import (
	_ "github.com/lib/pq"
)

func main() {

	// site1, err := site.NewSite("https://habr.com/ru/articles/481556/")

	// if err != nil {
	// 	os.Exit(1)
	// }

	// site2, err := site.NewChildSite("..", *site1)

	// if err != nil {
	// 	os.Exit(1)
	// }

	// cfg_queue, err := postgres.NewConfig(".url_queue.env")

	// if err != nil {
	// 	fmt.Printf("err: %s\n", err)
	// 	os.Exit(1)
	// }

	// db, err := postgres.NewConn(*cfg_queue)

	// if err != nil {
	// 	fmt.Printf("err: %s\n", err)
	// 	os.Exit(1)
	// }

	// rep := tasks.NewRepository(db)

	// // err = rep.AddTask([]site.Site{*site1, *site2})

	// // if err != nil {
	// // 	fmt.Println(err)
	// // 	os.Exit(1)
	// // }

	// res, err := rep.GetLimitTasks(20)

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// ids := make([]int, 0)

	// for _, item := range res {
	// 	ids = append(ids, item.Id)
	// }

	// err = rep.CompleteTasks(ids)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// res, err = rep.GetLimitTasks(20)

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// for _, item := range res {
	// 	fmt.Println(item)
	// }

}
