package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type dumResp struct {
	code  int
	url   string
	index int
}

var paths = map[string]dumResp{
	"/validpath1": dumResp{
		code:  302,
		index: 0,
		url:   "/validpath2",
	},
	"/validpath2": dumResp{
		code:  301,
		index: 1,
		url:   "/validpath3",
	},
	"/validpath3": dumResp{
		code:  302,
		index: 2,
		url:   "/validpath4",
	},
	"/validpath4": dumResp{
		code:  302,
		index: 3,
		url:   "/validpath5",
	},
	"/validpath5": dumResp{
		code:  200,
		index: 4,
	},
	"/invalidpath1": dumResp{
		code:  302,
		index: 0,
		url:   "/invalidpath2",
	},
	"/invalidpath2": dumResp{
		code:  400,
		index: 1,
		url:   "/invalidpath3",
	},
	"/longpath1": dumResp{
		code:  302,
		index: 0,
		url:   "/longpath2",
	},
	"/longpath2": dumResp{
		code:  302,
		index: 1,
		url:   "/longpath3",
	},
	"/longpath3": dumResp{
		code:  302,
		index: 2,
		url:   "/longpath4",
	},
	"/longpath4": dumResp{
		code:  302,
		index: 3,
		url:   "/longpath5",
	},
	"/longpath5": dumResp{
		code:  302,
		index: 4,
		url:   "/longpath6",
	},
	"/longpath6": dumResp{
		code:  302,
		index: 5,
		url:   "/longpath7",
	},
	"/longpath7": dumResp{
		code:  302,
		index: 6,
		url:   "/longpath8",
	},
	"/longpath8": dumResp{
		code:  302,
		index: 7,
		url:   "/longpath9",
	},
	"/longpath9": dumResp{
		code:  302,
		index: 8,
		url:   "/longpath10",
	},
	"/longpath10": dumResp{
		code:  302,
		index: 9,
		url:   "/longpath11",
	},
	"/longpath11": dumResp{
		code:  200,
		index: 10,
	},
}

func urlStackCreator(startURL string) []resp {
	resps := make([]resp, 0, 0)
	r, ok := paths[startURL]
	if !ok {
		return resps
	}
	resps = append(resps, resp{
		code: r.code,
		url:  startURL,
	})
	for ok {
		oldURL := r.url
		r, ok = paths[r.url]
		if ok {
			resps = append(resps, resp{
				code: r.code,
				url:  oldURL,
			})
		}
	}
	return resps
}

var tests = []struct {
	startURL string
	err      bool
	urlStack []resp
}{
	{
		startURL: "/validpath1",
		err:      false,
		urlStack: urlStackCreator("/validpath1"),
	},
	{
		startURL: "/invalidpath1",
		err:      true,
		urlStack: urlStackCreator("/invalidpath1"),
	},
	{
		startURL: "/longpath1",
		err:      false,
		urlStack: urlStackCreator("/longpath1"),
	},
}

func TestRoundTrip(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reponse, present := paths[r.URL.Path]
		if present {
			w.Header().Set("Location", reponse.url)
			w.WriteHeader(reponse.code)
			return
		}
		w.WriteHeader(200)
	}))
	tw := &transportWrapper{}
	client := &http.Client{
		Transport:     tw,
		CheckRedirect: checkRedirect(-1),
	}

	for _, test := range tests {
		tw.redirectTrail = []resp{}
		_, err := client.Get(ts.URL + test.startURL)
		if test.err && err == nil {
			t.Errorf("Expected: Error\nGot: nil\n")
			t.Fail()
			continue
		}
		if !test.err && err != nil {
			t.Errorf("Expected: Nil\nGot: %s\n", err)
			t.Fail()
			continue
		}
		if len(tw.redirectTrail) != len(test.urlStack) {
			t.Errorf("Expected: %v\nGot: %v\n", test.urlStack, tw.redirectTrail)
			t.Fail()
			continue
		}
		for i := 0; i < len(test.urlStack); i++ {
			if tw.redirectTrail[i].code != test.urlStack[i].code || tw.redirectTrail[i].url != ts.URL+test.urlStack[i].url {
				t.Errorf("Expected: %v\nGot: %v\n", test.urlStack, tw.redirectTrail)
				break
			}
		}
	}

	client = &http.Client{
		Transport:     tw,
		CheckRedirect: checkRedirect(9),
	}

	_, err := client.Get(ts.URL + "/longpath1")
	if err == nil {
		t.Errorf("Expected: Error\nGot: nil\n")
		t.Fail()
	}

	_, err = client.Get("abcd://efgh.ijk")
	if err == nil {
		t.Errorf("Expected: Error\nGot: nil\n")
		t.Fail()
	}
}
