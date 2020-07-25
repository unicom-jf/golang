package main

import (
	"bufio"
	"fmt"
	"goinaction/p_work/work"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type fileGetter struct {
	file string
}

//Task ....
/*
func (m *fileGetter) Task(workerID int) {
	//log.Println(workerID, " -- ", m.file)
	//fmt.Println("starting with ", m.file)
	fmt.Println("wget -c " + rootURL + m.file)

}
*/

func (m *fileGetter) Task(workerID int) {
	log.Println(workerID, " -- ", m.file)
	/*
		//defer wg.Done()
		fmt.Println("starting with ", m.file)
		out, err := os.Create(m.file)
		if err != nil {
			fmt.Println(err)
		}
		defer out.Close()
		resp, err := http.Get(rootURL + m.file)
		if err != nil {
			fmt.Println("get ", m.file, " with ", err)
			return
		}
		defer resp.Body.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			fmt.Println("save ", m.file, " with ", err)
			return
		}
	*/
	resp, err := http.Head(rootURL + m.file)

	if err != nil {
		fmt.Println(err)
		return
	}
	//defer resp.Body.Close()
	size := resp.ContentLength
	fmt.Println("ContentLength: ", size)
	resp.Body.Close()
	out, err := os.OpenFile(m.file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	info, err := out.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	done := info.Size()
	fmt.Println("size: ", size, "done: ", done)
	defer out.Close()
	if size == done {
		fmt.Println("already done")
		return
	}
	req, err := http.NewRequest("GET", rootURL+m.file, nil)
	req.Header.Add("Range", "bytes="+strconv.FormatInt(done, 10)+"-"+strconv.FormatInt(size, 10))
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("save WALinuxAgent-2.2.18-1.el7.noarch.rpm", " with ", err)
		return
	}
	fmt.Println(m.file, " -- done")
	//time.Sleep(time.Second)
}

const rootURL = "http://vault.centos.org/7.4.1708/extras/x86_64/Packages/"

var wg sync.WaitGroup
var files []string
var fileNumb int

func main() {

	//fileNumb = 0

	getFiles()
	//wg.Add(fileNumb)
	fmt.Println("waiting...")
	p := work.New(5)
	//var wg sync.WaitGroup
	wg.Add(len(files))
	fmt.Println(len(files))

	for _, name := range files {
		//name = ""
		fn := getFileName(name)
		np := fileGetter{
			file: fn,
		}
		//fmt.Println(np.file)
		p.Run(&np)
		wg.Done()
	}

	wg.Wait()
	p.Shutdown()
	/**
	for _, file := range files {
		//fmt.Println(file)
		fn := getFileName(file)
		go ftpGet(fn)
		//fmt.Println(fn[0:5])
	}
	**/
	//wg.Wait()
	fmt.Println("all done")

}
func getFileName(line string) string {
	//line := "[   ]	ansible-doc-2.3.1.0-3.el7.noarch.rpm	2017-09-07 21:09	279K"
	pos1 := strings.Index(line, "]")
	pos2 := strings.Index(line, ".rpm")
	//fmt.Println("pos: ", pos)
	line = line[pos1+1 : pos2+4]
	return strings.TrimSpace(line)
	//str = strings.Trim(str, " ")
	//fmt.Println(strings.TrimSpace(str))
}
func getFiles() {
	f, _ := os.Open("/Users/jf10/rpms.txt")
	defer f.Close()
	r := bufio.NewScanner(f)
	//i := 0
	//num := 0
	for r.Scan() {
		//files[i] = r.Text()
		//i++
		fileNumb++
		files = append(files, r.Text())
	}
	//fmt.Println(num)
}
func ftpGet(file string) {
	//defer wg.Done()
	fmt.Println("starting with ", file)
	out, err := os.Create(file)
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()
	resp, err := http.Get(rootURL + file)
	if err != nil {
		fmt.Println("get ", file, " with ", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("save ", file, " with ", err)
	}
	fmt.Println(file, " -- done")
}
