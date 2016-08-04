//
// watch-cli - watch configuration files and run command when files change
//
// (C) Philip Schlump, 2013-2015.
// Version: 1.0.1
// BuildNo: 003
//
// /Users/corwin/Projects/watch-cli
//
// Example:
// 		./watch-cli -c "make refresh_site"
//
// Issues:
// 		0. Adding files should work - API to add file, then re-init of watch-cli
//			/api/list/add-cfg-file?fn=, then...
//			Some sort of signal to watch-cli to re-init?
// 		1. For a farm need to send messages to "every" tab-server1.go -- Serve Farm
//

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	flags "github.com/jessevdk/go-flags"
	tr "github.com/pschlump/godebug" // "../go-lib/tr" // "www.2c-why.com/go-lib/tr"

	"github.com/howeyc/fsnotify"
)

var optsRecursive = false

// Desc: get a list of filenames and directorys
func getFilenames(dir string) (filenames, dirs []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil
	}
	for _, fstat := range files {
		if !strings.HasPrefix(string(fstat.Name()), ".") {
			if fstat.IsDir() {
				dirs = append(dirs, fstat.Name())
			} else {
				filenames = append(filenames, fstat.Name())
			}
		}
	}
	return
}

var db1 bool = false

type Jar struct { // A CookieJar
	lk      sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

// SetCookies handles the receipt of the cookies in a reply for the
// given URL.  It may or may not choose to save the cookies, depending
// on the jar's policy and implementation.
func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	// fmt.Printf ( "cookies=%v\n", cookies )
	jar.cookies[u.Host] = cookies
	jar.lk.Unlock()
}

// Cookies returns the cookies to send in a request for the given URL.
// It is up to the implementation to honor the standard cookie use
// restrictions such as in RFC 6265.
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}
func (jar *Jar) CookiesHost(host string) []*http.Cookie {
	return jar.cookies[host]
}

func (jar *Jar) SpillCookies() {
	for i, v := range jar.cookies {
		for j, w := range v {
			// fmt.Printf ( "Path(%s) [%d]: %v\n", i, j, w )
			if db1 {
				fmt.Printf("Path(%s) [%d]: %s\n", i, j, tr.SVar(w))
			}
		}
	}
}

func (jar *Jar) getX() string {
	for i, v := range jar.cookies {
		for j, w := range v {
			if false {
				fmt.Printf("Path(%s) [%d]: %v\n", i, j, w)
			}
			if w.Name == "XSRF-TOKEN" {
				return w.Value
			}
		}
	}
	return ""
}

func doGet(client *http.Client, url string) string {
	r1, e0 := client.Get(url)
	if e0 != nil {
		fmt.Printf("Error!!!!!!!!!!! %v, %s\n", e0, tr.LF())
		return "Error"
	}
	rv, e1 := ioutil.ReadAll(r1.Body)
	if e1 != nil {
		fmt.Printf("Error!!!!!!!!!!! %v, %s\n", e1, tr.LF())
		return "Error"
	}
	r1.Body.Close()
	// fmt.Printf("Register New User Response: %s\n",string(rv))
	// fmt.Printf("Register New User Response 6: %s\n",string(rv[6:]))
	// Xyzzy - let's convert this to a func that deal with it.
	if string(rv[0:6]) == ")]}',\n" {
		rv = rv[6:]
	}

	return string(rv)
}

// -----------------------------------------------------------------------------------------------------------------------------------
// Exists reports whether the named file or directory exists.
// -----------------------------------------------------------------------------------------------------------------------------------
func DirExists(name string) bool {
	if fstat, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	} else {
		if fstat.IsDir() {
			return true
		}
	}
	return false
}

// -------------------------------------------------------------------------------------------------
// Exists reports whether the named file or directory exists.
// -------------------------------------------------------------------------------------------------
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// -------------------------------------------------------------------------------------------------
// Get a list of filenames and directorys.
// -------------------------------------------------------------------------------------------------
func GetFilenames(dir string) (filenames, dirs []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil
	}
	for _, fstat := range files {
		if !strings.HasPrefix(string(fstat.Name()), ".") {
			if fstat.IsDir() {
				dirs = append(dirs, fstat.Name())
			} else {
				filenames = append(filenames, fstat.Name())
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------
func EscapeStr(v string, on bool) string {
	if on {
		return html.EscapeString(v)
	} else {
		return v
	}
}

// -------------------------------------------------------------------------------------------------
// xyzzy - need "fast" version of "CompareFiles" with some limits on what it will use "fast"
// compare for - .jpg,.gif,.png fiels - a fiel size before uses fast etc.   Compare Size?
// Compare name?  What is the "fast" compare for rsync? -- Calculate Hashes for each and
// keep them around?
// -------------------------------------------------------------------------------------------------
func CompareFiles(cmpFile string, refFile string) bool {
	cmp, err := ioutil.ReadFile(cmpFile)
	if err != nil {
		fmt.Printf("Unable to read %s\n", cmpFile)
		return false
	}

	if Exists(refFile) {
		ref, err := ioutil.ReadFile(refFile)
		if err != nil {
			fmt.Printf("Unable to read %s\n", refFile)
			return false
		}
		if len(ref) != len(cmp) { // xyzzy - Could be faster - just check lenths on disk - if diff then return false
			return false
		}
		if string(ref) != string(cmp) {
			return false
		}
	} else {
		return false
	}
	return true
}

// -------------------------------------------------------------------------------------------------
// Get a list of filenames and directories.
// -------------------------------------------------------------------------------------------------
func GetFilenamesRecrusive(dir string) (filenames, dirs []string, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}
	//for ii, fstat := range files {
	//	fmt.Printf("Top files %d:[%s]\n", ii, fstat.Name())
	//}
	for _, fstat := range files {
		if !strings.HasPrefix(string(fstat.Name()), ".") {
			if fstat.IsDir() {
				name := fstat.Name()
				dirs = append(dirs, dir+"/"+name)
				// fmt.Printf("Recursive dir [%s]\n", dir+"/"+name)
				tf, td, err := GetFilenamesRecrusive(dir + "/" + name)
				if err != nil {
					return nil, nil, err
				}
				filenames = append(filenames, tf...)
				dirs = append(dirs, td...)
			} else {
				name := fstat.Name()
				name = dir + "/" + name
				// fmt.Printf("dir %s ->%s<-\n", dir, name)
				filenames = append(filenames, name)
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------------------
func FilterArray(re string, inArr []string) (outArr []string) {
	var validID = regexp.MustCompile(re)

	outArr = make([]string, 0, len(inArr))
	for k := range inArr {
		if validID.MatchString(inArr[k]) {
			outArr = append(outArr, inArr[k])
		}
	}
	// fmt.Printf ( "output = %v\n", outArr )
	return
}

type SqlCfgLoaded struct {
	FileName string
	ErrorMsg string
}

var SqlCfgFilesLoaded []SqlCfgLoaded

var run_cmd = false
var mutex = &sync.Mutex{}

func getRunCmd() bool {
	mutex.Lock()
	x := run_cmd
	mutex.Unlock()
	return x
}
func getClearRunCmd() bool {
	mutex.Lock()
	x := run_cmd
	run_cmd = false
	mutex.Unlock()
	return x
}
func setRunCmd(x bool) {
	mutex.Lock()
	run_cmd = x
	mutex.Unlock()
}

type JSONCfg struct {
	FilesToWatch []string
}

func main() {

	tmp1, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	flist := tmp1[1:]

	// Get to the correct directory to run stuff from
	if opts.CdTo != "" {
		err := os.Chdir(opts.CdTo)
		if err != nil {
			fmt.Printf("Unable to change directories to %s, error: %s\n", opts.CdTo, err)
			os.Exit(1)
		}
	}

	var gCfg JSONCfg

	// Read in JSON config file if it exists
	if Exists(opts.Cfg) {
		fb, err := ioutil.ReadFile(opts.Cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading config file %s, %s\n", opts.Cfg, err)
			os.Exit(1)
		}

		err = json.Unmarshal(fb, &gCfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading config file %s, %s - file did not parse\n", opts.Cfg, err)
			os.Exit(1)
		}

		fmt.Printf("Config file %s read in, watching %s\n", opts.Cfg, gCfg.FilesToWatch)

		flist = append(flist, gCfg.FilesToWatch...)
	}

	// Delcare stuff --------------------------------------------------------------------------------------
	cli_buf := strings.Split(opts.Cmd, " ")

	// Do periodic Stuff ----------------------------------------------------------------------------------
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
				if getClearRunCmd() {
					// setRunCmd(false)
					cmd := exec.Command(cli_buf[0], cli_buf[1:]...)
					// cmd.Stdin = strings.NewReader("some input")
					var out bytes.Buffer
					cmd.Stdout = &out
					err := cmd.Run()
					if err != nil {
						fmt.Printf("Run Errors: %s", err)
					}
					fmt.Printf("%s\n", out.String())
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// Watch file for changes -----------------------------------------------------------------------------
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				fmt.Printf("Event: %+v\n", ev) // log.Println("event:", ev)
				name := ev.Name
				isRen := ev.IsRename()
				fmt.Printf("Caught an event, %s\n", tr.SVar(ev))

				setRunCmd(true)

				if isRen {
					err = watcher.Watch(name)
					if err != nil {
						fmt.Printf("Failed to set watch on %s, %s, %s\n", name, err, tr.LF())
					}
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	fmt.Printf("***************************************\n")
	fmt.Printf("* watching %s \n", flist)
	fmt.Printf("***************************************\n")
	for _, fn := range flist {
		if DirExists(fn) {
			var fns []string
			if !optsRecursive {
				fns, _ = GetFilenames(fn)
			} else {
				fns, _, _ = GetFilenamesRecrusive(fn)
			}
			for _, fn0 := range fns {
				err = watcher.Watch(fn0)
				if err != nil {
					fmt.Printf("Failed to set watch on %s, %s, %s\n", fn0, err, tr.LF())
				}
			}
		} else {
			err = watcher.Watch(fn)
			if err != nil {
				fmt.Printf("Failed to set watch on %s, %s, %s\n", fn, err, tr.LF())
			}
		}
	}

	<-done

	/* ... do stuff ... */
	watcher.Close()
}

/* vim: set noai ts=4 sw=4: */
