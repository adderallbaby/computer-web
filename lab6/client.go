package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

func main() {
	var (
		username, password, hostname, port string
	)
	fmt.Print("ftp-host: ")
	fmt.Scan(&hostname)
	fmt.Print("port: ")
	fmt.Scan(&port)
	fmt.Print("login: ")
	fmt.Scan(&username)
	fmt.Print("password: ")
	fmt.Scan(&password)

	c, err := ftp.Dial(hostname+":"+port, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = c.Login(username, password)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		coms := scanner.Text()
		cmd := strings.Split(coms, " ")
		switch cmd[0] {
		case "quit": // заканчиваем работу
			if err := c.Quit(); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		case "stor": // закидывает файл из данной папки на сервер
			f, err := os.Open(cmd[1])
			if err != nil {
				panic(err)
			}
			defer f.Close()
			name := findName(cmd[1])
			fmt.Println(name)
			err = c.Stor(name, f)
			if err != nil {
				panic(err)
			}

		case "retr": // крадем файл с сервера
			r, err := c.Retr(cmd[1])
			if err != nil {
				panic(err)
			}
			defer r.Close()
			buf, err := io.ReadAll(r)
			file, err := os.Create(cmd[1])
			if err != nil {
				fmt.Println("Unable to create file:", err)
			}
			defer file.Close()
			file.WriteString(string(buf))
		case "mdir": // создаем папку
			err := c.MakeDir(cmd[1])
			if err != nil {
				panic(err)
			}
		case "ddir": // удаляем файл или папку
			var err error
			if strings.Index(cmd[1], ".") == -1 {
				err = c.RemoveDir(cmd[1])
			} else {
				err = c.Delete(cmd[1])
			}
			if err != nil {
				panic(err)
			}
		case "list": // вывод содержимого всей папки или папки в папке
			if len(cmd) == 1 {
				lisst(0, "../", c)
			} else {
				lisst(0, cmd[1], c)
			}
		case "news": // запись новостей на сервер
			ti := time.Now().Format(time.RFC1123Z)
			t, _ := time.Parse(time.RFC1123Z, ti)
			text := parse(c)
			if text != "" {
				file, err := os.Create("test.txt")
				if err != nil {
					panic(err)
				}
				defer file.Close()
				file.WriteString(text)

				f, err := os.Open("test.txt")
				if err != nil {
					panic(err)
				}
				defer f.Close()
				err = c.Stor("Куйвашев Дмитрий"+" "+t.String()[:10]+" "+t.String()[11:13]+";"+t.String()[14:16]+";"+t.String()[17:19]+".txt", f)
				if err != nil {
					panic(err)
				}
				err = os.Remove("test.txt")
				if err != nil {
					panic(err)
				}
			}
		default:
			fmt.Println("Unknown command, try again")
		}
	}

}

func findName(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}

func lisst(iter int, folder string, c *ftp.ServerConn) {
	var r []*ftp.Entry
	var err error
	r, err = c.List(folder)
	if err != nil {
		panic(err)
	}
	for _, elem := range r {
		if elem.Name != ".DS_Store" {
			var tab string
			for i := 0; i < iter; i++ {
				tab += "  "
			}
			fmt.Println(tab + elem.Name)
			if strings.Index(elem.Name, ".") == -1 {
				lisst(iter+1, folder+"/"+elem.Name, c)
			}
		}
	}
}

func check(folder string, c *ftp.ServerConn) string {
	var res string
	var s []*ftp.Entry
	var err error
	if folder == "" {
		s, err = c.List("../")
	} else {
		s, err = c.List(folder)
	}
	if len(s) == 0 {
		return ""
	}
	if err != nil {
		panic(err)
	}
	for _, elem := range s {
		if strings.Index(elem.Name, ".") == -1 {
			if folder == "" {
				res += check(elem.Name, c)
			} else {
				res += check(folder+"/"+elem.Name, c)
			}
		} else {
			if strings.Index(elem.Name, "Булкин") != -1 {
				var r *ftp.Response
				if folder != "" {
					r, err = c.Retr(folder + "/" + elem.Name)
				} else {
					r, err = c.Retr(folder + elem.Name)
				}
				if err != nil {
					panic(err)
				}
				buf, err := io.ReadAll(r)
				if err != nil {
					panic(err)
				}
				r.Close()
				res += string(buf)
			}
		}
	}
	return res
}

func parse(c *ftp.ServerConn) string {
	exist := check("", c)
	var file string
	res, _ := http.Get("http://www.msk-times.ru/feed.php") // сюда свою ссылку
	content, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	feedData := string(content)

	for i := 0; i < 30; i++ {
		feedData = feedData[strings.Index(feedData, "<title>")+7:]
		title := feedData[:strings.Index(feedData, "</title>")]
		pubday := feedData[strings.Index(feedData, "<pubDate>")+14 : strings.Index(feedData, "<pubDate>")+16]
		pubTime := feedData[strings.Index(feedData, "<pubDate>")+26 : strings.Index(feedData, "<pubDate>")+34]
		pubMonth := "-09"
		pubYear := "2022"
		pubDate := pubYear + pubMonth + "-" + pubday
		description := feedData[strings.Index(feedData, "<![CDATA[")+9 : strings.Index(feedData, "]]>")+3]
		description = strings.ReplaceAll(description, "]]>", "     ")
		if strings.Index(exist, title) == -1 {
			file += title + " " + pubDate + " " + pubTime + "\n" + description + "\n\n"
		} else {
			fmt.Println(title + " already exists")
		}
	}
	return file
}
