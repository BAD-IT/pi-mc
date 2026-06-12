package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/earendil-works/pi-mc/internal/rpc"
)

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
