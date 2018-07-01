package main

import (
	"path/filepath"
	"strings"
	//"context"

	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	gohttp "net/http"
	"net/url"
	"os"

	"github.com/ipfs/go-ipfs-cmdkit/files"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

const (
	exampleHash = "QmducKywvonbyssngALAdq8aeGc5wDUvAEXEUFJa329XUw"
	shellUrl    = "localhost:3000"
)

//var path = "C:/project/test.txt"

type RadioButton struct {
	Name       string
	Value      string
	IsDisabled bool
	IsChecked  bool
	Text       string
}

type PageVariables struct {
	PageTitle        string
	PageRadioButtons []RadioButton
	Answer           string
}

func main() {
	gohttp.HandleFunc("/index", DisplayRadioButtons)
	gohttp.HandleFunc("/selected", UserSelected)
	gohttp.HandleFunc("/selected/file", UserSelectedFile)
	gohttp.HandleFunc("/selected/file/read", UserSelectedReadFile)
	log.Fatal(gohttp.ListenAndServe(":3000", nil))
}
func ShowFiles(w gohttp.ResponseWriter, r *gohttp.Request) {

}
func DisplayRadioButtons(w gohttp.ResponseWriter, r *gohttp.Request) {
	// Display some radio buttons to the user

	Title := "Which do you prefer?"
	MyRadioButtons := []RadioButton{
		RadioButton{"option", "create", false, false, "Create"},
		RadioButton{"option", "read", false, false, "Read"},
	}

	MyPageVariables := PageVariables{
		PageTitle:        Title,
		PageRadioButtons: MyRadioButtons,
	}

	t, err := template.ParseFiles("select.html") //parse the html file homepage.html
	if err != nil {                              // if there is an error
		log.Print("template parsing error: ", err) // log it
	}

	err = t.Execute(w, MyPageVariables) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                     // if there is an error
		log.Print("template executing error: ", err) //log it
	}

}

func UserSelected(w gohttp.ResponseWriter, r *gohttp.Request) {
	r.ParseForm()

	yourchoice := r.Form.Get("option")

	Title := "Your preferred choice"
	MyPageVariables := PageVariables{
		PageTitle: Title,
		Answer:    yourchoice,
	}

	t, err := template.ParseFiles("select.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, MyPageVariables)
	if err != nil {
		log.Print("template executing error: ", err)
	}
	switch MyPageVariables.Answer {
	case "create":
		var files []string

		root := "C:/project"
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})
		if err != nil {
			panic(err)
		}
		Title := "Which do you prefer?"
		MyRadioButtons := []RadioButton{

			RadioButton{"option", files[1], false, false, files[1]},
			RadioButton{"option", files[2], false, false, files[2]},
			RadioButton{"option", files[3], false, false, files[3]},
			RadioButton{"option", files[4], false, false, files[4]},
		}

		MyPageVariables := PageVariables{
			PageTitle:        Title,
			PageRadioButtons: MyRadioButtons,
		}

		t, err := template.ParseFiles("file.html") //parse the html file homepage.html
		if err != nil {                            // if there is an error
			log.Print("template parsing error: ", err) // log it
		}

		err = t.Execute(w, MyPageVariables) //execute the template and pass it the HomePageVars struct to fill in the gaps
		if err != nil {                     // if there is an error
			log.Print("template executing error: ", err) //log it
		}
	case "read":
		var files []string

		root := "C:/Users/chotu ram/Downloads/go-ipfs_v0.4.15_windows-amd64/go-ipfs/ipfs"
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})
		if err != nil {
			panic(err)
		}
		Title := "Which do you prefer?"
		MyRadioButtons := []RadioButton{

			RadioButton{"option", files[1], false, false, files[1]},
			RadioButton{"option", files[2], false, false, files[2]},
			RadioButton{"option", files[3], false, false, files[3]},
			RadioButton{"option", files[4], false, false, files[4]},
		}

		MyPageVariables := PageVariables{
			PageTitle:        Title,
			PageRadioButtons: MyRadioButtons,
		}

		t, err := template.ParseFiles("read.html") //parse the html file homepage.html
		if err != nil {                            // if there is an error
			log.Print("template parsing error: ", err) // log it
		}

		err = t.Execute(w, MyPageVariables) //execute the template and pass it the HomePageVars struct to fill in the gaps
		if err != nil {                     // if there is an error
			log.Print("template executing error: ", err) //log it

		}
	}

}

func UserSelectedFile(w gohttp.ResponseWriter, r *gohttp.Request) {
	r.ParseForm()

	yourchoice := r.Form.Get("option")

	Title := "Your preferred choice"
	MyPageVariables := PageVariables{
		PageTitle: Title,
		Answer:    yourchoice,
	}

	t, err := template.ParseFiles("file.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, MyPageVariables)
	if err != nil {
		log.Print("template executing error: ", err)
	}
	var _, er = os.Stat(MyPageVariables.Answer)
	if er != nil {
		fmt.Println(er)
	}

	var file, errr = os.OpenFile(MyPageVariables.Answer, os.O_RDWR, 0644)
	if isError(errr) {
		return
	}
	rc := io.Reader(file)

	//	Add(r)
	s := NewShell(shellUrl)
	mhash, err := s.Add(rc)
	fmt.Println(mhash)

	fmt.Println("==> done creating file", MyPageVariables.Answer)
}

func UserSelectedReadFile(w gohttp.ResponseWriter, r *gohttp.Request) {
	r.ParseForm()

	yourchoice := r.Form.Get("option")

	Title := "Your preferred choice"
	MyPageVariables := PageVariables{
		PageTitle: Title,
		Answer:    yourchoice,
	}

	t, err := template.ParseFiles("read.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, MyPageVariables)
	if err != nil {
		log.Print("template executing error: ", err)
	}
	var _, er = os.Stat(MyPageVariables.Answer)
	if er != nil {
		fmt.Println(er)
	}

	var file, errr = os.OpenFile(MyPageVariables.Answer, os.O_RDWR, 0644)
	if isError(errr) {
		return
	}
	var text = make([]byte, 1024)
	for {
		_, err = file.Read(text)
		if err == io.EOF {
			break
		}

		if err != nil && err != io.EOF {
			isError(err)
			break
		}
	}

	fmt.Println("==> done reading from file")
	fmt.Println(string(text))

}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

type Shell struct {
	url     string
	httpcli *gohttp.Client
}

type object struct {
	Hash string
}

func NewShell(url string) *Shell {
	c := &gohttp.Client{
		Transport: &gohttp.Transport{
			Proxy:             gohttp.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	}
	return NewShellWithClient(url, c)
}

func NewShellWithClient(url string, c *gohttp.Client) *Shell {
	if a, err := ma.NewMultiaddr(url); err == nil {
		_, host, err := manet.DialArgs(a)
		if err == nil {
			url = host
		}
	}
	return &Shell{
		url:     url,
		httpcli: c,
	}
}

func (s *Shell) Add(r io.Reader) (string, error) {
	return s.AddWithOpts(r, true, false)
}

func (s *Shell) AddWithOpts(r io.Reader, pin bool, rawLeaves bool) (string, error) {
	var rc io.ReadCloser
	if rclose, ok := r.(io.ReadCloser); ok {
		rc = rclose
	} else {
		rc = ioutil.NopCloser(r)
	}

	// handler expects an array of files
	fr := files.NewReaderFile("test.txt", "C:/project/test.txt", rc, nil)
	slf := files.NewSliceFile("test.txt", "C:/project/test.txt", []files.File{fr})
	fileReader := files.NewMultiFileReader(slf, true)
	req := NewRequest(context.Background(), s.url, "add")
	req.Body = fileReader
	req.Opts["progress"] = "false"
	if !pin {
		req.Opts["pin"] = "false"
	}

	if rawLeaves {
		req.Opts["raw-leaves"] = "true"
	}

	resp, err := req.Send(s.httpcli)
	if err != nil {
		return "", err
	}
	defer resp.Close()
	if resp.Error != nil {
		return "", err
	}

	var out object
	err = json.NewDecoder(resp.Output).Decode(&out)
	if err != nil {
		return "", err
	}

	return out.Hash, nil
}
func (r *Request) getURL() string {
	values := make(url.Values)
	for _, arg := range r.Args {
		values.Add("arg", arg)
	}
	for k, v := range r.Opts {
		values.Add(k, v)
	}
	return fmt.Sprintf("%s/%s?%s", r.ApiBase, r.Command, values.Encode())
}

func (r *Request) Send(c *gohttp.Client) (*Response, error) {
	url := r.getURL()
	req, err := gohttp.NewRequest("POST", url, r.Body)
	if err != nil {
		return nil, err
	}
	if fr, ok := r.Body.(*files.MultiFileReader); ok {
		req.Header.Set("Content-Type", "multipart/form-data; boundary="+fr.Boundary())
		req.Header.Set("Content-Disposition", "form-data: name=\"files\"")
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	contentType := resp.Header.Get("Content-Type")
	parts := strings.Split(contentType, ";")
	contentType = parts[0]
	nresp := new(Response)
	nresp.Output = resp.Body
	if resp.StatusCode >= gohttp.StatusBadRequest {
		e := &Error{
			Command: r.Command,
		}
		switch {
		case resp.StatusCode == gohttp.StatusNotFound:
			e.Message = "command not found"
		case contentType == "text/plain":
			out, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ipfs-shell: warning! response read error: %s\n", err)
			}
			e.Message = string(out)
		case contentType == "application/json":
			if err = json.NewDecoder(resp.Body).Decode(e); err != nil {
				fmt.Fprintf(os.Stderr, "ipfs-shell: warning! response unmarshall error: %s\n", err)
			}
		default:
			fmt.Fprintf(os.Stderr, "ipfs-shell: warning! unhandled response encoding: %s", contentType)
			out, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ipfs-shell: response read error: %s\n", err)
			}
			e.Message = fmt.Sprintf("unknown ipfs-shell error encoding: %q - %q", contentType, out)
		}
		nresp.Error = e
		nresp.Output = nil
		// drain body and close
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
	return nresp, nil
}

type Response struct {
	Output io.ReadCloser
	Error  *Error
}

type Error struct {
	Command string
	Message string
	Code    int
}

func NewRequest(ctx context.Context, url, command string, args ...string) *Request {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	opts := map[string]string{
		"encoding":       "json",
		"steam-channels": "true",
	}
	return &Request{
		ApiBase: url + "/selected",
		Command: command,
		Args:    args,
		Opts:    opts,
		Headers: make(map[string]string),
	}
}

type Request struct {
	ApiBase string
	Command string
	Args    []string
	Opts    map[string]string
	Body    io.Reader
	Headers map[string]string
}

func (r *Response) Close() error {
	if r.Output != nil {
		// always drain output (response body)
		ioutil.ReadAll(r.Output)
		return r.Output.Close()
	}
	return nil
}
