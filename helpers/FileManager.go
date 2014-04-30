package helpers

type ManagerSource interface {
	Update(string)
}

var managerSourceMapping = make(map[string]ManagerSource)
var filechanges = make(chan string)

func Manage(source ManagerSource, filename string) {
	managerSourceMapping[filename] = source
	go Listen(filename, filechanges)
}
