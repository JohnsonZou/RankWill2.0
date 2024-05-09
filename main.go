package main

import (
	"RankWillServer/backend"
	"RankWillServer/dao"
	"RankWillServer/web_server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	logfile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logfile)
}
func main() {
	_ = dao.InitDB()

	go backend.Serve()

	go web_server.GinRun()

	//go func() {
	//	c := model.Contest{
	//		TitleSlug: "weekly-contest-396",
	//	}
	//	ctx := backend.InitContext(context.Background())
	//	_ = service.FetchContest(ctx, &c)
	//	_ = service.Predict(ctx, &c)
	//}()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)
	fmt.Println(<-sig)
}
