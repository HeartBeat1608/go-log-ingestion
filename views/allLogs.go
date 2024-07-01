package views

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/HeartBeat1608/logmaster/internals"
)

type ViewManager struct {
	c *internals.DBConnectionManager
}

func NewViewManager(c *internals.DBConnectionManager) *ViewManager {
	return &ViewManager{c}
}

func writeString(w http.ResponseWriter, status int, text string) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(text))
}

type LogItem struct {
	Id        int32  `json:"id" db:"id"`
	Timestamp string `json:"timestamp" db:"timestamp"`
	Message   string `json:"message" db:"message"`
}
type LogResults struct {
	Logs     []LogItem
	Service  string
	Count    int
	Skip     int
	Limit    int
	Page     int
	Previous int
	Next     int
	Last     int
}

func (vm *ViewManager) RenderAllLogs(w http.ResponseWriter, r *http.Request) {
	var (
		filename = "templates/logs.html"
		service  = r.PathValue("service")
		l        = r.URL.Query().Get("limit")
		s        = r.URL.Query().Get("skip")
	)

	var limit int = 20
	if l != "" {
		x, err := strconv.Atoi(l)
		if err == nil {
			limit = max(x, 5)
		}
	}

	var skip int = 0
	if s != "" {
		x, err := strconv.Atoi(s)
		if err == nil {
			skip = max(x, 0)
		}
	}

	t, err := template.ParseFiles(filename)
	if err != nil {
		writeString(w, 500, err.Error())
		return
	}

	db, err := vm.c.GetConnection(service)
	if err != nil {
		writeString(w, 500, err.Error())
		return
	}

	row := db.QueryRowx("select count(id) as count from logs", nil)
	var count int
	if err = row.Scan(&count); err != nil {
		writeString(w, 500, err.Error())
		return
	}

	rows, err := db.Queryx("select id, timestamp, message from logs order by timestamp desc limit ?, ?", skip, limit)
	if err != nil {
		writeString(w, 500, err.Error())
		return
	}

	var res = LogResults{
		Service:  service,
		Logs:     []LogItem{},
		Count:    count,
		Limit:    limit,
		Skip:     skip,
		Page:     skip / limit,
		Next:     min(skip+limit, count),
		Previous: max(skip-limit, 0),
		Last:     count - limit,
	}

	for rows.Next() {
		var tmp LogItem
		if err = rows.StructScan(&tmp); err != nil {
			writeString(w, 500, err.Error())
			return
		}
		res.Logs = append(res.Logs, tmp)
	}

	if err = t.Execute(w, res); err != nil {
		writeString(w, 500, err.Error())
		return
	}
}
