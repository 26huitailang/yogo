package env

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/26huitailang/yogo/framework/contract"
)

type YogoEnv struct {
	folder string
	maps   map[string]string
}

func NewYogoEnv(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("NewYogoEnv: params error")
	}
	folder := params[0].(string)
	log.Println(folder)
	yogoEnv := &YogoEnv{
		folder: folder,
		maps:   map[string]string{"APP_ENV": contract.EnvDevelopment},
	}

	file := path.Join(folder, ".env")

	f, err := os.Open(file)
	if err == nil {
		defer f.Close()
		reader := bufio.NewReader(f)
		for {
			line, _, c := reader.ReadLine()
			if c == io.EOF {
				break
			}
			s := bytes.SplitN(line, []byte{'='}, 2)
			if len(s) < 2 {
				continue
			}
			key := string(s[0])
			val := string(s[1])
			yogoEnv.maps[key] = val
		}
	}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) < 2 {
			continue
		}
		yogoEnv.maps[pair[0]] = pair[1]
	}
	return yogoEnv, nil
}

func (env *YogoEnv) AppEnv() string {
	return env.Get("APP_ENV")
}

func (env *YogoEnv) IsExist(key string) bool {
	_, ok := env.maps[key]
	return ok
}

func (env *YogoEnv) Get(s string) string {
	if val, ok := env.maps[s]; ok {
		return val
	}
	return ""
}

func (env *YogoEnv) All() map[string]string {
	return env.maps
}
