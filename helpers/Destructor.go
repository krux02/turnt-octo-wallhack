package helpers

type Destructor interface {
	Delete()
}

type DependenceList struct {
	dependencies []Destructor
}

func (this *DependenceList) Delete() {
	for _, dep := range this.dependencies {
		dep.Delete()
	}
	this.dependencies = nil
}

func (this *DependenceList) Bind(dest Destructor) {
	this.dependencies = append(this.dependencies, dest)
}
