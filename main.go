package main

import (
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const url = "https://www.google.com"

type UrlStatus struct {
	URL        string
	StatusCode int
	Up         bool
}

type model struct {
	AllURLs    []UrlStatus
	Loading    bool
	InsertMode bool
	NewUrl     string
}

func checkInterval() tea.Cmd {
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return t
	})
}

func (m model) Init() tea.Cmd {
	return checkInterval()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "/":
			if !m.InsertMode {
				m.InsertMode = true
			} else {
				m.NewUrl += "/"
			}
			return m, nil
		case "enter":
			m.InsertMode = false
			m.AllURLs = append(m.AllURLs, struct {
				URL        string
				StatusCode int
				Up         bool
			}{URL: m.NewUrl, StatusCode: 0, Up: false})
			m.NewUrl = ""
			return m, nil
		case "esc":
			m.InsertMode = false
			m.NewUrl = ""
			return m, nil
		case "backspace":
			if m.InsertMode {
				m.NewUrl = m.NewUrl[:len(m.NewUrl)-1]
			}
			return m, nil
		default:
			if m.InsertMode {
				m.NewUrl += msg.String()
			}
		}
	case time.Time:
		for i, url := range m.AllURLs {
			statusCode, up := checkURL(url.URL)
			m.AllURLs[i].StatusCode = statusCode
			m.AllURLs[i].Up = up
		}
		return m, checkInterval()
	default:
		return m, nil
	}

	return m, nil
}

func checkURL(url string) (int, bool) {
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return 0, false
	}
	return resp.StatusCode, true
}

func (m model) View() string {
	s := "URL Monitor\n"
	if m.InsertMode {
		s += fmt.Sprintf("Insert mode: %s", m.NewUrl)
	} else {
		s += fmt.Sprintf("Monitor mode")
		for _, url := range m.AllURLs {
			isUp := "🔴"
			if url.Up {
				isUp = "🟢"
			}
			s += fmt.Sprintf("\n <%s> %s - %d", isUp, url.URL, url.StatusCode)
		}
	}

	s += "\n\n"
	if m.InsertMode {
		s += "esc:cancel  enter:save"
	} else {
		s += "Press / to insert new url"
	}
	s += "  ctrl+c:quit"
	return s
}

func main() {
	p := tea.NewProgram(model{})
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
	}
}
