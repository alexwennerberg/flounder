// Copied from https://git.sr.ht/~sircmpwn/kineto/tree/master/item/main.go
package main

import (
	"fmt"
	"html/template"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"git.sr.ht/~adnano/go-gemini"
)

func proxyGemini(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("404 Not found"))
		return
	}
	var spath []string
	if r.URL.Path == "/" {
		http.Redirect(w, r, "gemini.circumlunar.space", http.StatusSeeOther)
		return
	} else if r.URL.Path == "/robots.txt" {
		temp := c.TemplatesDirectory
		http.ServeFile(w, r, path.Join(temp, "proxy-robots.txt"))
		return
	} else {
		spath = strings.SplitN(r.URL.Path, "/", 3)
	}
	req := gemini.Request{}
	var err error
	req.Host = spath[1]
	if len(spath) > 2 {
		req.URL, err = url.Parse(fmt.Sprintf("gemini://%s/%s", spath[1], spath[2]))
	} else {
		req.URL, err = url.Parse(fmt.Sprintf("gemini://%s/", spath[1]))
	}
	req.URL.RawQuery = r.URL.RawQuery
	client := gemini.Client{
		Timeout: 60 * time.Second,
	}

	if h := (url.URL{Host: req.Host}); h.Port() == "" {
		req.Host += ":1965"
	}

	resp, err := client.Do(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Gateway error: %v", err)
		return
	}
	defer resp.Body.Close()

	switch resp.Status {
	case 10, 11:
		// TODO accept input
		w.WriteHeader(http.StatusInternalServerError)
		return
	case 20:
		break // OK
	case 30, 31:
		to, err := url.Parse(resp.Meta)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(fmt.Sprintf("Gateway error: bad redirect %v", err)))
		}
		next := req.URL.ResolveReference(to)
		if next.Scheme != "gemini" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("This page is redirecting you to %s", next.String())))
			return
		}
		next.Host = r.URL.Host
		next.Scheme = r.URL.Scheme
		next.Path = fmt.Sprintf("/%s/%s", req.URL.Host, next)
		w.Header().Add("Location", next.String())
		w.WriteHeader(http.StatusFound)
		w.Write([]byte("Redirecting to " + next.String()))
		return
	case 40, 41, 42, 43, 44:
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "The remote server returned %d: %s", resp.Status, resp.Meta)
		return
	case 50, 51:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "The remote server returned %d: %s", resp.Status, resp.Meta)
		return
	case 52, 53, 59:
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "The remote server returned %d: %s", resp.Status, resp.Meta)
		return
	default:
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "Proxy does not understand Gemini response status %d", resp.Status)
		return
	}

	m, _, err := mime.ParseMediaType(resp.Meta)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("Gateway error: %d %s: %v",
			resp.Status, resp.Meta, err)))
		return
	}

	_, raw := r.URL.Query()["raw"]
	acceptsGemini := strings.Contains(r.Header.Get("Accept"), "text/gemini")
	if m != "text/gemini" || raw || acceptsGemini {
		w.Header().Add("Content-Type", resp.Meta)
		io.Copy(w, resp.Body)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	parse, _ := gemini.ParseText(resp.Body)
	htmlDoc := textToHTML(req.URL, parse)
	if strings.HasSuffix(r.URL.Path, "/") {
		r.URL.Path = path.Dir(r.URL.Path)
	}
	data := struct {
		SiteBody  template.HTML
		PageTitle string
		GeminiURI *url.URL
		URI       *url.URL
		Config    Config
	}{template.HTML(htmlDoc.Content), htmlDoc.Title, req.URL, r.URL, c}

	err = t.ExecuteTemplate(w, "user_page.html", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}
}
