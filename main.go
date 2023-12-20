package main

import (
	"RankWillServer/service"
	"log"
	"sync"
	"time"
)

func main() {
	t1 := time.Now().UnixMilli()
	client := service.GenNewClient()
	wg := sync.WaitGroup{}
	cnt := 0
	for i := 0; i < 900; i += 50 {
		go func(x int16) {
			wg.Add(1)
			for k := 0; k < 100; k++ {
				for j := x; j < x+50; j++ {
					pg, _ := service.QueryContestRankByNameAndPage(&client, "weekly-contest-370", j)
					t := time.Now().UnixMilli()
					cnt++
					if (t-t1)%10000 < 2000 {
						time.Sleep(2000 * time.Millisecond)
					}
					pg.UserNum = 1
					log.Printf("%d %v\n", cnt, 1000*float64(cnt)/float64(t-t1))
				}
			}
			wg.Done()
		}(int16(i))
	}
	wg.Wait()
	t2 := time.Now().UnixMilli()
	log.Println(t2 - t1)
}
