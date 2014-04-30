package helpers

type ManagerSource interface {
	Update(string)
}

var managerSourceMapping = map[string]ManagerSource{}
var filechanges = make(chan string)

func Manage(source ManagerSource, filename string) {
	managerSourceMapping[filename] = source
	go Listen(filename, filechanges)
}

func UpdateManagers() {
	b := true
	for b {
		select {
		case filename := <-filechanges:
			ms := managerSourceMapping[filename]
			ms.Update(filename)
		default:
			b = false
		}
	}
}
