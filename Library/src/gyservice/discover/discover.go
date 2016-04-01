package discover

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"github.com/coreos/go-etcd/etcd"
	"gylogger"
	"google.golang.org/grpc"
	. "gyservice/service"
)

const (
	DEFAULT_ETCD = "http://uat-sz.lansion.cn:4001"
	DEFAULT_SERVICE_PATH = "/yanzhi"
	DEFAULT_NAME_FILE = "/yanzhi/names"
	DEFAULT_DIAL_TIMEOUT = 10 * time.Second
	RETRY_DELAY = 10 * time.Second
)

type client struct {
	key  string
	conn *grpc.ClientConn
}

type service struct {
	clients []client
	idx     uint32
}

type service_pool struct {
	services          map[string]*service
	service_names     map[string]bool // store names.txt
	enable_name_check bool
	client_pool       sync.Pool
	sync.RWMutex
}

var (
	_default_pool service_pool
)

func init() {
	_default_pool.init()
	_default_pool.connect_all(DEFAULT_SERVICE_PATH)
	go _default_pool.watcher()
}

func (p *service_pool) init() {
	// etcd client
	machines := []string{DEFAULT_ETCD}
	if env := os.Getenv("ETCD_HOST"); env != "" {
		machines = strings.Split(env, ";")
	}
	p.client_pool.New = func() interface{} {
		return etcd.NewClient(machines)
	}

	p.services = make(map[string]*service)
}

// get stored handlermap name
func (p *service_pool) load_names() {
	p.service_names = make(map[string]bool)
	client := p.client_pool.Get().(*etcd.Client)
	defer func() {
		p.client_pool.Put(client)
	}()

	// get the keys under directory
	logger.Debug("reading names:", DEFAULT_NAME_FILE)
	resp, err := client.Get(DEFAULT_NAME_FILE, false, false)
	if err != nil {
		logger.Debug(err)
		return
	}

	// validation check
	if resp.Node.Dir {
		logger.Debug("names is not a file")
		return
	}

	// split names
	names := strings.Split(resp.Node.Value, "\n")
	logger.Debug("all handlermap names:", names)
	for _, v := range names {
		p.service_names[DEFAULT_SERVICE_PATH + "/" + strings.TrimSpace(v)] = true
	}

	p.enable_name_check = true
}

// connect to all services
func (p *service_pool) connect_all(directory string) {
	client := p.client_pool.Get().(*etcd.Client)
	defer func() {
		p.client_pool.Put(client)
	}()

	// get the keys under directory
	resp, err := client.Get(directory, true, true)
	if err != nil {
		logger.Debug(err)
		return
	}

	// validation check
	if !resp.Node.Dir {
		logger.Debug("not a directory")
		return
	}

	for _, node := range resp.Node.Nodes {
		if node.Dir {
			// handlermap directory
			for _, service := range node.Nodes {
				p.add_service(service.Key, service.Value)
			}
		}
	}
	logger.Debug("services add complete")
}

// watcher for data change in etcd directory
func (p *service_pool) watcher() {
	client := p.client_pool.Get().(*etcd.Client)
	defer func() {
		p.client_pool.Put(client)
	}()

	for {
		ch := make(chan *etcd.Response, 10)
		go func() {
			for {
				if resp, ok := <-ch; ok {
					if resp.Node.Dir {
						continue
					}
					key, value := resp.Node.Key, resp.Node.Value
					if value == "" {
						logger.Debugf("node delete: %v.", key)
						p.remove_service(key)
					} else {
						logger.Debugf("node add: %v %v.", key, value)
						p.add_service(key, value)
					}
				} else {
					return
				}
			}
		}()

		logger.Debug("Watching:", DEFAULT_SERVICE_PATH)
		_, err := client.Watch(DEFAULT_SERVICE_PATH, 0, true, ch, nil)
		if err != nil {
			logger.Debug(err)
		}
		<-time.After(RETRY_DELAY)
	}
}

// add a handlermap
func (p *service_pool) add_service(key, value string) {
	p.Lock()
	defer p.Unlock()
	service_name := filepath.Dir(key)
	// name check
	if p.enable_name_check && !p.service_names[service_name] {
		logger.Debugf("handlermap not in names: %v, ignored.", service_name)
		return
	}

	if p.services[service_name] == nil {
		p.services[service_name] = &service{}
		logger.Debugf("new handlermap type: %v.", service_name)
	}
	service := p.services[service_name]

	if conn, err := grpc.Dial(value, grpc.WithTimeout(DEFAULT_DIAL_TIMEOUT), grpc.WithInsecure()); err == nil {
		service.clients = append(service.clients, client{key, conn})
		logger.Debugf("handlermap added: %s : %s.", key, value)
	} else {
		logger.Debugf("handlermap not added: %s:%s. Error: %v.", key, value, err)
	}
}

// remove a handlermap
func (p *service_pool) remove_service(key string) {
	p.Lock()
	defer p.Unlock()
	service_name := filepath.Dir(key)
	service := p.services[service_name]
	if service == nil {
		logger.Debugf("no such handlermap %v.", service_name)
		return
	}

	for k := range service.clients {
		if service.clients[k].key == key {
			// deletion
			service.clients = append(service.clients[:k], service.clients[k + 1:]...)
			logger.Debugf("handlermap removed %v.", key)
			return
		}
	}
}

// provide a specific key for a handlermap, eg:
// path:/backends/snowflake, id:s1
//
// handlermap must be stored like /backends/xxx_service/xxx_id
func (p *service_pool) get_service_with_id(path string, id string) *grpc.ClientConn {
	p.RLock()
	defer p.RUnlock()
	service := p.services[path]
	if service == nil {
		return nil
	}

	if len(service.clients) == 0 {
		return nil
	}

	fullpath := string(path) + "/" + id
	for k := range service.clients {
		if service.clients[k].key == fullpath {
			return service.clients[k].conn
		}
	}

	return nil
}

func (p *service_pool) get_service(path string) *grpc.ClientConn {
	p.RLock()
	defer p.RUnlock()
	service := p.services[path]
	if service == nil {
		return nil
	}

	if len(service.clients) == 0 {
		return nil
	}
	idx := int(atomic.AddUint32(&service.idx, 1))
	return service.clients[idx % len(service.clients)].conn
}

// choose a handlermap randomly
func GetService(path string) *grpc.ClientConn {
	return _default_pool.get_service(path)
}

// get a specific handlermap instance with given path and id
func GetServiceWithId(path string, id string) *grpc.ClientConn {
	return _default_pool.get_service_with_id(path, id)
}

func GetServiceClient(path string) ServiceClient {
	conn := _default_pool.get_service(path)
	if conn != nil {
		return NewServiceClient(conn)
	}
	return nil
}