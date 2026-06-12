package ui

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/earendil-works/pi-mc/internal/config"
	"github.com/earendil-works/pi-mc/internal/rpc"
)

type EditorFinishedMsg struct {
	File string
	Err  error
}

type AppModel struct {
	leftPane   LeftPaneModel
	rightPane  RightPaneModel
	fileTree   FileTreeModel
	footer     FooterModel
	
	focusRight bool
	showTree   bool
	
	width  int
	height int

	rpcClient *rpc.PiRpcClient

	models       []string
	modelIndex   int
	thinkingLvls []string
	thinkingIdx  int
}

func NewApp() AppModel {
	cfg := config.LoadConfig()
	applyTheme(cfg.Theme)

	cmdStr := os.Getenv("PI_CMD")
	if cmdStr == "" {
		cmdStr = "go run ./cmd/mock-pi"
	}
	parts := strings.Fields(cmdStr)

	client, err := rpc.NewPiRpcClient(parts[0], parts[1:]...)
	if err != nil {
		log.Fatalf("Failed to initialize rpc client: %v", err)
	}

	return AppModel{
		leftPane:     NewLeftPane(),
		rightPane:    NewRightPane(),
		fileTree:     NewFileTree(),
		footer:       NewFooter(),
		focusRight:   true,
		showTree:     false,
		rpcClient:    client,
		models:       []string{"gemini-1.5-pro", "gemini-2.0-flash"},
		modelIndex:   0,
		thinkingLvls: []string{"none", "low", "high"},
		thinkingIdx:  0,
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
		if cmd := m.handleKey(msg); cmd != nil {
			return m, cmd
		}
		
	case EditorFinishedMsg:
		if msg.Err == nil {
			b, _ := os.ReadFile(msg.File)
			os.Remove(msg.File)
			content := strings.TrimSpace(string(b))
			if content != "" {
				m.rightPane.AddMessage("▸ (from editor)\n" + content)
				m.rpcClient.SendCommand(rpc.Command{
					Type: "prompt",
					Payload: map[string]string{"text": content},
				})
			}
		}
		
	case PromptMsg:
		msgStr := string(msg)
		if strings.HasPrefix(msgStr, "/") {
			parts := strings.Fields(msgStr)
			cmd := parts[0]
			switch cmd {
			case "/model":
				if len(parts) > 1 {
					m.rpcClient.SendCommand(rpc.Command{
						Type: "set_model", Payload: map[string]string{"model": parts[1]},
					})
					m.rightPane.AddMessage("⚙️ Model set to " + parts[1])
				}
			case "/thinking":
				if len(parts) > 1 {
					m.rpcClient.SendCommand(rpc.Command{
						Type: "set_thinking", Payload: map[string]string{"level": parts[1]},
					})
					m.rightPane.AddMessage("⚙️ Thinking level set to " + parts[1])
				}
			case "/compact":
				m.rpcClient.SendCommand(rpc.Command{Type: "compact"})
				m.rightPane.AddMessage("⚙️ Memory compacted")
			default:
				m.rightPane.AddMessage("❌ Unknown command: " + cmd)
			}
		} else {
			m.rightPane.AddMessage("▸ " + msgStr)
			m.rpcClient.SendCommand(rpc.Command{
				Type: "prompt",
				Payload: map[string]string{"text": msgStr},
			})
		}
		
	case rpcEventMsg:
		cmds = append(cmds, m.handleRPCEvent(msg))
		
	case rpcErrorMsg:
		m.rightPane.AddMessage("❌ RPC Error: " + msg.err.Error())
		m.rightPane.AddMessage("Backend disconnected. Press Ctrl+C to exit.")
		
	case tea.WindowSizeMsg:
		m.handleResize(msg)
	}
	
	m.leftPane.SetActive(!m.focusRight && !m.showTree)
	m.fileTree.SetActive(!m.focusRight && m.showTree)
	m.rightPane.SetActive(m.focusRight)
	
	if m.focusRight {
		rp, cmd := m.rightPane.Update(msg)
		m.rightPane = rp.(RightPaneModel)
		cmds = append(cmds, cmd)
	} else if m.showTree {
		ft, cmd := m.fileTree.Update(msg)
		m.fileTree = ft.(FileTreeModel)
		cmds = append(cmds, cmd)
	} else {
		lp, cmd := m.leftPane.Update(msg)
		m.leftPane = lp.(LeftPaneModel)
		cmds = append(cmds, cmd)
	}
	
	return m, tea.Batch(cmds...)
}

func (m *AppModel) handleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "f10":
		m.rpcClient.Stop()
		return tea.Quit
	case "tab":
		m.focusRight = !m.focusRight
	case "esc":
		m.focusRight = true
	case "f5":
		m.showTree = !m.showTree
	case "ctrl+g":
		f, _ := os.CreateTemp("", "pi-mc-*.txt")
		f.Close()
		
		editor := os.Getenv("VISUAL")
		if editor == "" {
			editor = os.Getenv("EDITOR")
		}
		if editor == "" {
			editor = "vim" // fallback
		}
		
		c := exec.Command(editor, f.Name())
		return tea.ExecProcess(c, func(err error) tea.Msg {
			return EditorFinishedMsg{File: f.Name(), Err: err}
		})
	case "f2":
		m.modelIndex = (m.modelIndex + 1) % len(m.models)
		m.rpcClient.SendCommand(rpc.Command{
			Type: "set_model",
			Payload: map[string]string{"model": m.models[m.modelIndex]},
		})
	case "shift+tab":
		m.thinkingIdx = (m.thinkingIdx + 1) % len(m.thinkingLvls)
		m.rpcClient.SendCommand(rpc.Command{
			Type: "set_thinking",
			Payload: map[string]string{"level": m.thinkingLvls[m.thinkingIdx]},
		})
	}
	return nil
}

func (m *AppModel) handleRPCEvent(msg rpcEventMsg) tea.Cmd {
	if msg.Type == "message_update" {
		var update rpc.MessageUpdate
		if err := json.Unmarshal(msg.Payload, &update); err == nil {
			m.rightPane.AppendStream(update.Delta)
		}
	}
	
	lp, _ := m.leftPane.Update(msg)
	m.leftPane = lp.(LeftPaneModel)

	return listenForRPCEvents(m.rpcClient)
}

func (m *AppModel) handleResize(msg tea.WindowSizeMsg) {
	m.width = msg.Width
	m.height = msg.Height
	
	// The footer is exactly 1 row tall, so the panes get the rest.
	paneHeight := m.height - 1
	if paneHeight < 0 {
		paneHeight = 0
	}
	
	leftWidth := int(float64(m.width) * 0.35)
	rightWidth := m.width - leftWidth
	
	m.leftPane.SetSize(leftWidth, paneHeight)
	m.fileTree.SetSize(leftWidth, paneHeight)
	m.rightPane.SetSize(rightWidth, paneHeight)
	
	f, _ := m.footer.Update(msg)
	m.footer = f.(FooterModel)
}

func (m AppModel) View() string {
	leftView := m.leftPane.View()
	if m.showTree {
		leftView = m.fileTree.View()
	}

	panes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftView,
		m.rightPane.View(),
	)
	
	footerText := m.footer.View()
	statusText := lipgloss.NewStyle().
		Background(CurrentTheme.FooterBg).
		Foreground(CurrentTheme.FooterFg).
		Render(fmt.Sprintf(" Model: %s | Thinking: %s ", m.models[m.modelIndex], m.thinkingLvls[m.thinkingIdx]))

	availWidth := m.width - lipgloss.Width(footerText) - lipgloss.Width(statusText)
	if availWidth < 0 { 
		availWidth = 0 
	}
	
	filler := lipgloss.NewStyle().Background(CurrentTheme.FooterBg).Render(strings.Repeat(" ", availWidth))
	
	fullFooter := footerText + filler + statusText
	
	return appBgStyle.Width(m.width).Height(m.height).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			panes,
			fullFooter,
		),
	)
}
