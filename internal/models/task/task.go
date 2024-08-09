package task

type Task struct {
	ID string `json:"id"`
	Info
}

type Info struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Response struct {
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}
