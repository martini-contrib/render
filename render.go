// Package render is a middleware for Martini that provides easy JSON serialization and HTML template rendering.
//
//  package main
//
//  import (
//    "encoding/xml"
//
//    "github.com/go-martini/martini"
//    "github.com/martini-contrib/render"
//  )
//
//  type Greeting struct {
//    XMLName xml.Name `xml:"greeting"`
//    One     string   `xml:"one,attr"`
//    Two     string   `xml:"two,attr"`
//  }
//
//  func main() {
//    m := martini.Classic()
//    m.Use(render.Renderer()) // reads "templates" directory by default
//
//    m.Get("/html", func(r render.Render) {
//      r.HTML(200, "mytemplate", nil)
//    })
//
//    m.Get("/json", func(r render.Render) {
//      r.JSON(200, "hello world")
//    })
//
//    m.Get("/xml", func(r render.Render) {
//      r.XML(200, Greeting{One: "hello", Two: "world"})
//    })
//
//    m.Run()
//  }
package render

import (
	//	"encoding/json"
	//	"encoding/xml"

	"html/template"
	"log"
	//	"io"
	"io/ioutil"
	//	"net/http"
	"os"
	"path/filepath"
	"strings"

	//"github.com/go-martini/martini"
	"github.com/gin-gonic/gin/render"
)

const (
	ContentType    = "Content-Type"
	ContentLength  = "Content-Length"
	ContentBinary  = "application/octet-stream"
	ContentText    = "text/plain"
	ContentJSON    = "application/json"
	ContentHTML    = "text/html"
	ContentXHTML   = "application/xhtml+xml"
	ContentXML     = "text/xml"
	defaultCharset = "UTF-8"
)

// Provides a temporary buffer to execute templates into and catch errors.
//var bufpool *bpool.BufferPool

// Included helper functions for use when rendering html
//var helperFuncs = template.FuncMap{
//	"yield": func() (string, error) {
//		return "", fmt.Errorf("yield called with no layout defined")
//	},
//	"current": func() (string, error) {
//		return "", nil
//	},
//}

// Render is a service that can be injected into a Martini handler. Render provides functions for easily writing JSON and
// HTML templates out to a http Response.
//type Render interface {
//	// JSON writes the given status and JSON serialized version of the given value to the http.ResponseWriter.
//	JSON(status int, v interface{})
//	// HTML renders a html template specified by the name and writes the result and given status to the http.ResponseWriter.
//	HTML(status int, name string, v interface{}, htmlOpt ...HTMLOptions)
//	// XML writes the given status and XML serialized version of the given value to the http.ResponseWriter.
//	XML(status int, v interface{})
//	// Data writes the raw byte array to the http.ResponseWriter.
//	Data(status int, v []byte)
//	// Text writes the given status and plain text to the http.ResponseWriter.
//	Text(status int, v string)
//	// Error is a convenience function that writes an http status to the http.ResponseWriter.
//	Error(status int)
//	// Status is an alias for Error (writes an http status to the http.ResponseWriter)
//	Status(status int)
//	// Redirect is a convienience function that sends an HTTP redirect. If status is omitted, uses 302 (Found)
//	Redirect(location string, status ...int)
//	// Template returns the internal *template.Template used to render the HTML
//	Template() *template.Template
//	// Header exposes the header struct from http.ResponseWriter.
//	Header() http.Header
//}

type Render struct {
	Templates       map[string]*template.Template
	TemplatesDir    string
	Layout          string
	Exts            []string
	TemplateFuncMap map[string]interface{}
	Debug           bool
}

func New() Render {
	r := Render{

		Templates: map[string]*template.Template{},
		// TemplatesDir holds the location of the templates
		TemplatesDir: "app/views/",
		// Layout is the file name of the layout file
		Layout: "layouts/base",
		// Ext is the file extension of the rendered templates
		Exts: []string{".html"},
		// Template's function map
		TemplateFuncMap: nil,
		// Debug enables debug mode
		Debug: false,
	}

	return r
}

// Delims represents a set of Left and Right delimiters for HTML template rendering
//type Delims struct {
//	// Left delimiter, defaults to {{
//	Left string
//	// Right delimiter, defaults to }}
//	Right string
//}

//// Options is a struct for specifying configuration options for the render.Renderer middleware
//type Options struct {
//	// Directory to load templates. Default is "templates"
//	Directory string
//	// Layout template name. Will not render a layout if "". Defaults to "".
//	Layout string
//	// Extensions to parse template files from. Defaults to [".tmpl"]
//	Extensions []string
//	// Funcs is a slice of FuncMaps to apply to the template upon compilation. This is useful for helper functions. Defaults to [].
//	Funcs []template.FuncMap
//	// Delims sets the action delimiters to the specified strings in the Delims struct.
//	Delims Delims
//	// Appends the given charset to the Content-Type header. Default is "UTF-8".
//	Charset string
//	// Outputs human readable JSON
//	IndentJSON bool
//	// Outputs human readable XML
//	IndentXML bool
//	// Prefixes the JSON output with the given bytes.
//	PrefixJSON []byte
//	// Prefixes the XML output with the given bytes.
//	PrefixXML []byte
//	// Allows changing of output to XHTML instead of HTML. Default is "text/html"
//	HTMLContentType string
//}

// HTMLOptions is a struct for overriding some rendering Options for specific HTML call
type HTMLOptions struct {
	// Layout template name. Overrides Options.Layout.
	Layout string
}

// Renderer is a Middleware that maps a render.Render service into the Martini handler chain. An single variadic render.Options
// struct can be optionally provided to configure HTML rendering. The default directory for templates is "templates" and the default
// file extension is ".tmpl".
//
// If MARTINI_ENV is set to "" or "development" then templates will be recompiled on every request. For more performance, set the
// MARTINI_ENV environment variable to "production"
//func Renderer(options ...Options) martini.Handler {
//	opt := prepareOptions(options)
//	cs := prepareCharset(opt.Charset)
//	t := compile(opt)
//	bufpool = bpool.NewBufferPool(64)
//	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
//		var tc *template.Template
//		if martini.Env == martini.Dev {
//			// recompile for easy development
//			tc = compile(opt)
//		} else {
//			// use a clone of the initial template
//			tc, _ = t.Clone()
//		}
//		c.MapTo(&renderer{res, req, tc, opt, cs}, (*Render)(nil))
//	}
//}

func prepareCharset(charset string) string {
	if len(charset) != 0 {
		return "; charset=" + charset
	}

	return "; charset=" + defaultCharset
}

//func prepareOptions(options []Options) Options {
//	var opt Options
//	if len(options) > 0 {
//		opt = options[0]
//	}
//
//	// Defaults
//	if len(opt.Directory) == 0 {
//		opt.Directory = "templates"
//	}
//	if len(opt.Extensions) == 0 {
//		opt.Extensions = []string{".tmpl"}
//	}
//	if len(opt.HTMLContentType) == 0 {
//		opt.HTMLContentType = ContentHTML
//	}
//
//	return opt
//}

func (r Render) compile() {
	log.Printf("compile")
	t := template.New(r.TemplatesDir)
	// t.Delims(options.Delims.Left, options.Delims.Right)
	// parse an initial template in case we don't have any
	template.Must(t.Parse("SlothNinja"))

	filepath.Walk(r.TemplatesDir, func(path string, info os.FileInfo, perr error) (err error) {
		var (
			rp, ext, extension, name string
			tmpl                     *template.Template
			buf                      []byte
		)
		if perr != nil {
			err = perr
			return
		}

		if rp, err = filepath.Rel(r.TemplatesDir, path); err != nil {
			return
		}

		ext = getExt(rp)

		for _, extension = range r.Exts {
			if ext == extension {

				if buf, err = ioutil.ReadFile(path); err != nil {
					panic(err)
				}

				name = filepath.ToSlash((rp[0 : len(rp)-len(ext)]))
				if r.Debug {
					log.Printf("[GIN-debug] %s\n", name)
				}
				tmpl = t.New(name)

				// add our funcmaps
				tmpl.Funcs(r.TemplateFuncMap)

				// Bomb out if parse fails. We don't want any silent server starts.
				//r.Templates[name] = template.Must(tmpl.Funcs(helperFuncs).Parse(string(buf)))
				r.Templates[name] = template.Must(tmpl.Parse(string(buf)))
			}
		}

		return
	})
}

func getExt(s string) string {
	if strings.Index(s, ".") == -1 {
		return ""
	}
	return "." + strings.Join(strings.Split(s, ".")[1:], ".")
}

//type renderer struct {
//	http.ResponseWriter
//	req             *http.Request
//	t               *template.Template
//	opt             Options
//	compiledCharset string
//}
//
//func (r *renderer) JSON(status int, v interface{}) {
//	var result []byte
//	var err error
//	if r.opt.IndentJSON {
//		result, err = json.MarshalIndent(v, "", "  ")
//	} else {
//		result, err = json.Marshal(v)
//	}
//	if err != nil {
//		http.Error(r, err.Error(), 500)
//		return
//	}
//
//	// json rendered fine, write out the result
//	r.Header().Set(ContentType, ContentJSON+r.compiledCharset)
//	r.WriteHeader(status)
//	if len(r.opt.PrefixJSON) > 0 {
//		r.Write(r.opt.PrefixJSON)
//	}
//	r.Write(result)
//}
//
//func (r *renderer) HTML(status int, name string, binding interface{}, htmlOpt ...HTMLOptions) {
//	opt := r.prepareHTMLOptions(htmlOpt)
//	// assign a layout if there is one
//	if len(opt.Layout) > 0 {
//		r.addYield(name, binding)
//		name = opt.Layout
//	}
//
//	buf := bufpool.Get()
//	defer bufpool.Put(buf)
//
//	err := r.t.ExecuteTemplate(buf, name, binding)
//	if err != nil {
//		http.Error(r, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	// template rendered fine, write out the result
//	r.Header().Set(ContentType, r.opt.HTMLContentType+r.compiledCharset)
//	r.WriteHeader(status)
//	io.Copy(r, buf)
//}
//
//func (r *renderer) XML(status int, v interface{}) {
//	var result []byte
//	var err error
//	if r.opt.IndentXML {
//		result, err = xml.MarshalIndent(v, "", "  ")
//	} else {
//		result, err = xml.Marshal(v)
//	}
//	if err != nil {
//		http.Error(r, err.Error(), 500)
//		return
//	}
//
//	// XML rendered fine, write out the result
//	r.Header().Set(ContentType, ContentXML+r.compiledCharset)
//	r.WriteHeader(status)
//	if len(r.opt.PrefixXML) > 0 {
//		r.Write(r.opt.PrefixXML)
//	}
//	r.Write(result)
//}
//
//func (r *renderer) Data(status int, v []byte) {
//	if r.Header().Get(ContentType) == "" {
//		r.Header().Set(ContentType, ContentBinary)
//	}
//	r.WriteHeader(status)
//	r.Write(v)
//}
//
//func (r *renderer) Text(status int, v string) {
//	if r.Header().Get(ContentType) == "" {
//		r.Header().Set(ContentType, ContentText+r.compiledCharset)
//	}
//	r.WriteHeader(status)
//	r.Write([]byte(v))
//}
//
//// Error writes the given HTTP status to the current ResponseWriter
//func (r *renderer) Error(status int) {
//	r.WriteHeader(status)
//}
//
//func (r *renderer) Status(status int) {
//	r.WriteHeader(status)
//}
//
//func (r *renderer) Redirect(location string, status ...int) {
//	code := http.StatusFound
//	if len(status) == 1 {
//		code = status[0]
//	}
//
//	http.Redirect(r, r.req, location, code)
//}
//
//func (r *renderer) Template() *template.Template {
//	return r.t
//}
//
//func (r *renderer) addYield(name string, binding interface{}) {
//	funcs := template.FuncMap{
//		"yield": func() (template.HTML, error) {
//			buf := bufpool.Get()
//			defer bufpool.Put(buf)
//
//			err := r.t.ExecuteTemplate(buf, name, binding)
//			// return safe html here since we are rendering our own template
//			return template.HTML(buf.String()), err
//		},
//		"current": func() (string, error) {
//			return name, nil
//		},
//	}
//	r.t.Funcs(funcs)
//}
//
//func (r *renderer) prepareHTMLOptions(htmlOpt []HTMLOptions) HTMLOptions {
//	if len(htmlOpt) > 0 {
//		return htmlOpt[0]
//	}
//
//	return HTMLOptions{
//		Layout: r.opt.Layout,
//	}
//}

func (r Render) Init() Render {
	log.Printf("Init")
	r.compile()
	//	log.Printf("Init debug: %v", r.Debug)
	//	layout := r.TemplatesDir + r.Layout + r.Ext
	//	log.Printf("Init layout: %v", layout)
	//
	//	viewDirs, _ := filepath.Glob(r.TemplatesDir + "**/*" + r.Ext)
	//	log.Printf("Init viewDirs: %v", viewDirs)
	//
	//	for _, view := range viewDirs {
	//		renderName := r.getRenderName(view)
	//		if r.Debug {
	//			log.Printf("[GIN-debug] %-6s %-25s --> %s\n", "LOAD", view, renderName)
	//		}
	//		r.AddFromFiles(renderName, layout, view)
	//	}
	//
	return r
}

func (r Render) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r.Templates[name],
		Data:     data,
	}
}
