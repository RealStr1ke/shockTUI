package src;

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
);

type model struct {
	Tabs       []string
	TabContent []string
	activeTab  int
}

func getPages(dir string) (map[string]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err;
	}
	
	pages := make(map[string]string)
	for _, file := range files {
		// pages = append(pages, file.Name()[:len(file.Name())-3]);
		content, err := os.ReadFile(filepath.Join(dir, file.Name()));
		if err != nil {
			return nil, err;
		}
		pages[file.Name()[:len(file.Name())-3]] = string(content);
	}

	return pages, nil;
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}

	return m, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

// View returns the rendered view of the model.
// It generates a string representation of the tabs and the active tab's content.
// The rendered view includes the tabs and the content of the active tab.
// The tabs are styled based on their active/inactive state.
// The active tab is highlighted with the activeTabStyle, while the inactive tabs are styled with the inactiveTabStyle.
// The rendered view is constructed by joining the rendered tabs horizontally using lipgloss.JoinHorizontal.
// The active tab's content is retrieved from the TabContent map using the activeTab index.
// The width of the active tab's content is adjusted to fit the window by subtracting the horizontal frame size of the windowStyle.
// The final rendered view is returned as a string.
func (m model) View() string {
	// Create a strings.Builder to store the rendered view
	doc := strings.Builder{}

	// Create an empty slice to store the rendered tabs
	var renderedTabs []string

	// Iterate over the tabs in the model
	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab

		// Determine the style based on the tab's active/inactive state
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}

		// Get the border of the style
		border, _, _, _, _ := style.GetBorder()

		// Adjust the border based on the tab's position and active/inactive state
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}

		// Apply the adjusted border to the style
		style = style.Border(border)

		// Render the tab using the style and append it to the renderedTabs slice
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	// Join the rendered tabs horizontally using lipgloss.JoinHorizontal
	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	// Write the row to the strings.Builder
	doc.WriteString(row)
	doc.WriteString("\n")

	// Retrieve the content of the active tab from the TabContent map
	activeTabContent := m.TabContent[m.activeTab]

	// Adjust the width of the active tab's content to fit the window
	adjustedContent := windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(activeTabContent)

	// Write the adjusted content to the strings.Builder using the docStyle
	doc.WriteString(docStyle.Render(adjustedContent))

	// Return the final rendered view as a string
	return doc.String()
}

func StartTUI() {
	pages, err := getPages("pages");
	if err != nil {
		panic(err);
	}
	tabs := make([]string, 0, len(pages));
	tabContent := make([]string, 0, len(pages));
	for page, content := range pages {
		tabs = append(tabs, page);
		tabContent = append(tabContent, content);
	}
	
	m := model{Tabs: tabs, TabContent: tabContent}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}