package gozdrofitapi

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"
)

func TestHttpApi_Authenticate_Success(t *testing.T) {
	mockHttpResponseBody := "{\"User\":{\"Member\":{\"Id\":99,\"HomeClubId\":99,\"DefaultClubId\":99}}}"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCookie(
			authenticationCookie,
			"token",
			r.Host,
			w)
		_, _ = fmt.Fprintf(w, mockHttpResponseBody)
	}))
	defer server.Close()
	url, err := url.Parse(server.URL)
	require.NoError(t, err)
	api := NewHttpApi(*url, http.Client{Jar: &TestJar{}}, true)
	request := LoginRequest{true, "login", "password"}

	resp, err := api.Authenticate(request)

	require.NoError(t, err)
	expected := &LoginResponse{User: User{Member{
		Id:            99,
		HomeClubId:    99,
		DefaultClubId: 99,
	}}}

	assert.Equal(t, resp, expected)
	assert.True(t, api.Authenticated())
}

func TestHttpApi_DailyClasses_Success(t *testing.T) {
	mockHttpResponseBody := "{\"CalendarData\":[" +
		"{\"Classes\":[{\"Id\":1,\"Status\":\"Bookable\",\"StatusReason\":null,\"Name\":\"TBC\",\"StartTime\":\"2021-08-12T17:00:00\",\"BookingIndicator\":{\"Limit\":16,\"Available\":7},\"Users\":[{\"Id\":1,\"IsCurrentUser\":true}]}]}," +
		"{\"Classes\":[{\"Id\":2,\"Status\":\"Awaitable\",\"StatusReason\":null,\"Name\":\"Tabata\",\"StartTime\":\"2021-08-12T18:00:00\",\"BookingIndicator\":{\"Limit\":16,\"Available\":0},\"Users\":[]}]}" +
		"]}"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectCookie(t, authenticationCookie, r)

		_, _ = fmt.Fprintf(w, mockHttpResponseBody)
	}))
	defer server.Close()
	url, err := url.Parse(server.URL)
	require.NoError(t, err)
	api := NewHttpApi(*url, NewDefaultHttpClient(), true)
	setDummyAuthentication(api)
	request := DailyClassesRequest{99, Date{time.Now()}}

	resp, err := api.DailyClasses(request)

	require.NoError(t, err)
	expected := &DailyClassesResponse{
		CalendarData: []CalendarData{
			{
				Classes: []Class{
					{
						Id:               1,
						Status:           ClassStatusBookable,
						Name:             "TBC",
						StartTime:        DateTime{Time: time.Date(2021, 8, 12, 17, 0, 0, 0, time.UTC)},
						BookingIndicator: BookingIndicator{16, 7},
						Users: []ClassUser{
							{
								Id:            1,
								IsCurrentUser: true,
							},
						},
					}}},
			{
				Classes: []Class{
					{
						Id:               2,
						Status:           ClassStatusAwaitable,
						Name:             "Tabata",
						StartTime:        DateTime{Time: time.Date(2021, 8, 12, 18, 0, 0, 0, time.UTC)},
						BookingIndicator: BookingIndicator{16, 0},
						Users:            []ClassUser{},
					}}},
		}}
	assert.Equal(t, resp, expected)
}

func TestHttpApi_BookClass(t *testing.T) {
	mockHttpResponseBody := "{\"Tickets\":[{\"TimeTableEventId\":1,\"Name\":\"Pumba®\",\"StartTime\":\"2021-08-12T20:00:00\",\"ZoneName\":\"Zdrofit Centrum Krucza\",\"UserName\":\"Imie Nazwisko\",\"UserNumber\":\"5555555\",\"UserId\":111111,\"Trainer\":\"IMIE NAZWISKO\"}],\"ClassId\":1,\"UserId\":111111}"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectCookie(t, authenticationCookie, r)

		_, _ = fmt.Fprintf(w, mockHttpResponseBody)
	}))
	defer server.Close()
	url, err := url.Parse(server.URL)
	require.NoError(t, err)
	api := NewHttpApi(*url, NewDefaultHttpClient(), true)
	setDummyAuthentication(api)
	request := BookClassRequest{99}

	err = api.BookClass(request)

	assert.NoError(t, err)
}

func TestHttpApi_CancelClassBooking(t *testing.T) {
	mockHttpResponseBody := "{\"ClassId\":1,\"UserId\":111111}"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectCookie(t, authenticationCookie, r)

		_, _ = fmt.Fprintf(w, mockHttpResponseBody)
	}))
	defer server.Close()
	url, err := url.Parse(server.URL)
	require.NoError(t, err)
	api := NewHttpApi(*url, NewDefaultHttpClient(), true)
	setDummyAuthentication(api)
	request := CancelBookingRequest{99}

	err = api.CancelClassBooking(request)

	assert.NoError(t, err)
}

func setCookie(name string, value string, domain string, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	})
}

func expectCookie(t *testing.T, name string, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == name {
			return
		}
	}
	t.Errorf("Cookie %v is expected but wasn't provided in the request", name)
}

func setDummyAuthentication(api Api) {
	cookie := &http.Cookie{
		Name:    authenticationCookie,
		Value:   "token",
		Expires: time.Now().Add(time.Hour * 24),
	}
	cookies := []*http.Cookie{cookie}
	httpApi, _ := api.(*httpApi)
	httpApi.httpClient.Jar.SetCookies(&httpApi.baseUrl, cookies)
}

type TestJar struct {
	m      sync.Mutex
	perURL map[string][]*http.Cookie
}

func (j *TestJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.m.Lock()
	defer j.m.Unlock()
	if j.perURL == nil {
		j.perURL = make(map[string][]*http.Cookie)
	}
	j.perURL[u.Host] = cookies
}

func (j *TestJar) Cookies(u *url.URL) []*http.Cookie {
	j.m.Lock()
	defer j.m.Unlock()
	return j.perURL[u.Host]
}
