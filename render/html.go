/*

HTML RENDER

*/

package render

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	// beeViewPathTemplates caching map and supported template file extensions per view
	beeViewPathTemplates = make(map[string]map[string]*template.Template)

	templatesLock sync.RWMutex

	// beeTemplateExt stores the template extension which will build
	beeTemplateExt = []string{"tpl", "html", "gohtml", "htm"}

	// beeTemplatePreprocessors stores associations of extension -> preprocessor handler
	beeTemplateEngines = map[string]templatePreProcessor{}

	beeTemplateFS = defaultFSFunc

	defaultViewPath = "views"
	defaultRunMode  = "dev"

	defaultDelims = Delims{Left: "{{", Right: "}}"}

	//utf8 text/html ContentType
	htmlContentType = []string{"text/html; charset=utf-8"}
)

//HTMLRender a simple html render
type HTMLRender struct {
	ViewPath   string
	Name       string
	AutoRender bool
	FuncMap    template.FuncMap
	Data       interface{}
	Delims     Delims
	RunMode    string
}

// NewHTMLRender return a Render  interface
func (m HTMLRender) NewHTMLRender(dir string, funcMap template.FuncMap, delims Delims, autoRender bool, runMode string) Render {
	if runMode == "" {
		runMode = defaultRunMode
	}
	render := HTMLRender{
		ViewPath:   dir,
		AutoRender: autoRender,
		FuncMap:    funcMap,
		Delims:     delims,
		RunMode:    runMode,
	}
	defaultDelims = render.Delims
	for key, item := range funcMap {
		templateFuncMap[key] = item
	}

	// add view path
	addViewPath(render.ViewPath)
	return render
}

// Render render implement
func (m HTMLRender) Render(w http.ResponseWriter, name string, data interface{}) error {
	if !m.AutoRender {
		return nil
	}
	m.Name = name
	m.Data = data

	b, err := m.renderBytes()
	if err != nil {
		return err
	}
	m.WriteContentType(w)
	_, err = w.Write(b)
	return err
}

// WriteContentType implement
func (m HTMLRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, htmlContentType)
}

// renderBytes renderBytes
func (m HTMLRender) renderBytes() ([]byte, error) {
	buf, err := m.renderTemplate()
	return buf.Bytes(), err
}

// renderTemplate
func (m HTMLRender) renderTemplate() (bytes.Buffer, error) {
	var buf bytes.Buffer

	if m.RunMode == defaultRunMode {
		files := []string{m.Name}
		BuildTemplate(m.ViewPath, files...)
	}
	return buf, ExecuteTemplate(&buf, m.Name, m.ViewPath, m.RunMode, m.Data)
}

// addViewPath addViewPath
func addViewPath(viewPath string) error {
	if _, exist := beeViewPathTemplates[viewPath]; exist {
		return nil
	}
	beeViewPathTemplates[viewPath] = make(map[string]*template.Template)
	err := BuildTemplate(viewPath)
	return err
}

// ExecuteTemplate applies the template with name  to the specified data object,
// writing the output to wr.
// A template will be executed safely in parallel.
func ExecuteTemplate(wr io.Writer, name, viewPath string, runMode string, data interface{}) error {
	if viewPath == "" {
		viewPath = defaultViewPath
	}
	return ExecuteViewPathTemplate(wr, name, viewPath, runMode, data)
}

// ExecuteViewPathTemplate applies the template with name and from specific viewPath to the specified data object,
// writing the output to wr.
// A template will be executed safely in parallel.
func ExecuteViewPathTemplate(wr io.Writer, name string, viewPath string, runMode string, data interface{}) error {
	if runMode == defaultRunMode {
		templatesLock.RLock()
		defer templatesLock.RUnlock()
	}
	if beeTemplates, ok := beeViewPathTemplates[viewPath]; ok {
		if t, ok := beeTemplates[name]; ok {
			var err error
			if t.Lookup(name) != nil {
				err = t.ExecuteTemplate(wr, name, data)
			} else {
				err = t.Execute(wr, data)
			}
			return err
		}
		panic("can't find template file in the path:" + viewPath + "/" + name)
	}
	panic("Unknown view path:" + viewPath)
}

// AddFuncMap let user to register a func in the template.
func AddFuncMap(key string, fn interface{}) error {
	templateFuncMap[key] = fn
	return nil
}

type templatePreProcessor func(root, path string, funcs template.FuncMap) (*template.Template, error)

type templateFile struct {
	root  string
	files map[string][]string
}

// visit will make the paths into two part,the first is subDir (without tf.root),the second is full path(without tf.root).
// if tf.root="views" and
// paths is "views/errors/404.html",the subDir will be "errors",the file will be "errors/404.html"
// paths is "views/admin/errors/404.html",the subDir will be "admin/errors",the file will be "admin/errors/404.html"
func (tf *templateFile) visit(paths string, f os.FileInfo, err error) error {
	if f == nil {
		return err
	}
	if f.IsDir() || (f.Mode()&os.ModeSymlink) > 0 {
		return nil
	}
	if !HasTemplateExt(paths) {
		return nil
	}

	replace := strings.NewReplacer("\\", "/")
	file := strings.TrimLeft(replace.Replace(paths[len(tf.root):]), "/")
	subDir := filepath.Dir(file)

	tf.files[subDir] = append(tf.files[subDir], file)
	return nil
}

// HasTemplateExt return this path contains supported template extension of beego or not.
func HasTemplateExt(paths string) bool {
	for _, v := range beeTemplateExt {
		if strings.HasSuffix(paths, "."+v) {
			return true
		}
	}
	return false
}

// AddTemplateExt add new extension for template.
func AddTemplateExt(ext string) {
	for _, v := range beeTemplateExt {
		if v == ext {
			return
		}
	}
	beeTemplateExt = append(beeTemplateExt, ext)
}

// BuildTemplate will build all template files in a directory.
// it makes beego can render any template file in view directory.
func BuildTemplate(dir string, files ...string) error {
	var err error
	fs := beeTemplateFS()
	f, err := fs.Open(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.New("dir open err")
	}
	defer f.Close()

	beeTemplates, ok := beeViewPathTemplates[dir]
	if !ok {
		panic("Unknown view path: " + dir)
	}

	self := &templateFile{
		root:  dir,
		files: make(map[string][]string),
	}
	err = Walk(fs, dir, func(path string, f os.FileInfo, err error) error {
		return self.visit(path, f, err)
	})
	if err != nil {
		log.Printf("Walk() returned %v\n", err)
		return err
	}
	buildAllFiles := len(files) == 0
	for _, v := range self.files {
		for _, file := range v {
			if buildAllFiles || InSlice(file, files) {
				templatesLock.Lock()
				ext := filepath.Ext(file)
				var t *template.Template
				if len(ext) == 0 {
					t, err = getTemplate(self.root, fs, file, v...)
				} else if fn, ok := beeTemplateEngines[ext[1:]]; ok {
					t, err = fn(self.root, file, templateFuncMap)
				} else {
					t, err = getTemplate(self.root, fs, file, v...)
				}
				if err != nil {
					log.Printf("parse template err: %v %v \n", file, err)
					templatesLock.Unlock()
					return err
				}
				beeTemplates[file] = t
				templatesLock.Unlock()
			}
		}
	}
	return nil
}

func getTplDeep(root string, fs http.FileSystem, file string, parent string, t *template.Template) (*template.Template, [][]string, error) {
	var fileAbsPath string
	var rParent string
	var err error
	if strings.HasPrefix(file, "../") {
		rParent = filepath.Join(filepath.Dir(parent), file)
		fileAbsPath = filepath.Join(root, filepath.Dir(parent), file)
	} else {
		rParent = file
		fileAbsPath = filepath.Join(root, file)
	}
	f, err := fs.Open(fileAbsPath)
	if err != nil {
		panic("can't find template file:" + file)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, [][]string{}, err
	}
	t, err = t.New(file).Parse(string(data))
	if err != nil {
		return nil, [][]string{}, err
	}
	reg := regexp.MustCompile(defaultDelims.Left + "[ ]*template[ ]+\"([^\"]+)\"")
	allSub := reg.FindAllStringSubmatch(string(data), -1)
	for _, m := range allSub {
		if len(m) == 2 {
			tl := t.Lookup(m[1])
			if tl != nil {
				continue
			}
			if !HasTemplateExt(m[1]) {
				continue
			}
			_, _, err = getTplDeep(root, fs, m[1], rParent, t)
			if err != nil {
				return nil, [][]string{}, err
			}
		}
	}
	return t, allSub, nil
}

func getTemplate(root string, fs http.FileSystem, file string, others ...string) (t *template.Template, err error) {
	t = template.New(file).Delims(defaultDelims.Left, defaultDelims.Right).Funcs(templateFuncMap)
	var subMods [][]string
	t, subMods, err = getTplDeep(root, fs, file, "", t)
	if err != nil {
		return nil, err
	}
	t, err = _getTemplate(t, root, fs, subMods, others...)

	if err != nil {
		return nil, err
	}
	return
}

func _getTemplate(t0 *template.Template, root string, fs http.FileSystem, subMods [][]string, others ...string) (t *template.Template, err error) {
	t = t0
	for _, m := range subMods {
		if len(m) == 2 {
			tpl := t.Lookup(m[1])
			if tpl != nil {
				continue
			}
			//first check filename
			for _, otherFile := range others {
				if otherFile == m[1] {
					var subMods1 [][]string
					t, subMods1, err = getTplDeep(root, fs, otherFile, "", t)
					if err != nil {
						log.Printf("template parse file err: %v \n", err)
					} else if len(subMods1) > 0 {
						t, err = _getTemplate(t, root, fs, subMods1, others...)
					}
					break
				}
			}
			//second check define
			for _, otherFile := range others {
				var data []byte
				fileAbsPath := filepath.Join(root, otherFile)
				f, err := fs.Open(fileAbsPath)
				if err != nil {
					f.Close()
					log.Printf("template file parse error, not success open file: %v \n", err)
					continue
				}
				data, err = ioutil.ReadAll(f)
				f.Close()
				if err != nil {
					log.Printf("template file parse error, not success open file: %v \n", err)
					continue
				}
				reg := regexp.MustCompile(defaultDelims.Left + "[ ]*define[ ]+\"([^\"]+)\"")
				allSub := reg.FindAllStringSubmatch(string(data), -1)
				for _, sub := range allSub {
					if len(sub) == 2 && sub[1] == m[1] {
						var subMods1 [][]string
						t, subMods1, err = getTplDeep(root, fs, otherFile, "", t)
						if err != nil {
							log.Printf("template parse file err: %v\n", err)
						} else if len(subMods1) > 0 {
							t, err = _getTemplate(t, root, fs, subMods1, others...)
							if err != nil {
								log.Printf("template parse file err: %v\n", err)
							}
						}
						break
					}
				}
			}
		}

	}
	return
}

type templateFSFunc func() http.FileSystem

func defaultFSFunc() http.FileSystem {
	return FileSystem{}
}

// SetTemplateFSFunc set default filesystem function
func SetTemplateFSFunc(fnt templateFSFunc) {
	beeTemplateFS = fnt
}

func InSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
