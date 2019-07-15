
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
	fmt.Println("Benckmark")
	conn, err := net.DialTimeout("tcp", "127.0.0.1:19000",time.Second)
	if err != nil{
		fmt.Println("Error connecting:", err)
		return err
	}

	defer conn.Close()
	done := make(chan string)
    go HandleWrite(conn, done)
	go HandleRead(conn, done)
	fmt.Println(<-done)
	fmt.Println(<-done)
	return nil
}



func main() {
	ip := flag.String("Ip","127.0.0.1", "Ip address")
	port :=flag.Int("Port", 19000, "Port")
	//timeout := flag.Int("Timeout", 1000, "Timeout")
	//size := flag.Int("Size", 4096, "size for server to return")
	n := flag.Int("n", 1000, "Number of requests to perform")
	c := flag.Int("c", 1000, "Number of multiple requests to make at a time")
	file := flag.String("File","./request.log", "Log file")
	flag.Parse()

	fmt.Println("requests=",*n, " concurrency=", *c)
	logFile, err := os.OpenFile(*file, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
	defer logFile.Close()
	if err != nil{
		fmt.Println("Error create request log")
	}
	log.SetOutput(logFile)

	var wg sync.WaitGroup
	jobs := make(chan int , *n)
	for i:=0; i < *n ; i++{
		jobs <- i
	}

	t := time.Now().UnixNano()
	for i:=0; i <*c; i++{
		wg.Add(1)
		fmt.Println("i:", i)
		go func(){
			fmt.Println("jobs.size()", len(jobs))
			for id :=range jobs{

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
			wg.Done()

		}()
	}
	close(jobs) //关闭任务channel
	wg.Wait()
	totoal := (time.Now().UnixNano() - t) /int64(time.Second)
	fmt.Printf("Benchmark done, total cost:%d s \n", totoal)
}
