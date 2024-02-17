package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/websocket"
)

const TermBufferSize = 8192

const Api_key = "dTAu1iOvOfxQ63BZsYQpDqvyHMjeD8itjZ7GTs"

type termErr struct {
	Cause string `json:"cause"`
}

type newTermReq struct {
	Host     string `query:"host" form:"host" json:"host"`
	Port     int    `query:"port" form:"port" json:"port"`
	Username string `query:"user" form:"user" json:"user"`
	Password string `query:"pwd" form:"pwd" json:"pwd"`
	Rows     int    `query:"rows" form:"rows" json:"rows"`
	Cols     int    `query:"cols" form:"cols" json:"cols"`
}

type setTermWindowSizeReq struct {
	Rows int `query:"rows" form:"rows" json:"rows"`
	Cols int `query:"cols" form:"cols" json:"cols"`
}

type dataCfg struct {
	Key   string `query:"key"`
	Url   string `query:"url"`
	Token string `query:"token"`
}

func listTermHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, termErr{"nothing here!"})
}

func createTermHandler(c echo.Context) error {
	req := new(newTermReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, termErr{err.Error()})
	}
	if req.Host == "" {
		return c.JSON(http.StatusBadRequest, termErr{"Host not provided"})
	}
	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, termErr{"User or password not provided"})
	}
	term, err := termStore.New(TermOption{
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, termErr{err.Error()})
	}
	return c.JSON(http.StatusOK, term)
}

func setTermWindowSizeHandler(c echo.Context) error {
	req := new(setTermWindowSizeReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, termErr{err.Error()})
	}
	if req.Rows == 0 || req.Cols == 0 {
		return c.JSON(http.StatusBadRequest, termErr{"Rows or cols can't be zero"})
	}
	term, err := termStore.Get(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, termErr{err.Error()})
	}
	defer termStore.Put(term.Id)
	err = term.SetWindowSize(req.Rows, req.Cols)
	if err != nil {
		return c.JSON(http.StatusBadRequest, termErr{err.Error()})
	}
	return c.JSON(http.StatusOK, term)
}

func getConfig(c echo.Context) error {
	req := new(dataCfg)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusOK, "404: Not Found")
	}
	if req.Key == "" || req.Key != Api_key {
		return c.String(http.StatusOK, "404: Not Found")
	}
	if req.Url == "" || req.Token == "" {
		return c.String(http.StatusOK, "404: Not Found")
	}
	res, err := fetchData(req.Url, req.Token)
	if err != nil {
		return c.String(http.StatusOK, "404: Not Found")
	}
	return c.String(http.StatusOK, res)
}

func linkTermDataHandler(c echo.Context) error {
	term, err := termStore.Lookup(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, termErr{err.Error()})
	}
	websocket.Handler(func(ws *websocket.Conn) {
		defer func() {
			termStore.Put(term.Id)
			ws.Close()
		}()
		go func() {
			b := [TermBufferSize]byte{}
			for {
				n, err := term.Stdout.Read(b[:])
				if err != nil {
					if !errors.Is(err, io.EOF) {
						websocket.Message.Send(ws, fmt.Sprintf("\nError: %s", err.Error()))
					}
					return
				}
				if n == 0 {
					continue
				}
				websocket.Message.Send(ws, string(b[:n]))
			}
		}()
		go func() {
			b := [TermBufferSize]byte{}
			for {
				n, err := term.Stderr.Read(b[:])
				if err != nil {
					if !errors.Is(err, io.EOF) {
						websocket.Message.Send(ws, fmt.Sprintf("\nError: %s", err.Error()))
					}
					return
				}
				if n == 0 {
					continue
				}
				websocket.Message.Send(ws, string(b[:n]))
			}
		}()
		for {
			b := ""
			err := websocket.Message.Receive(ws, &b)
			if err != nil {
				return
			}
			_, err = term.Stdin.Write([]byte(b))
			if err != nil {
				if !errors.Is(err, io.EOF) {
					websocket.Message.Send(ws, fmt.Sprintf("\nError: %s", err.Error()))
				}
				return
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func runCmd(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer func() {
			ws.Close()
		}()
		for {
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				return
			}
			if msg == "pong" {
				websocket.Message.Send(ws, "ping")
			} else {
				var raw map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &raw); err != nil {
					websocket.Message.Send(ws, "red;;;error processing data!")
					websocket.Message.Send(ws, "blue;;;All job Done!")
					return
				}
				_, ok := raw["key"]
				if !ok || raw["key"] != Api_key {
					websocket.Message.Send(ws, "red;;;error api_key not valid!")
					websocket.Message.Send(ws, "blue;;;All job Done!")
					return
				}
				cmd := raw["cmd"].(string)
				hosts := strings.Split(raw["cfg"].(string), "\n")
				cmd_list := []string{}
				if strings.Contains(cmd, "cus_cmd") {
					split_cmd := strings.Split(cmd, "\n")
					cus_cmd := strings.Split(strings.Replace(split_cmd[0], "cus_cmd = ", "", -1), ", ")
					cmd = strings.Join(split_cmd[1:], "\n")
					for i := 0; i <= len(hosts); i++ {
						cmd_list = append(cmd_list, strings.Replace(cmd, "cus_cmd", cus_cmd[i], -1))
					}
					cmd = strings.Join(cmd_list, "\n")
				} else {
					cmd_list = strings.Split(cmd, "\n")
				}
				var wg sync.WaitGroup
				for i := range hosts {
					wg.Add(1)
					go func(i int) {
						defer wg.Done()
						if hosts[i] == "" {
							return
						}
						acc_host := strings.Split(hosts[i], "||")
						user_pass := strings.Split(acc_host[0], ":")
						o, e, err := execCmd(acc_host[1], user_pass[0], user_pass[1], cmd_list)
						websocket.Message.Send(ws, fmt.Sprintf("blue;;;%s~# %s", acc_host[1], cmd))
						if err != nil {
							websocket.Message.Send(ws, fmt.Sprintf("red;;;%s", err.Error()))
							return
						}
						for x := range o {
							websocket.Message.Send(ws, fmt.Sprintf("green;;;%s", o[x]))
						}
						for y := range e {
							websocket.Message.Send(ws, fmt.Sprintf("yellow;;;%s", e[y]))
						}
					}(i)
				}
				wg.Wait()
				websocket.Message.Send(ws, "blue;;;All job Done!")
			}

		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func execCmd(host string, user string, auth string, cmds []string) ([]string, []string, error) {
	cfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(auth),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}
	conn, err := ssh.Dial("tcp", host, cfg)
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()
	out_msg := []string{}
	err_msg := []string{}
	for _, cmd := range cmds {
		session, err := conn.NewSession()
		if err != nil {
			return nil, nil, err
		}
		defer session.Close()
		var stdoutBuf, stderrBuf bytes.Buffer
		session.Stdout = &stdoutBuf
		session.Stderr = &stderrBuf
		err = session.Run(cmd)
		if err != nil {
			return nil, nil, err
		}
		if stdout := stdoutBuf.String(); stdout != "" {
			out_msg = append(out_msg, stdout)
		}
		if stderr := stderrBuf.String(); stderr != "" {
			err_msg = append(err_msg, stderr)
		}
		session.Close()
	}
	return out_msg, err_msg, nil
}

func fetchData(url string, token string) (string, error) {
	var client = &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		b, err := io.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	return "", errors.New("error fetch data!")
}
