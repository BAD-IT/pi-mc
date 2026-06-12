package ui

import (
	"encoding/json"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/earendil-works/pi-mc/internal/rpc"
)

type AppModel struct {
	leftPane   LeftPaneModel
	rightPane  RightPaneModel
	footer     FooterModel
	
	focusRight bool
	
	width  int
	height int

	rpcClient *rpc.PiRpcClient
}

type rpcEventMsg rpc.Event
type rpcErrorMsg struct{ err error }

func listenForRPCEvents(client *rpc.PiRpcClient) tea.Cmd {
	return func() tea.Msg {
		select {
		case ev, ok := <-client.Events():
			if !ok {
				return rpcErrorMsg{fmt.Errorf("rpc channel closed")}
			}
			return rpcEventMsg(ev)
		case err := <-client.Errors():
			return rpcErrorMsg{err}
		}
	}
}

func NewApp() AppModel {
	client, err := rpc.NewPiRpcClient("go", "run", "./cmd/mock-pi")
	if err != nil {
		log.Fatalf("Failed to initialize rpc client: %v", err)
	}

	return AppModel{
		leftPane:   NewLeftPane(),
		rightPane:  NewRightPane(),
		footer:     NewFooter(),
		focusRight: true,
		rpcClient:  client,
	}
}

func (m AppModel) Init() tea.Cmd {
	err := m.rpcClient.Start()
	if err != nil {
		log.Fatalf("Failed to start rpc client: %v", err)
	}
	
	return tea.Batch(
		m.rightPane.Init(),
		listenForRPCEvents(m.rpcClient),
	)
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "f10":
			m.rpcClient.Stop()
			return m, tea.Quit
		case "tab":
			m.focusRight = !m.focusRight
		case "esc":
			m.focusRight = true
		}
		
	case PromptMsg:
		m.rightPane.AddMessage("▸ " + string(msg))
		m.rpcClient.SendCommand(rpc.Command{
			Type: "prompt",
			Payload: map[string]string{"text": string(msg)},
		})
		
	case rpcEventMsg:
		if msg.Type == "message_update" {
			var update rpc.MessageUpdate
			if err := json.Unmarshal(msg.Payload, &update); err == nil {
				m.rightPane.AppendStream(update.Delta)
			}
		}
		// Also forward to leftPane so it can process state updates
		lp, _ := m.leftPane.Update(msg)
		m.leftPane = lp.(LeftPaneModel)

		cmds = append(cmds, listenForRPCEvents(m.rpcClient))
		
	case rpcErrorMsg:
		// ignore for now
		
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		paneHeight := m.height - 2
		if paneHeight < 0 {
			paneHeight = 0
		}
		
		leftWidth := int(float64(m.width) * 0.35)
		rightWidth := m.width - leftWidth
		
		m.leftPane.SetSize(leftWidth, paneHeight)
		m.rightPane.SetSize(rightWidth, paneHeight)
		
		f, _ := m.footer.Update(msg)
		m.footer = f.(FooterModel)
	}
	
	m.leftPane.SetActive(!m.focusRight)
	m.rightPane.SetActive(m.focusRight)
	
	if m.focusRight {
		rp, cmd := m.rightPane.Update(msg)
		m.rightPane = rp.(RightPaneModel)
		cmds = append(cmds, cmd)
	} else {
		lp, cmd := m.leftPane.Update(msg)
		m.leftPane = lp.(LeftPaneModel)
		cmds = append(cmds, cmd)
	}
	
	return m, tea.Batch(cmds...)
}

func (m AppModel) View() string {
	panes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.leftPane.View(),
		m.rightPane.View(),
	)
	
	return lipgloss.JoinVertical(
		lipgloss.Left,
		panes,
		m.footer.View(),
	)
}
