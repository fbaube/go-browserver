package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	// "github.com/elijahmorg/lhtmx/htmx" 
	"github.com/fbaube/go-browserver/htmx"
	"github.com/labstack/echo/v4"
)

var TZ *time.Location

func init() {
     var e error
     TZ, e  = time.LoadLocation("Local")
     if  e != nil { panic(e) }
}

// hRenderTodosRoute is the handler for GET ("/")
func hRenderTodosRoute(c echo.Context) error {
     	fmt.Println("Called hRenderTodosRoute: GET(\"/\")")
	return c.HTML(http.StatusOK, htmx.RenderTodos(htmx.FakeServerSideTodosDB))
}

// hToggleTodoRoute is the handler for POST("/toggle/:id")
func hToggleTodoRoute(c echo.Context) error {
     	fmt.Println("Called hToggleTodoRoute: POST(\"/toggle/:id\")")
     // fmt.Printf("TOGGLE:" + ctxAsString(c))
        c.Echo().Logger.Info("TOGGLE:" + ctxAsString(c))
        println("TOGGLE:" + ctxAsString(c))
	// js.Global().Get("console").Call("log", "Hello from Go WebAssembly!")
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedTodo htmx.Todo
	for i, todo := range htmx.FakeServerSideTodosDB {
		if todo.ID == id {
			htmx.FakeServerSideTodosDB[i].Done = !todo.Done
			updatedTodo = htmx.FakeServerSideTodosDB[i]
			break
		}
	}
	return c.HTML(http.StatusOK, htmx.CreateTodoNode(updatedTodo).Render())
}

// hAddTodoRoute is the handler for POST("/add")
func hAddTodoRoute(c echo.Context) error {
	fmt.Println("Called hAddTodoRoute: POST(\"/add\")")
	todoTitle := c.FormValue("newTodo")
	fmt.Println("TodoTitle: ", todoTitle)
	if todoTitle == "" {
		return c.String(http.StatusBadRequest, "no title provided")
	}
	// Create a single value
	todo := htmx.Todo{ID: len(htmx.FakeServerSideTodosDB) + 1, Title: todoTitle, 
	     	Done: false, TimeID: time.Now().UnixNano() }
	if todoTitle != "" {
		htmx.FakeServerSideTodosDB = append(htmx.FakeServerSideTodosDB, todo)
	}
	fmt.Println("hello world: from addTodoRoute: writing response")
	err := c.HTML(http.StatusOK, htmx.RenderBody(htmx.FakeServerSideTodosDB))

	return err
}

// hSyncTodos handles POST("/sync") on the server side. 
// So, titles in our local client's DB override same-names
// in htmx.FakeServerSideTodosDB, and then the result is
// "pushed" back into the server-side DB. 
func hSyncTodos(c echo.Context) error {
	var todos []htmx.Todo
	err := c.Bind(&todos)
	if err != nil {
		return err
	}
	// This call assumes that todos is newer, 
	// and its contents override same-titles 
	// in htmx.FakeServerSideTodosDB, our
	// let's-pretend server-side datastore. 
	todos, _ = htmx.MergeChanges(todos, htmx.FakeServerSideTodosDB) 
	
	// This is where we "write back to the server DB".
	htmx.FakeServerSideTodosDB = todos

	// Dump out the new state of the server-side DB.
	c.JSON(http.StatusOK, htmx.FakeServerSideTodosDB)
	fmt.Println("Called hSyncTodos: POST(\"/sync\")")
	return nil
}

// hGetTodos handles GET("/sync") on the server side.
// So, it returns everything found in htmx.FakeServerSideTodosDB, 
// our let's-pretend server-side datastore.
func hGetTodos(c echo.Context) error {
	fmt.Println("Called hGetTodos: GET(\"/sync\")")
	return c.JSON(http.StatusOK, htmx.FakeServerSideTodosDB)
}

// ServerDelay middleware inserts artificial latency
func ServerDelay(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		time.Sleep(1000 * time.Millisecond)
		return next(c)
	}
}

// =====================
//  FUNCS for DATA SYNC 
// =====================

// SyncToServer syncs the current data state with the server in a Go routine.
//
// "This is a very rudimentary data sync method and it definitely has issues.
// It is another idea I’d like to follow up with to see how to do this better
// and more robustly. That is not the main point here, so I just got something
// very basic working."
// .
func SyncToServer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		SyncData()
		return next(c)
	}
}

func SyncData() {
	go syncDataRoutine()
}

func syncDataRoutine() {
	b := bytes.NewBuffer([]byte(""))
	json.NewEncoder(b).Encode(htmx.FakeServerSideTodosDB)

	// First we send our entire (presumedly updated) DB
	// to the server. 
	resp, err := http.Post("http://localhost:3000/sync",
	      "application/json", b)
	if err != nil {
		fmt.Println("error syncing data: ", err)
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("bad status code for sync")
		fmt.Println("error syncing data: ", err)
		return
	}
	// Now we pretend to receive it at the server. 
	todos := make([]htmx.Todo, 0)
	err = json.NewDecoder(resp.Body).Decode(&todos)
	if err != nil {
		fmt.Println("error decoding response: ", err)
	}
	// Here we pretend to use the just-received (and updated)
	// DB while merging in older stuff from todos. 
	todos, err = htmx.MergeChanges(htmx.FakeServerSideTodosDB, todos) // local, server 
	if err != nil {
		fmt.Println("error merging: ", err)
	}
	// Then we put the whole mess in our (fake) server-side DB.
	htmx.FakeServerSideTodosDB = todos
}

// GetData gets data from the server. 
func GetData() error {

	fmt.Println("get data from server for syncing")
	// This calls hGetTodos 
	resp, err := http.Get("http://localhost:3000/sync")
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("bad status code for sync")
		fmt.Println(err)
		return err
	}
	// todos is a temp var for fetching from the server.
	todos := make([]htmx.Todo, 0)
	json.NewDecoder(resp.Body).Decode(&todos)

	// This is basickly a no-op, cos it assigns data
	// fetched from our fake server DB back to same.
	todos, err = htmx.MergeChanges(htmx.FakeServerSideTodosDB, todos) // local, server 
	if err != nil {
		return err
	}

	htmx.FakeServerSideTodosDB = todos
	return nil
}

// = = = = = =

/*
// Context represents the context of the current HTTP request.
// It holds request and response objects, path, path parameters,
// data, and registered handler.
type Context interface {
		Request() *http.Request
		Response() *Response // "".Writer http.ResponseWriter
		// Path returns the registered path for the handler.
		Path() string
		
		// Param returns path parameter by name.
		Param(name string) string
 		// ParamNames returns path parameter names.
		ParamNames() []string
		// ParamValues returns path parameter values.
		ParamValues() []string
		// QueryParam returns the query param for the provided name.
		QueryParam(name string) string
		// QueryParams returns the query parameters as `url.Values`.
		QueryParams() url.Values
		// QueryString returns the URL query string.
		QueryString() string
		// FormValue returns the form field value for the provided name.
		FormValue(name string) string
		// FormParams returns the form parameters as `url.Values`.
		FormParams() (url.Values, error)
		// FormFile returns the multipart form file for the provided name
		FormFile(name string) (*multipart.FileHeader, error)
		// MultipartForm returns the multipart form.
		MultipartForm() (*multipart.Form, error)
		// Cookie returns the named cookie provided in the request.
		Cookie(name string) (*http.Cookie, error)
		// Cookies returns the HTTP cookies sent with the request.
		Cookies() []*http.Cookie
		// Get retrieves data from the context.
		Get(key string) interface{}
		// Bind binds the request body into provided type `i`. 
		// The default binder does it based on Content-Type header.
		Bind(i interface{}) error
		// Validate validates provided `i`.
		// It is usually called after `Context#Bind()`.
		// Validator must be registered using `Echo#Validator`.
		Validate(i interface{}) error
	}
	context struct {
		request  *http.Request
		response *Response
		path     string
		pnames   []string
		pvalues  []string
		query    url.Values
		handler  HandlerFunc
		store    Map
		echo     *Echo
	}
)
*/

func ctxAsString(p echo.Context) string {
     return fmt.Sprintf("CTX<%+v>", *(p.Request()))
}

