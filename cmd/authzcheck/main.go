package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jonathanselby/authorisation-dual-control/internal/authz"
	"os"
)

const version = "0.1.0"

func main() {
	var pol, req string
	var v bool
	flag.StringVar(&pol, "policy", "", "policy JSON")
	flag.StringVar(&req, "request", "", "request JSON")
	flag.BoolVar(&v, "version", false, "print version")
	flag.Parse()
	if v {
		fmt.Println("authzcheck " + version)
		return
	}
	if pol == "" || req == "" {
		flag.Usage()
		os.Exit(2)
	}
	var p authz.Policy
	var r authz.Request
	if e := read(pol, &p); e != nil {
		bad(e)
	}
	if e := read(req, &r); e != nil {
		bad(e)
	}
	d, e := authz.Evaluate(p, r)
	if e != nil {
		bad(e)
	}
	json.NewEncoder(os.Stdout).Encode(d)
	if !d.Allowed {
		os.Exit(1)
	}
}
func read(p string, v any) error {
	b, e := os.ReadFile(p)
	if e != nil {
		return e
	}
	return json.Unmarshal(b, v)
}
func bad(e error) { fmt.Fprintln(os.Stderr, e); os.Exit(3) }
