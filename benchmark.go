
package main
import (
"flag"
"os"
"fmt"
"log"
"sync"
"time"
"net"
)

func HandleWrite(conn net.Conn, done chan string){
	_, err := conn.Write([]byte("hello world! xxx x xjf \r\n"))
	if err != nil{
		fmt.Println("Error to send message:", err.Error())

	}
	done <- "Sent"
}

func HandleRead(conn net.Conn, done chan string){
	buf := make([]byte,  1024)
	_, err := conn.Read(buf)
	if err != nil{
		fmt.Println("Error to read message:", err.Error())
		//return
	}
	done <- "Read"
}


func Benckmark( host string,   port int,  msg string) error {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:19000",time.Second)
	if err != nil{
		fmt.Println("Error connecting:", err)
		return err
	}

	defer conn.Close()
	done := make(chan string)
    go HandleWrite(conn, done)
	go HandleRead(conn, done)
	<-done //fmt.Println(<-done)
	<-done //fmt.Println(<-done)
	return nil
}



func main() {
	ip := flag.String("Ip","127.0.0.1", "Ip address")
	port :=flag.Int("Port", 19000, "Port")
	//timeout := flag.Int("Timeout", 1000, "Timeout")
	//size := flag.Int("Size", 4096, "size for server to return")
	count := flag.Int("Count", 1000, "Total request to be send")
	max := flag.Int("Max", 1000, "Max request to be send at one time")
	file := flag.String("File","./request.log", "Log file")
	flag.Parse()

	logFile, err := os.OpenFile(*file, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
	defer logFile.Close()
	if err != nil{
		fmt.Println("Error create request log")
	}
	log.SetOutput(logFile)

	var wg sync.WaitGroup
	jobs := make(chan int , *count)
	for i:=0; i < *count ; i++{
		jobs <- i
	}

	fmt.Println("Do benchmark")
	t := time.Now().UnixNano()
	for i:=0; i <*max; i++{
		wg.Add(1)
		go func(){
			defer wg.Done()
			for id := range jobs{
				st := time.Now().UnixNano()
				var result []byte
				//request := make([]byte, *size)
				err := Benckmark(*ip, *port, "hello")
				if err != nil{
					failCost := (time.Now().UnixNano() -st)/ int64(time.Microsecond)
					log.Printf("%d|%d|%d|%d|%s", id, 1, 0, failCost, err.Error())
					return
				}

				cost := (time.Now().UnixNano() -st)/ int64(time.Microsecond)
				//jobid, succ or fail , length, costtime, error msg
				log.Printf("%d|%d|%d|%d", id, 0, len(result), cost)

			}
		}()
	}
	close(jobs)
	wg.Wait()
	totoal := (time.Now().UnixNano() - t) /int64(time.Second)
	fmt.Printf("Benchmark done, total cost:%d s \n", totoal)
}
