package prober

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"dtool/utils"
)

type Data struct {
	target string
	dict   map[string]bool
}

var dataset map[string][]string

func retrieve_ip(pool chan string, filename string) {
	cnt := 0
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("cannot open file %s\n", filename)
		return
	}
	fmt.Println("sending msg ...")
	reader := bufio.NewReader(f)
	for {
		s, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		s = strings.Trim(s, "\n")
		pool <- s
		cnt++
		if cnt%10 == 0 {
			fmt.Println(cnt)
		}
	}
	close(pool)
}

func active_probe(n int, addr string) Data {
	target_ip := addr[:len(addr)-3]
	data := Data{target_ip, make(map[string]bool)}
	stop := 0
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	for i := 0; i < n; i++ {
		subdomain := strings.Join([]string{strings.Replace(target_ip, ".", "-", -1), "echo", strconv.Itoa(i), timestamp}, "-")
		rdns_ip, err := utils.SendQuery(addr, subdomain)
		if err == nil {
			data.dict[rdns_ip] = true
		} else {
			stop += 1
		}
		if stop == 3 {
			break
		}
	}
	return data
}

func upstream_prober(ip_pool chan string, data_pool chan Data, wg *sync.WaitGroup) {
	for {
		if s, ok := <-ip_pool; ok {
			addr := s + ":53"
			data := active_probe(20, addr)
			if data.dict != nil {
				data_pool <- data
			}
		} else {
			break
		}
	}
	wg.Done()
}

func create_probers(num int, ip_pool chan string, data_pool chan Data, wg *sync.WaitGroup) {
	for i := 0; i < num; i++ {
		wg.Add(1)
		go upstream_prober(ip_pool, data_pool, wg)
	}
}

func store_data(pool chan Data, wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		var temp []string
		if data, ok := <-pool; ok {
			if len(data.dict) > 0 {
				for rdns := range data.dict {
					temp = append(temp, rdns)
				}
				dataset[data.target] = temp
			}
		} else {
			break
		}
	}
	wg.Done()
}

func Get_upstream_file(filename string, prober_num int) {
	dataset = map[string][]string{}
	ip_pool := make(chan string, 500)
	data_pool := make(chan Data, 200)
	var probe_tasks sync.WaitGroup
	var store_task sync.WaitGroup

	go retrieve_ip(ip_pool, filename)
	create_probers(prober_num, ip_pool, data_pool, &probe_tasks)
	go store_data(data_pool, &store_task)
	probe_tasks.Wait()
	close(data_pool)
	store_task.Wait()
	utils.OutputJSON(dataset)
}

func Get_upstream_ip(ip string) {
	dataset = make(map[string][]string)
	var temp []string
	addr := ip + ":53"
	data := active_probe(10, addr)
	for rdns := range data.dict {
		temp = append(temp, rdns)
	}
	dataset[data.target] = temp
	utils.OutputJSON(dataset)
}
