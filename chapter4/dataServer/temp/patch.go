package temp

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func patch(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	b, e := ioutil.ReadFile(infoFile)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	datFile := infoFile + ".dat"
	f, e := os.OpenFile(datFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, e = io.Copy(f, r.Body)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	i := strings.Split(string(b), ":")
	size, _ := strconv.ParseInt(i[1], 0, 64)
	actual := info.Size()
	if actual > size {
		//	os.Remove(datFile)
		//	os.Remove(infoFile)
		log.Println("actual size", actual, "exceeds", size)
		//	w.WriteHeader(http.StatusInternalServerError)
	}
}
