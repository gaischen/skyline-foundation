package discover

//
type Discover interface {
}

type Service interface {

}

type Watcher interface {
	Watch(func(service Service))
}
