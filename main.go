package main

func main() {
	cfg := new(ConfigServer)
	cfg.init()
	worker := new(jobs)
	worker.cfg = cfg
	worker.init()
	server := new(Server)
	server.cfg = cfg
	server.worker = worker
	server.start()
}
