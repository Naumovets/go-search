package tasks

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestAddTask(t *testing.T) {
	// db, mock, err := sqlmock.Newx()

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer db.Close()

	// type args struct {
	// 	tasks  []entities.Task
	// 	unique int
	// }

	// type mockBehavour func(args args, id int)

	// testTable := []struct {
	// 	name         string
	// 	mockBehavour mockBehavour
	// 	args         args
	// 	flag         bool
	// }{
	// 	{
	// 		name: "2-youtube-1-vk",
	// 		args: args{
	// 			tasks: []entities.Task{
	// 				{
	// 					Id:      1,
	// 					SiteURL: "https://youtube.com/",
	// 					Status:  0,
	// 				},
	// 				{
	// 					Id:      2,
	// 					SiteURL: "https://youtube.com/",
	// 					Status:  0,
	// 				},
	// 				{

	// 					Id:      3,
	// 					SiteURL: "https://vk.com/",
	// 					Status:  0,
	// 				},
	// 			},
	// 			unique: 2,
	// 		},
	// 		mockBehavour: func(args args, id int) {
	// 			mock.ExpectQuery("INSERT INTO task (url)").WithArgs(args.tasks)
	// 		},
	// 	},
	// }
}
