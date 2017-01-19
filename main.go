package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"

	"github.com/ajays20078/go-http-logger"
	"github.com/gorilla/mux"
)

const goToPath = "/go/"
const defaultURL = "http://thingsforhackers.com/blog"
const fullURLKey = "X-Full-URL"
const dumpURLKey = "X-Dump-URL"
const updateTokenKey = "X-Update-Token"

//Flags groups command line args together
type Flags struct {
	dbPath      string
	portNum     int
	updateToken string
	debug       bool
}

var flags Flags

func init() {
	if usr, err := user.Current(); err != nil {
		log.Fatal(err)
	} else {
		flag.StringVar(&flags.dbPath,
			"dbPath", fmt.Sprintf("%s/url.db", usr.HomeDir), "Path to Dbase file")
	}
	flag.IntVar(&flags.portNum, "portNum", 8080, "Port to listen on")
	flag.StringVar(&flags.updateToken, "updateToken", "", "Optional token required for DB modification operations")
	flag.BoolVar(&flags.debug, "debug", false, "Enable debug output")
}

func main() {

	flag.Parse()

	if flags.debug {
		fmt.Printf("Port: %d\n", flags.portNum)
		fmt.Printf("DbPath: %s\n", flags.dbPath)
		fmt.Printf("Token: %s\n", flags.updateToken)
	}

	us := new(URLstore)

	if err := us.Start(flags.dbPath); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create URLStore: %s\n", err.Error())
		os.Exit(1)
	}
	defer us.Stop()

	r := setUpRoutes(us)

	if flags.debug {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", flags.portNum), httpLogger.WriteLog(r, os.Stdout)))
	} else {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", flags.portNum), r))
	}
}

/*
writeResponse wraps up setting of the response status and a message
*/
func writeResponse(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	if len(msg) > 0 {
		io.WriteString(w, fmt.Sprintf("<h1>%s</h1>", msg))
	}
}

//setUpRoutes will setup handler methods for the routes
func setUpRoutes(us *URLstore) *mux.Router {

	r := mux.NewRouter().StrictSlash(true)

	//Default handler
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, defaultURL, http.StatusFound)
	})

	//Go handler
	//Handles both GET - do a redirect and PUT add a new shortName mapping
	r.HandleFunc(goToPath+"{shortName}", func(w http.ResponseWriter,
		r *http.Request) {

		vars := mux.Vars(r)
		shortName := vars["shortName"]
		longURL, err := us.Retrieve(shortName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Lookup failed")
			return
		}

		if r.Method == http.MethodPut {
			//Setting
			updateToken := r.Header.Get(updateTokenKey)
			if flags.updateToken != "" && flags.updateToken != updateToken {
				writeResponse(w, http.StatusUnauthorized, "")
				return
			}
			if longURL != "" {
				writeResponse(w, http.StatusConflict,
					fmt.Sprintf("%s already mapped to %s", shortName, longURL))
				return
			}
			//Read new url from header
			//fmt.Println(r.Header())
			longURL = r.Header.Get(fullURLKey)
			if longURL == "" {
				//Not set
				writeResponse(w, http.StatusBadRequest,
					fmt.Sprintf("%s is missing from header", fullURLKey))
				return
			}
			if err := us.Store(shortName, longURL); err != nil {
				writeResponse(w, http.StatusInternalServerError,
					fmt.Sprintf("Failed to store %s", shortName))
				return
			}
			writeResponse(w, http.StatusCreated, "")
		} else {
			//Assume GET
			if longURL == "" {
				//Empty, i.e. shortName not known
				writeResponse(w, http.StatusNotFound,
					fmt.Sprintf("Can't map %s to a URL", shortName))
				return
			}
			if r.Header.Get(dumpURLKey) == "true" {
				writeResponse(w, http.StatusOK, fmt.Sprintf("%s ---> %s", shortName, longURL))
				return
			}
			http.Redirect(w, r, longURL, http.StatusFound)
			return
		}

	}).Methods(http.MethodGet, http.MethodPut)

	http.Handle("/", r)

	return r
}
