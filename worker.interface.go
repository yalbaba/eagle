package eagle

type Task interface{
	Handle(...interface{})(...interface{},error)
}

