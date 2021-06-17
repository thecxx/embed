package internal

type P struct {
	Symbol string
	Path   string
}

type M interface {
	Requires() []P
	Dependents() []M
	Statement() string
}

type GzipModule struct {
}

func (g GzipModule) Requires() []P {
	return nil
}

func (g GzipModule) Dependents() []M {
	return nil

}

func (g GzipModule) Statement() string {
	return ""
}

type Page struct {
	ms []M
}

func (p *Page) InsertModule(m M) {
	p.ms = append(p.ms, m)
}

func (p *Page) DeleteModule(m M) {

}

type Package struct {
	Pages []Page
	Name  string
	Path  string
}

func (p *Package) InsertPage(page *Page) {

}

func (p *Package) DeletePage(page *Page) {

}
