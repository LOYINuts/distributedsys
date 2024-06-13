package registry

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

const ServerPort = ":3000"
const ServicesURL = "http://localhost" + ServerPort + "/services"

type myregistry struct {
	registrations []Registration //已经注册的服务
	mutex         *sync.Mutex    //多个线程可能并发访问registrations，加锁保证线程安全
}

func (r *myregistry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	return nil
}

var reg = myregistry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.Mutex),
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")
	switch r.Method {
	case http.MethodPost:
		// dec存储着json格式的请求的body
		dec := json.NewDecoder(r.Body)
		var tmpreg Registration
		// 解码存在tmpreg
		err := dec.Decode(&tmpreg)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v with URL: %s\n", tmpreg.ServiceName,
			tmpreg.ServiceURL)
		err = reg.add(tmpreg)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
