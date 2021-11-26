package jd_cookie

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/cdle/sillyGirl/core"
	"github.com/gorilla/websocket"
)

var jd_cookie = core.NewBucket("jd_cookie")

var mhome sync.Map

func initLogin() {
	core.BeforeStop = append(core.BeforeStop, func() {
		for {
			running := false
			mhome.Range(func(_, _ interface{}) bool {
				running = true
				return false
			})
			if !running {
				break
			}
			time.Sleep(time.Second)
		}
	})
	go RunServer()
	core.AddCommand("", []core.Function{
		{
			Rules: []string{`raw ^登录$`, `raw ^登陆$`, `raw ^h$`},
			Handle: func(s core.Sender) interface{} {
				if groupCode := jd_cookie.Get("groupCode"); !s.IsAdmin() && groupCode != "" && s.GetChatID() != 0 && !strings.Contains(groupCode, fmt.Sprint(s.GetChatID())) {
					return nil
				}
				if c == nil {
					tip := jd_cookie.Get("tip")
					if tip == "" {
						if s.IsAdmin() {
							s.Reply(jd_cookie.Get("tip", "阿东又不行了。")) //已支持阿东前往了解，https://github.com/rubyangxg/jd-qinglong
							return nil
						} else {
							tip = "暂时无法使用短信登录。"
						}
					}
					s.Reply(tip)
					return nil
				}
				go func() {
					stop := false
					phone := ""

					uid := time.Now().UnixNano()
					cry := make(chan string, 1)
					mhome.Store(uid, cry)
					var deadline = time.Now().Add(time.Second * time.Duration(200))
					var cookie *string
					sendMsg := func(msg string) {
						c.WriteJSON(map[string]interface{}{
							"time":         time.Now().Unix(),
							"self_id":      jd_cookie.GetInt("selfQid"),
							"post_type":    "message",
							"message_type": "private",
							"sub_type":     "friend",
							"message_id":   time.Now().UnixNano(),
							"user_id":      uid,
							"message":      msg,
							"raw_message":  msg,
							"font":         456,
							"sender": map[string]interface{}{
								"nickname": "傻妞",
								"sex":      "female",
								"age":      18,
							},
						})
					}
					if s.GetImType() == "wxmp" {
						cancel := false
						for {
							if phone != "" {
								break
							}
							if cancel {
								break
							}
							s.Await(s, func(s core.Sender) interface{} {
								message := s.GetContent()
								if message == "退出" {
									cancel = true
									return "取消登录"
								}
								if regexp.MustCompile(`^\d{11}$`).FindString(message) == "" {
									return "请输入格式正确的手机号，或者对我说“退出”。"
								}
								phone = message
								return "请输入收到的验证码哦～"
							})
						}
						if cancel {
							return
						}
					}

					defer func() {
						cry <- "stop"
						mhome.Delete(uid)
						if cookie != nil {
							s.SetContent(*cookie)
							core.Senders <- s
						}
						sendMsg("q")
					}()
					go func() {
						for {
							msg := <-cry
							fmt.Println(msg)
							if msg == "stop" {
								break
							}
							msg = strings.Replace(msg, "登陆", "登录", -1)
							if strings.Contains(msg, "不占资源") {
								msg += "\n" + "4.取消"
							}
							if strings.Contains(msg, "无法回复") {
								sendMsg("帮助")
							}
							if strings.Contains(msg, "无法发送验证码") {
								s.Reply("获取验证码失败，请重新输入手机号码")
								continue
							}
							if strings.Contains(msg, "手机号格式有误") {
								s.Reply("手机号格式有误，请重新输入手机号码")
								continue
							}
							{
								res := regexp.MustCompile(`剩余操作时间：(\d+)`).FindStringSubmatch(msg)
								if len(res) > 0 {
									remain := core.Int(res[1])
									deadline = time.Now().Add(time.Second * time.Duration(remain))
								}
							}
							lines := strings.Split(msg, "\n")
							new := []string{}
							for _, line := range lines {
								if !strings.Contains(line, "剩余操作时间") {
									new = append(new, line)
								}
							}
							msg = strings.Join(new, "\n")
							if strings.Contains(msg, "数字编号") {
								sendMsg("1")
								continue
							}
							if strings.Contains(msg, "登录方式") {
								sendMsg("1")
								continue
							}
							if phone != "" && (strings.Contains(msg, "请输入手机号（可输入”退出“结束登录）") || strings.Contains(msg, "请输入11位手机号（可输入”退出“结束登录）")) {
								sendMsg(phone)
								continue
							}
							if strings.Contains(msg, "pt_key") {
								cookie = &msg
								stop = true
								s.SetContent("q")
								core.Senders <- s
							}
							if cookie == nil {
								if strings.Contains(msg, "已点击登录") {
									continue
								}
								s.Reply(msg)
							}
						}
					}()
					sendMsg("h")
					for {
						if stop == true {
							break
						}
						if deadline.Before(time.Now()) {
							stop = true
							s.Reply("登录超时")
							break
						}
						s.Await(s, func(s core.Sender) interface{} {
							msg := s.GetContent()
							if msg == "查询" || strings.Contains(msg, "pt_pin=") {
								s.Continue()
								return nil
							}
							iw := core.Int(msg)
							if msg == "q" || msg == "exit" || msg == "退出" || msg == "10" || msg == "4" || (fmt.Sprint(iw) == msg && iw > 1 && iw < 11) {
								stop = true
								if cookie == nil {
									return "取消登录"
								} else {
									return "登录成功，请加我为好友，定时推送资产消息与任务完成消息给你"
								}
							}
							if phone != "" {
								if regexp.MustCompile(`^\d{6}$`).FindString(msg) == "" {
									return "请输入格式正确的验证码，或者对我说“退出”。"
								} else {
									s.Reply("八九不离十登录成功啦，60秒后对我说“查询”已确认登录成功。")
								}
							}
							sendMsg(s.GetContent())
							return nil
						}, `[\s\S]+`)
					}
				}()
				if s.GetImType() == "wxmp" {
					return "请输入11位手机号："
				}
				return nil
			},
		},
	})
	// if jd_cookie.GetBool("enable_aaron", false) {
	// core.Senders <- &core.Faker{
	// 	Message: "ql cron disable https://github.com/Aaron-lv/sync.git",
	// }
	// core.Senders <- &core.Faker{
	// 	Message: "ql cron disable task Aaron-lv_sync_jd_scripts_jd_city.js",
	// }
	// }
}

var c *websocket.Conn

func RunServer() {
	addr := jd_cookie.Get("adong_addr")
	if addr == "" {
		return
	}
	defer func() {
		time.Sleep(time.Second * 2)
		RunServer()
	}()
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws/event"}
	logs.Info("连接阿东 %s", u.String())
	var err error
	c, _, err = websocket.DefaultDialer.Dial(u.String(), http.Header{
		"X-Self-ID":     {fmt.Sprint(jd_cookie.GetInt("selfQid"))},
		"X-Client-Role": {"Universal"},
	})
	if err != nil {
		logs.Warn("连接阿东错误:", err)
		return
	}
	defer c.Close()
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logs.Info("read:", err)
				return
			}

			type AutoGenerated struct {
				Action string `json:"action"`
				Echo   string `json:"echo"`
				Params struct {
					UserID  interface{} `json:"user_id"`
					Message string      `json:"message"`
				} `json:"params"`
			}
			ag := &AutoGenerated{}
			json.Unmarshal(message, ag)
			if ag.Action == "send_private_msg" {
				UserID, _ := strconv.ParseInt(ag.Params.UserID.(string), 10, 64)
				if cry, ok := mhome.Load(UserID); ok {
					fmt.Println(ag.Params.Message)
					cry.(chan string) <- ag.Params.Message
				}
			}
			logs.Info("recv: %s", message)
		}
	}()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(`{}`))
			if err != nil {
				logs.Info("阿东错误:", err)
				c = nil
				return
			}
		}
	}
}

func decode(encodeed string) string {
	decoded, _ := base64.StdEncoding.DecodeString(encodeed)
	return string(decoded)
}
