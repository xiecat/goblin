package plugin

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"goblin/internal/plugin/dump"
	"goblin/internal/plugin/inject"
	"goblin/internal/plugin/replace"
	"goblin/pkg/utils"

	"gopkg.in/yaml.v2"
	log "unknwon.dev/clog/v2"
)

func (base *BasePlugin) GenPlugin(pluginFile string) error {
	if utils.FileExist(pluginFile) {
		return errors.New("File exist please check: " + pluginFile)
	}
	data, _ := yaml.Marshal(base)
	log.Warn("Will Write Plugin file: %s\n", pluginFile)
	err := ioutil.WriteFile(pluginFile, data, 0644) // 写入
	if err != nil {
		log.Fatal("%s can't  write Plugin file\n", pluginFile)
		return err
	}
	return nil

}

func GenDefaultPlugin(pluginFile string) error {
	base := &BasePlugin{
		Name:        path.Base(pluginFile),
		Version:     "0.0.1",
		Description: "this is a description",
		WriteDate:   time.Now().Format("2006-01-02"),
		Author:      "goblin",
		Rule: []*Rule{
			{
				URL:   "/login.php",
				Match: "word",
				Replace: []*replace.Replace{
					{
						Request: &replace.Request{
							Method: []string{"GET", "POST"},
							Header: map[string]string{
								"X-Forwarded-For": "127.0.0.1",
								"X-Real-IP":       "127.0.0.1",
							},
						},
						Response: &replace.Response{
							Status: 200,
							Header: map[string]string{
								"GoblinServer": Version,
							},
							Cookie: &replace.Cookie{
								Secure:   false,
								HttpOnly: false,
								SameSite: http.SameSiteNoneMode,
							},
							Body: &replace.Body{
								ReplaceStr: []*replace.ReplaceStr{
									{
										Old:   "Hello World",
										New:   "Hello goblin",
										Count: -1,
									},
								},
							},
						},
					},
				},
			},
			{
				URL:   "/dump",
				Match: "word",
				Dump: []*dump.Dump{
					{
						Request: &dump.Request{
							Method: []string{"POST"},
						},
						Response: &dump.Response{
							Status: 200,
						},
					},
				},
			},
			{
				URL:   "/test.js",
				Match: "word",
				InjectJs: &inject.InjectJs{
					EvilJs: "aaa.js",
				},
			},
		},
	}
	return base.GenPlugin(pluginFile)
}

func LoadPlugin(pathFile string) (*BasePlugin, error) {
	configFile, err := os.Open(pathFile)
	if err != nil {
		log.Fatal("readPlugin(%s) os.Open failed: %v", pathFile, err)
	}
	defer configFile.Close()

	content, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal("readPlugin(%s) ioutil.ReadAll failed: %v", pathFile, err)
		return nil, err
	}
	plugin := &BasePlugin{}

	err = yaml.Unmarshal(content, plugin)
	if err != nil {
		return nil, fmt.Errorf("[%s] yaml format err:%s ", pathFile, err.Error())
	}
	return plugin, nil
}

func (base *BasePlugin) SetInitConfig(PluginVar *TmpVariable) {
	for _, rule := range base.Rule {
		// InjectJs 检查
		if rule.InjectJs != nil {
			if rule.InjectJs.EvilJs != "" {
				if utils.FileExist(rule.InjectJs.EvilJs) {
					fname := path.Base(rule.InjectJs.EvilJs)
					b, err := ioutil.ReadFile(rule.InjectJs.EvilJs) // just pass the file name
					if err != nil {
						log.Fatal("%s", err.Error())
					}
					fname = "/" + PluginVar.Static + "/" + fname
					if _, ok := StaticFiles[fname]; ok {
						log.Fatal("fname: %s duplicate files pleace check", fname)
					}
					rule.InjectJs.EvilJs = fname
					StaticFiles[fname] = b
				}
				// 替换并且检查替换的变量
				tmpl, err := template.New("test").Parse(rule.InjectJs.EvilJs)
				if err != nil {
					log.Fatal(err.Error())
				}
				var tpl bytes.Buffer
				err = tmpl.Execute(&tpl, PluginVar)
				if err != nil {
					log.Fatal(err.Error())
				}
				log.Info("[plugin] Parse:rule.InjectJs.EvilJs: %s ==> %s ", rule.InjectJs.EvilJs, tpl.String())
				rule.InjectJs.EvilJs = tpl.String()
			}
		}
		//

		// Replace 检查
		if rule.Replace != nil {
			for _, rp := range rule.Replace {
				if rp.Response != nil {
					if rp.Response.Location != "" {
						// 替换并且检查替换的变量
						tplStr, err := utils.TempStr(rp.Response.Location, PluginVar)
						if err != nil {
							log.Fatal(err.Error())
						}
						log.Info("[plugin] Parse:Response.Header.Location: %s ==> %s ", rp.Response.Location, tplStr)
						rp.Response.Location = tplStr
					}

					if rp.Response.Header != nil {
						if location, ok := rp.Response.Header["Location"]; ok {
							// 替换并且检查替换的变量
							tplStr, err := utils.TempStr(location, PluginVar)
							if err != nil {
								log.Fatal(err.Error())
							}
							log.Info("[plugin] Parse:Response.Header.Location: %s ==> %s ", location, tplStr)
							rp.Response.Header["Location"] = tplStr
						}
					}
					// BodyFile
					if rp.Response.Body != nil {
						if rp.Response.Body.File != "" {
							fname := rp.Response.Body.File
							if _, ok := replace.BodyFiles[fname]; ok {
								log.Trace("[plugin] Parse:Response.Body.File fname: %s duplicate files", fname)
							} else {
								b, err := ioutil.ReadFile(fname) // just pass the file name
								if err != nil {
									log.Fatal("%s", err.Error())
								}
								log.Trace("[plugin] Parse:Response.Body.File  load %s", fname)
								// todo 二进制数据可能会有问题
								tplStr, err := utils.TempStr(string(b), PluginVar)
								if err != nil {
									log.Fatal(err.Error())
								}
								replace.BodyFiles[fname] = []byte(tplStr)
							}
						}
						// replace 可能为nil
						if rp.Response.Body.ReplaceStr != nil {
							// str 处理
							rpStr := rp.Response.Body.ReplaceStr
							if len(rpStr) > 0 {
								for key, str := range rpStr {
									tplStr, err := utils.TempStr(str.New, PluginVar)
									if err != nil {
										log.Fatal(err.Error())
									}
									log.Info("[plugin] Parse:rp.Response.Body.ReplaceStr: %s ==> %s ", str.New, tplStr)
									rp.Response.Body.ReplaceStr[key].New = tplStr
								}

							}
						}

						if rp.Response.Body.Append != "" {
							appstr := rp.Response.Body.Append
							// 替换并且检查替换的变量
							tplStr, err := utils.TempStr(appstr, PluginVar)
							if err != nil {
								log.Fatal(err.Error())
							}
							log.Info("[plugin] Parse:Response.Body.Append: %s ==> %s ", appstr, tplStr)
							rp.Response.Body.Append = tplStr
						}
					}
				}

			}
		}

	}
}

// PrintPlugin 打印插件
func (base *BasePlugin) PrintPlugin() {
	data, _ := yaml.Marshal(base)
	fmt.Println(string(data))
}

// CheckPlugin 检查插件格式是否正确
// Replace 中 body 的 File 和其他的不共存
func (base *BasePlugin) CheckPlugin() error {
	if base.Name == "" {
		return fmt.Errorf("[plugin: %s] Name cannot be empty", base.Name)
	}
	if base.Version == "" {
		return fmt.Errorf("[plugin: %s] Version cannot be empty ", base.Name)
	}
	if base.Description == "" {
		return fmt.Errorf("[plugin: %s] Description cannot be empty ", base.Name)
	}
	if base.WriteDate == "" {
		return errors.New(" WriteDate cannot be empty ")
	}
	if base.Author == "" {
		return fmt.Errorf("[plugin: %s] Author cannot be empty ", base.Name)
	}
	if base.Rule == nil || len(base.Rule) == 0 {
		return fmt.Errorf("[plugin: %s] Rule cannot be empty or nil ", base.Name)
	}
	// todo 检查分开做
	for _, rule := range base.Rule {
		if rule.URL == "" {
			return fmt.Errorf("[plugin: %s] rule.URL cannot be empty or nil  ", base.Name)
		}
		if rule.Match == "" {
			return fmt.Errorf("[plugin: %s.rule.%s] rule.Match cannot be empty or nil : [ %s ]", base.Name, rule.URL, strings.Join(MatchType, ","))
		}
		if !utils.EleInArray(rule.Match, MatchType) {
			return fmt.Errorf("[plugin: %s.rule.%s] rule.Match:%s must be: %s", base.Name, rule.URL, rule.Match, strings.Join(MatchType, ","))
		}
		// Replace 检查
		if rule.Replace != nil {
			for _, rp := range rule.Replace {
				if rp.Request == nil {
					return fmt.Errorf("[plugin: %s.rule.%s] Request cannot be empty or nil ", base.Name, rule.URL)
				}
				// 判断请求方法
				for _, m := range rp.Request.Method {
					if !utils.EleInArray(m, replace.Method) {
						return fmt.Errorf("[plugin: %s.rule.%s] method :%s not in[ %s ]", base.Name, rule.URL, m, strings.Join(replace.Method, ","))
					}
				}

				if rp.Response == nil {
					return nil
				}
				if rp.Response.Status < 0 || rp.Response.Status > 599 {
					return fmt.Errorf("[plugin: %s.rule.%s] 就你家网站响应码是这? Status:%d", base.Name, rule.URL, rp.Response.Status)
				}
				if rp.Response.Header != nil {
					//// Location 检查
					if lc, ok := rp.Response.Header["Location"]; ok {
						if rp.Response.Status/100 != 3 {
							log.Warn("[plugin: %s.rule.%s] Location:%s ,Status: %d, not set status code 3xx", base.Name, rule.URL, lc, rp.Response.Status)
						}
					}
				}
				if rp.Response.Cookie != nil {
					ck := rp.Response.Cookie
					if ck.SameSite > 4 || ck.SameSite < 0 {
						log.Fatal("[plugin: %s.rule.Cookie] URL:%s, SameSite value[1-4]: %s", base.Name, rule.URL, ck.SameSite)
					}
				}

				if rp.Response.Body != nil {
					if rp.Response.Body.File != "" {
						if !utils.FileExist(rp.Response.Body.File) {
							return fmt.Errorf("[plugin: %s.rule.%s] Response.Body.File not find file：%s", base.Name, rule.URL, rp.Response.Body.File)
						}
						if len(rp.Response.Body.ReplaceStr) != 0 {
							return fmt.Errorf("[plugin: %s.rule.%s] Response.Body File cannot coexist with replacestr and append ", base.Name, rule.URL)
						}
						if rp.Response.Body.Append != "" {
							return fmt.Errorf("[plugin: %s.rule.%s] Response.Body File cannot coexist with replacestr and append ", base.Name, rule.URL)
						}

					}
				}
				//判断请求方法
			}

		}
		//dump 检查
		if rule.Dump != nil {
			for _, rp := range rule.Dump {
				//判断请求方法
				for _, m := range rp.Request.Method {
					if !utils.EleInArray(m, replace.Method) {
						return fmt.Errorf("[plugin: %s.rule.%s] method :%s not in[ %s ]", base.Name, rule.URL, m, strings.Join(replace.Method, ","))
					}
				}
				if rp.Response == nil {
					return nil
				}
				if rp.Response.Status < 0 || rp.Response.Status > 599 {
					return fmt.Errorf("[plugin: %s.rule.%s] 就你家网站响应码是这? Status: %d", base.Name, rule.URL, rp.Response.Status)
				}
			}
		}
	}

	return nil
}
