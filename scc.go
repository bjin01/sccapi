package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	//BaseURLV1 is a basic
	BaseURLV1 = "https://scc.suse.com/connect"
	//Defaultheader is required for scc
	Defaultheader = "application/vnd.scc.suse.com.v4+json"
	//Defaultheader1 optional
	Defaultheader1 = "gzip, deflate"
)

//Geter interface for all GET func
type Geter interface {
	GetResults(parturl string, prodresp *[]ResultsGet) ([]ResultsGet, http.Header)
	MyPagination(myHeader http.Header, url string)
}

//Config for yaml file
type Config struct {
	Username string `yaml:"user_name"`
	Password string `yaml:"password"`
}

//Client is a struct
type Client struct {
	BaseURL    string
	Uname      string
	Pword      string
	HTTPClient *http.Client
}

//ResultsGet a universal struct for GET responses
type ResultsGet struct {
	ID                  int      `json:"id"`
	Name                string   `json:"name"`
	Status              string   `json: "status"`
	ExpiresAt           string   `json: "expires_at"`
	SystemsCount        int      `json: "systems_count"`
	VirtualSystemsCount int      `json: "virtual_count"`
	Identifier          string   `json:"identifier"`
	Version             string   `json:"version"`
	Regcode             string   `json: "regcode"`
	Productclasses      []string `json: "product_classes"`
	Login               string   `json: "login"`
	Password            string   `json: "password"`
	LastSeenAt          string   `json: "last_seen_at"`
	DistroTarget        string   `json: "distro_target"`
	URL                 string   `json: "url"`
	InstallerUpdates    bool     `json: "installer_updates"`
}

//ParseFlags to read in -config parameter
func ParseFlags() (f string, i string, e error) {
	var configPath string
	var information string
	flag.StringVar(&configPath, "config", "", "path to scc credential config file")
	flag.StringVar(&information, "get", "", "enter data parts to query scc, e.g. products or subscriptions")
	flag.Parse()

	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("config file exists\n")
		if information != "" {
			return configPath, information, nil
		}
	} else {
		fmt.Printf("File does not exist\n")
		log.Fatal("No config file found. Exit with error")
	}

	return configPath, information, nil

}

//NewConfig to read in the yaml file with scc cred
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

//NewClient is a func
func NewClient(c *Config) *Client {
	return &Client{
		BaseURL: BaseURLV1,
		Uname:   c.Username,
		Pword:   c.Password,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

//GetResults to list licensed prod
func (c *Client) GetResults(parturl string, prodresp *[]ResultsGet) ([]ResultsGet, http.Header) {
	newURL := c.BaseURL + parturl
	req, err := http.NewRequest("GET", newURL, nil)
	req.Header.Add("accept", Defaultheader)
	//req.Header.Add("Accept-Encoding", Defaultheader1)
	req.SetBasicAuth(c.Uname, c.Pword)
	resp, err := c.HTTPClient.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	//fmt.Printf("Result: %s\n", body)
	//fmt.Printf("Headers:  %v\n", len(resp.Header))

	jsonErr := json.Unmarshal(body, prodresp)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return *prodresp, resp.Header
}

//MyPagination for paging
func (c *Client) MyPagination(myHeader http.Header, url string) {

	for a, b := range myHeader {
		if a == "Link" {
			//fmt.Printf("%v: \t%v\n", a, b)
			if len(b) != 0 {
				for _, x := range b {
					//fmt.Printf(c)
					//links := strings.Split(c, " ")
					httpslinks := strings.SplitN(x, ">", 2)
					//fmt.Printf("%v\n", httpslinks[0])
					totalpages := strings.SplitN(httpslinks[0], "=", -1)
					tpages, _ := strconv.Atoi(totalpages[1])
					if tpages != 0 {
						//fmt.Printf("%v\n", tpages)
						for p := 2; p < tpages; p++ {
							newurl := url + "?page=" + strconv.Itoa(p)
							prodresp := &[]ResultsGet{}
							retbody, _ := c.GetResults(newurl, prodresp)
							PrintResults(retbody)

						}
					}

				}
			}
		}
	}
}

//PrintResults to print the responses
func PrintResults(i []ResultsGet) {
	var isslice bool
	if len(i) != 0 {
		isslice = true
	} else {
		fmt.Printf("%v\n", i)
	}

	if isslice == true {
		for _, b := range i {
			fmt.Printf("\tID: %v\n", b.ID)
			if b.Name != "" {
				fmt.Printf("\tName: %v\n", b.Name)
			}
			if b.Status != "" {
				fmt.Printf("\tStatus: %v\n", b.Status)
			}
			if b.ExpiresAt != "" {
				fmt.Printf("\tExpires at: %v\n", b.ExpiresAt)
			}
			if b.Regcode != "" {
				fmt.Printf("\tRegistration Code: %v\n", b.Regcode)
			}
			if len(b.Productclasses) != 0 {
				fmt.Printf("\tProduct Class: %#v\n", b.Productclasses)
			}
			if b.SystemsCount != 0 {
				fmt.Printf("\tSystem Count: %v\n", b.SystemsCount)
			}
			if b.VirtualSystemsCount != 0 {
				fmt.Printf("\tVirtual Count: %v\n", b.VirtualSystemsCount)
			}
			if b.Identifier != "" {
				fmt.Printf("\tIdentifier: %v\n", b.Identifier)
			}
			if b.Version != "" {
				fmt.Printf("\tVersion: %v\n", b.Version)
			}
			if b.Login != "" {
				fmt.Printf("\tLogin: %v\n", b.Login)
			}
			if b.Password != "" {
				fmt.Printf("\tPassword: %v\n", b.Password)
			}
			if b.LastSeenAt != "" {
				fmt.Printf("\tLast see at: %v\n", b.LastSeenAt)
			}

			if b.DistroTarget != "" {
				fmt.Printf("\tDistro Target: %v\n", b.DistroTarget)
			}
			if b.URL != "" {
				fmt.Printf("\tUrl: %v\n", b.URL)
			}
			if b.InstallerUpdates {
				fmt.Printf("\tInstaller Updates: %v\n", b.InstallerUpdates)
			}

			fmt.Println()
		}
	}
}

func main() {

	configPath, information, err := ParseFlags()
	credential, err := NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if configPath == "" && information == "" {
		log.Fatal("Either config file or get data information not provided. Exit with err")
	}

	if strings.ToLower(configPath) != "" {
		myclient := NewClient(credential)
		var url string
		if information == "installer" {
			url = "/repositories/" + strings.ToLower(information)
		} else {
			url = "/organizations/" + strings.ToLower(information)
		}
		prodresp := &[]ResultsGet{}
		retbody, paginationHeader := myclient.GetResults(url, prodresp)
		PrintResults(retbody)
		var MyGet Geter
		MyGet = myclient
		MyGet.MyPagination(paginationHeader, url)
	}

}
