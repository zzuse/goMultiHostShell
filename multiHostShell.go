package main

import (
	//"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

var (
	usrConfigure string
	shellCmd     string
)

type Config struct {
	Hostdetails []HostsType
}

type HostsType struct {
	Host string
	User string
	Pass string
}

func init() {
	flag.StringVar(&shellCmd, "c", "", "To perform a shell command.")
	flag.StringVar(&usrConfigure, "u", "hosts.json", "To configure host/user/pass command.")
}

func Usage() {
	fmt.Printf(`Usage of multissh:
    -c string
    To perform a shell command on all the blade
    -u string
	To read a configure file
	eg: multissh -c "ls" -u "./hosts.json"
    Be careful use this for rm command or something like that.
    `)
	fmt.Println("Don't use this do harmful things")
	os.Exit(1)
}

func main() {
	flag.Parse()
	if os.Args == nil || shellCmd == "" || usrConfigure == "" {
		Usage()
	}

	content, err := ioutil.ReadFile(usrConfigure)
	if err != nil {
		fmt.Print("Error:", err)
	}
	content2 := string(content)

	var conf Config
	dec := json.NewDecoder(strings.NewReader(content2))
	if err := dec.Decode(&conf); err == io.EOF {
		log.Fatal(err)
	} else if err != nil {
		log.Fatal(err)
	}

	idx := conf.Hostdetails
	if idx == nil {
		log.Fatal("no host in json")
	} else {
		response := make(chan string)
		log.Println("json have ", len(idx), " host in ", conf.Hostdetails)
		//fmt.Printf("idx len: %d\n", len(idx))
		for _, n := range idx {
			log.Println("HOSTs", n.Host)
			log.Println("USERs", n.User)

			go dial(n.Host, n.User, n.Pass, 22, 1<<15, shellCmd, response)

		}
		for j := 0; j < len(idx); j++ {
			select {
			case res := <-response:
				fmt.Println(res)
			}
		}
		close(response)
	}
}

func dial(HOST string, USER string, PASS string, PORT int, SIZE int, shellCmd string, res chan string) {
	var auths []ssh.AuthMethod
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
	}
	if PASS != "" {
		auths = append(auths, ssh.Password(PASS))
	}
	config := ssh.ClientConfig{
		User: USER,
		Auth: auths,
	}
	addr := fmt.Sprintf("%s:%d", HOST, PORT)
	conn, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		log.Fatalf("unable to connect to [%s]: %v", addr, err)
	}
	defer conn.Close()

	// Create a session
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("unable to create session: %s", err)
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	running := true
	for running {
		log.Println("------------enter", HOST, "-------------------------")
		log.Println(shellCmd)
		//session.Shell()
		//session.Wait()
		//b, err := session.Output(shellCmd)

		err := session.Run(shellCmd)
		if err != nil {
			log.Fatalf("failed to execute: %s", err)
		}
		shellCmd = string("exit")
		log.Println("HOST: ", shellCmd, HOST)
		log.Println("-------------", shellCmd, HOST, "-------------------------")
		if shellCmd == "exit" {
			running = false
		}
	}
	res <- "done" + HOST
}
