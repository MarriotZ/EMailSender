package template

import (
	"email/config"
	"errors"
	"log"
	"sync"
)

var data *tempData

type tempData struct {
	sync.RWMutex
	templates map[string]*Template
}

type Template struct {
	ID          uint64
	Name        string
	ParamsCount int
}

func init() {
	data = &tempData{
		templates: make(map[string]*Template, 0),
	}
	if list, ok := config.Conf.Get("SesMailTemplate").([]interface{}); ok {
		for _, item := range list {
			if mp, ok := item.(map[string]interface{}); ok {
				id, _ := mp["id"].(int)
				t := &Template{
					ID:          uint64(id),
					Name:        mp["name"].(string),
					ParamsCount: mp["paramscount"].(int),
				}
				Register(t.Name, t)
			}
		}
	}
}

func Register(key string, temp *Template) error {
	if temp == nil {
		err := errors.New("邮件模板对象不能为空")
		log.Println(err)
		return err
	}
	data.Lock()
	defer data.Unlock()
	if _, ok := data.templates[key]; ok {
		err := errors.New("邮件模板不能重复添加")
		log.Println(err)
		return err
	}
	data.templates[key] = temp
	return nil
}

func GetTemplate(key string) (*Template, error) {
	data.RLock()
	defer data.RUnlock()
	t, ok := data.templates[key]
	if ok {
		return t, nil
	} else {
		err := errors.New("模板不存在")
		return nil, err
	}
}
