package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hunntr/w3env/internal/store"
	"github.com/spf13/cobra"
)

var (
	sBorder = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238"))

	sActiveBorder = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99"))

	sHeader = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		PaddingLeft(1)

	sSelected = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	sActiveProfile = lipgloss.NewStyle().
		Foreground(lipgloss.Color("84"))

	sKey = lipgloss.NewStyle().
		Foreground(lipgloss.Color("81"))

	sDim = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	sBold = lipgloss.NewStyle().Bold(true)

	sOK = lipgloss.NewStyle().
		Foreground(lipgloss.Color("84")).
		PaddingLeft(1)

	sErr = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		PaddingLeft(1)

	sHelp = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		PaddingLeft(1)
)

type panel int

const (
	panelProfiles panel = iota
	panelVars
)

type inputMode int

const (
	modeNormal inputMode = iota
	modeNewProfile
	modeSetKey
	modeSetValue
	modeEditValue
	modeRename
	modeConfirmDelete
)

type tui struct {
	s          *store.State
	profiles   []string
	pCursor    int
	varKeys    []string
	vCursor    int
	focus      panel
	mode       inputMode
	ti         textinput.Model
	pendingKey string
	width      int
	height     int
	status     string
	statusErr  bool
	reveal          bool
	activated       string
	deactivatedVars map[string]string
}

func newTUI(s *store.State) tui {
	ti := textinput.New()
	ti.CharLimit = 512
	m := tui{s: s, ti: ti}
	m.syncProfiles()
	return m
}

func (m *tui) syncProfiles() {
	m.profiles = m.s.ProfileNames()
	if m.pCursor >= len(m.profiles) && len(m.profiles) > 0 {
		m.pCursor = len(m.profiles) - 1
	}
	m.syncVars()
}

func (m *tui) syncVars() {
	m.varKeys = nil
	if len(m.profiles) == 0 {
		return
	}
	p := m.s.Profiles[m.profiles[m.pCursor]]
	for k := range p.Vars {
		m.varKeys = append(m.varKeys, k)
	}
	sort.Strings(m.varKeys)
	if m.vCursor >= len(m.varKeys) && len(m.varKeys) > 0 {
		m.vCursor = len(m.varKeys) - 1
	}
}

func (m tui) Init() tea.Cmd { return textinput.Blink }


func (m tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case tea.KeyMsg:
		if m.mode != modeNormal {
			return m.handleInput(msg)
		}
		return m.handleNav(msg)
	}
	return m, nil
}

func (m tui) handleNav(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.status = ""

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "tab":
		if m.focus == panelProfiles {
			m.focus = panelVars
		} else {
			m.focus = panelProfiles
		}

	case "up", "k":
		if m.focus == panelProfiles {
			if m.pCursor > 0 {
				m.pCursor--
				m.vCursor = 0
				m.syncVars()
			}
		} else {
			if m.vCursor > 0 {
				m.vCursor--
			}
		}

	case "down", "j":
		if m.focus == panelProfiles {
			if m.pCursor < len(m.profiles)-1 {
				m.pCursor++
				m.vCursor = 0
				m.syncVars()
			}
		} else {
			if m.vCursor < len(m.varKeys)-1 {
				m.vCursor++
			}
		}

	case "enter":
		if m.focus == panelProfiles && len(m.profiles) > 0 {
			name := m.profiles[m.pCursor]
			m.s.Active = name
			m.activated = name
			if err := m.s.Save(); err != nil {
				m.status, m.statusErr = err.Error(), true
			} else {
				m.status, m.statusErr = fmt.Sprintf("switched to %q", name), false
			}
		}

	case "n":
		m.startInput(modeNewProfile, "", "profile name")

	case "r":
		if len(m.profiles) == 0 {
			break
		}
		if m.focus == panelProfiles {
			m.startInput(modeRename, m.profiles[m.pCursor], "new name")
		}

	case "s":
		if len(m.profiles) == 0 {
			m.status, m.statusErr = "create a profile first (n)", true
			break
		}
		m.focus = panelProfiles
		m.startInput(modeSetKey, "", "variable name (e.g. RPC_URL)")

	case "e":
		if m.focus == panelVars && len(m.varKeys) > 0 {
			key := m.varKeys[m.vCursor]
			val := m.s.Profiles[m.profiles[m.pCursor]].Vars[key]
			m.pendingKey = key
			m.startInput(modeEditValue, val, key)
		}

	case "d", "backspace":
		if m.focus == panelProfiles && len(m.profiles) > 0 {
			m.startInput(modeConfirmDelete, "", fmt.Sprintf("delete %q? (y/N)", m.profiles[m.pCursor]))
		} else if m.focus == panelVars && len(m.varKeys) > 0 {
			name := m.profiles[m.pCursor]
			key := m.varKeys[m.vCursor]
			p := m.s.Profiles[name]
			delete(p.Vars, key)
			m.s.Profiles[name] = p
			m.s.Save()
			m.syncVars()
			m.status, m.statusErr = fmt.Sprintf("removed %s", key), false
		}

	case "v":
		m.reveal = !m.reveal
	}

	return m, nil
}

func (m *tui) startInput(mode inputMode, value, placeholder string) {
	m.mode = mode
	m.ti.SetValue(value)
	m.ti.Placeholder = placeholder
	m.ti.CursorEnd()
	m.ti.Focus()
}

func (m tui) handleInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeNormal
		m.ti.Blur()
		return m, nil

	case "enter":
		raw := m.ti.Value()
		val := strings.TrimSpace(raw)

		switch m.mode {
		case modeNewProfile:
			if val == "" {
				break
			}
			if _, exists := m.s.Profiles[val]; exists {
				m.status, m.statusErr = fmt.Sprintf("%q already exists", val), true
			} else {
				m.s.Profiles[val] = store.Profile{Vars: make(map[string]string)}
				m.s.Save()
				m.syncProfiles()
				for i, n := range m.profiles {
					if n == val {
						m.pCursor = i
						break
					}
				}
				m.syncVars()
				m.status, m.statusErr = fmt.Sprintf("created %q", val), false
			}
			m.mode = modeNormal
			m.ti.Blur()

		case modeSetKey:
			if val == "" {
				break
			}
			m.pendingKey = val
			m.mode = modeSetValue
			m.ti.SetValue("")
			m.ti.Placeholder = fmt.Sprintf("value for %s", val)
			return m, nil

		case modeSetValue, modeEditValue:
			if m.pendingKey != "" && len(m.profiles) > 0 {
				name := m.profiles[m.pCursor]
				p := m.s.Profiles[name]
				if p.Vars == nil {
					p.Vars = make(map[string]string)
				}
				p.Vars[m.pendingKey] = raw // keep raw value (no trim)
				m.s.Profiles[name] = p
				m.s.Save()
				m.syncVars()
				for i, k := range m.varKeys {
					if k == m.pendingKey {
						m.vCursor = i
						break
					}
				}
				m.status, m.statusErr = fmt.Sprintf("set %s", m.pendingKey), false
				m.pendingKey = ""
			}
			m.mode = modeNormal
			m.ti.Blur()

		case modeRename:
			if val == "" || len(m.profiles) == 0 {
				m.mode = modeNormal
				m.ti.Blur()
				break
			}
			oldName := m.profiles[m.pCursor]
			if val == oldName {
				m.mode = modeNormal
				m.ti.Blur()
				break
			}
			if _, exists := m.s.Profiles[val]; exists {
				m.status, m.statusErr = fmt.Sprintf("%q already exists", val), true
			} else {
				p := m.s.Profiles[oldName]
				m.s.Profiles[val] = p
				delete(m.s.Profiles, oldName)
				if m.s.Active == oldName {
					m.s.Active = val
				}
				m.s.Save()
				m.syncProfiles()
				for i, n := range m.profiles {
					if n == val {
						m.pCursor = i
						break
					}
				}
				m.status, m.statusErr = fmt.Sprintf("renamed -> %q", val), false
			}
			m.mode = modeNormal
			m.ti.Blur()

		case modeConfirmDelete:
			if strings.ToLower(val) == "y" && len(m.profiles) > 0 {
				name := m.profiles[m.pCursor]
				if m.s.Active == name {
					p := m.s.Profiles[name]
					vars := make(map[string]string, len(p.Vars))
					for k, v := range p.Vars {
						vars[k] = v
					}
					m.deactivatedVars = vars
					m.s.Active = ""
				}
				delete(m.s.Profiles, name)
				m.s.Save()
				m.syncProfiles()
				m.status, m.statusErr = fmt.Sprintf("deleted %q", name), false
			}
			m.mode = modeNormal
			m.ti.Blur()
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.ti, cmd = m.ti.Update(msg)
	return m, cmd
}


func (m tui) View() string {
	if m.width == 0 {
		return ""
	}

	leftW := 28
	rightW := m.width - leftW - 7
	if rightW < 20 {
		rightW = 20
	}
	panelH := m.height - 5
	if panelH < 4 {
		panelH = 4
	}

	lBorder := sBorder.Width(leftW).Height(panelH)
	rBorder := sBorder.Width(rightW).Height(panelH)
	switch {
	case m.focus == panelProfiles, m.mode == modeNewProfile,
		m.mode == modeSetKey, m.mode == modeSetValue,
		m.mode == modeRename, m.mode == modeConfirmDelete:
		lBorder = sActiveBorder.Width(leftW).Height(panelH)
	case m.focus == panelVars, m.mode == modeEditValue:
		rBorder = sActiveBorder.Width(rightW).Height(panelH)
	}

	panels := lipgloss.JoinHorizontal(lipgloss.Top,
		lBorder.Render(m.profilePanel(leftW, panelH)),
		" ",
		rBorder.Render(m.varPanel(rightW, panelH)),
	)

	var bottom string
	switch {
	case m.mode != modeNormal:
		bottom = m.inputBar()
	case m.status != "":
		if m.statusErr {
			bottom = sErr.Render("✗ " + m.status)
		} else {
			bottom = sOK.Render("✓ " + m.status)
		}
	default:
		bottom = m.helpBar()
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		sHeader.Render("w3env"),
		panels,
		bottom,
	)
}

func (m tui) profilePanel(_ , h int) string {
	title := sBold.Render("profiles")
	if len(m.profiles) == 0 {
		return strings.Join([]string{
			title, "",
			sDim.Render("  no profiles"),
			sDim.Render("  press n to create"),
		}, "\n")
	}

	visible := h - 2
	start := 0
	if m.pCursor >= visible {
		start = m.pCursor - visible + 1
	}

	lines := []string{title, ""}
	for i := start; i < len(m.profiles) && i-start < visible; i++ {
		n := m.profiles[i]
		isActive := n == m.s.Active
		isSel := i == m.pCursor

		prefix := "  "
		name := n
		if isActive {
			name = sActiveProfile.Render(n)
			prefix = sActiveProfile.Render("● ")
		}
		if isSel {
			prefix = sSelected.Render("▶ ")
			if isActive {
				name = sActiveProfile.Render(n)
			} else {
				name = sSelected.Render(n)
			}
		}
		count := sDim.Render(fmt.Sprintf(" (%d)", len(m.s.Profiles[n].Vars)))
		lines = append(lines, prefix+name+count)
	}
	return strings.Join(lines, "\n")
}

func (m tui) varPanel(w, h int) string {
	if len(m.profiles) == 0 {
		return sDim.Render("no profile selected")
	}
	name := m.profiles[m.pCursor]
	p := m.s.Profiles[name]

	title := sBold.Render(name)
	if name == m.s.Active {
		title += " " + sActiveProfile.Render("(active)")
	}
	if m.reveal {
		title += " " + sDim.Render("(revealed)")
	}

	if len(m.varKeys) == 0 {
		return strings.Join([]string{
			title, "",
			sDim.Render("  empty"),
			sDim.Render("  press s to add a variable"),
		}, "\n")
	}

	maxK := 0
	for _, k := range m.varKeys {
		if len(k) > maxK {
			maxK = len(k)
		}
	}
	if maxK > 22 {
		maxK = 22
	}

	visible := h - 2
	start := 0
	if m.vCursor >= visible {
		start = m.vCursor - visible + 1
	}

	lines := []string{title, ""}
	for i := start; i < len(m.varKeys) && i-start < visible; i++ {
		k := m.varKeys[i]
		v := p.Vars[k]
		if !m.reveal && isSensitive(k) {
			v = "******"
		}
		maxV := w - maxK - 7
		if maxV > 0 && len(v) > maxV {
			v = v[:maxV-1] + "…"
		}
		pad := strings.Repeat(" ", maxK-len(k))
		isSel := i == m.vCursor && m.focus == panelVars

		var line string
		if isSel {
			line = "  " + sSelected.Render(k) + pad + sDim.Render(" = ") + sSelected.Render(v)
		} else {
			line = "  " + sKey.Render(k) + pad + sDim.Render(" = ") + v
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (m tui) inputBar() string {
	var label string
	switch m.mode {
	case modeNewProfile:
		label = "new profile"
	case modeSetKey:
		label = "key"
	case modeSetValue:
		label = fmt.Sprintf("value for %s", m.pendingKey)
	case modeEditValue:
		label = fmt.Sprintf("edit %s", m.pendingKey)
	case modeRename:
		label = "rename"
	case modeConfirmDelete:
		if len(m.profiles) > 0 {
			label = sErr.Render(fmt.Sprintf("delete %q?", m.profiles[m.pCursor]))
		}
	}
	return "  " + sBold.Render(label+"  ") + m.ti.View() + sDim.Render("  esc cancel")
}

func (m tui) helpBar() string {
	if m.focus == panelProfiles {
		return sHelp.Render("↑↓/jk nav  enter activate  n new  r rename  d delete  tab->vars  q quit")
	}
	return sHelp.Render("↑↓/jk nav  s set var  e edit  d delete  v reveal  tab->profiles  q quit")
}


func runTUI(cmd *cobra.Command, args []string) {
	s, err := store.Load()
	if err != nil {
		fatal(err.Error())
	}
	m := newTUI(s)

	opts := []tea.ProgramOption{tea.WithAltScreen()}
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err == nil {
		defer tty.Close()
		opts = append(opts, tea.WithInput(tty), tea.WithOutput(tty))
	} else {
		opts = append(opts, tea.WithOutput(os.Stderr))
	}

	p := tea.NewProgram(m, opts...)
	result, err := p.Run()
	if err != nil {
		fatal(err.Error())
	}
	final := result.(tui)
	if final.deactivatedVars != nil {
		writeDeactivationVars(final.deactivatedVars)
	} else {
		activeName := final.activated
		if activeName == "" {
			activeName = final.s.Active
		}
		if activeName != "" {
			fresh, err := store.Load()
			if err == nil {
				writeActivation(fresh, activeName)
			}
		}
	}
}

var tuiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Open the interactive TUI (also runs with no subcommand)",
	Run:   runTUI,
}

func init() {
	rootCmd.AddCommand(tuiCmd)
	rootCmd.Run = runTUI
}
