package src

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	glamour "github.com/charmbracelet/glamour"
	lipgloss "github.com/charmbracelet/lipgloss"
	help "github.com/charmbracelet/bubbles/help"
	key "github.com/charmbracelet/bubbles/key"
);

type guide struct {
	keys		keyMap
	help		help.Model
	lastKey		string
	quitting	bool
};

type model struct {
	Pages		[]page
	activePage	int
	guide		guide
};

type page struct {
	Name	string
	Content	string
	Order	int
};

type keyMap struct {
	Left		key.Binding
	Right		key.Binding
	Help 		key.Binding
	Quit 		key.Binding
};

var keys = keyMap{
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "Previous page"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l", "tab"),
		key.WithHelp("→/l", "Next page"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "Toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "Quit"),
	),
};

// Box unicode chars: ┌ ┐ └ ┘ ─ │ ┬ ┴ ├ ┤ ┼
var (
	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2).Width(80);
	highlightColor = lipgloss.AdaptiveColor{Light: "#8839EF", Dark: "#CBA6F7"};
	inactiveColor = lipgloss.AdaptiveColor{Light: "#313244", Dark: "#6C7086"};
	activePageStyle = lipgloss.NewStyle().Bold(true).Foreground(highlightColor);
	inactivePage = lipgloss.NewStyle().Foreground(inactiveColor);
	separatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF"));
	pageRowBorder = lipgloss.RoundedBorder();
	pageRowStyle = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(0, 1).Align(lipgloss.Center).Width(74);
	pageContentRenderer, _ = glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(74), glamour.WithPreservedNewLines());
	pageWindowStyle = lipgloss.NewStyle().BorderForeground(highlightColor).Border(lipgloss.RoundedBorder()).UnsetBorderTop().Width(74).AlignVertical(lipgloss.Top);
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F38BA8"));
	goodbyeStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Foreground(lipgloss.Color("#F38BA8")).Padding(0, 1);
	goodbyeMessage = "Thanks for checking out my portfolio TUI! :D";
);

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Help};
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right},
		{k.Quit, k.Help},
	};
}

func getPages(dir string) ([]page, error) {
	files, err := os.ReadDir(dir);
	if (err != nil) {
		log.Fatal(err);
	}
	
	pages := make(map[int]page);
	for _, file := range files {
		content, err := os.ReadFile(filepath.Join(dir, file.Name()));
		if (err != nil) {
			log.Fatal(err);
		}
		order, err := strconv.Atoi(file.Name()[:1])
		if err != nil {
			log.Fatal(err)
		}
		pages[order] = page{
			Name: file.Name()[4:len(file.Name()) - 3],
			Content: string(content),
			Order: order,
		};
	}

	// Sort the pages by order
	var orderedPages []page;
	for _, p := range pages {
		orderedPages = append(orderedPages, p);
	}
	sort.Slice(orderedPages, func(i, j int) bool {
		return orderedPages[i].Order < orderedPages[j].Order;
	});

	return orderedPages, nil;
}

func initialModel() model {
	pages, err := getPages("pages");
	if (err != nil) {
		log.Fatal(err);
	}

	return model{
		Pages:      pages,
		activePage: 0,
		guide: 		guide{
			keys: keys,
			help: help.New(),
			lastKey: "",
			quitting: false,
		},
	};
}

func (m model) Init() tea.Cmd {
	return nil;
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handles msg types
	switch msg := msg.(type) {
		// Handle window size msgs
		case tea.WindowSizeMsg:
			m.guide.help.Width = msg.Width;

		// Handle key presses msgs
		case tea.KeyMsg:
			// Handles the actual key that was pressed
			switch {
				// Quit the program
				case key.Matches(msg, m.guide.keys.Quit):
					m.guide.quitting = true;
					return m, tea.Quit;
				
				// Move to the next tab on the right
				case key.Matches(msg, m.guide.keys.Right):
					m.guide.lastKey = "→";
					if (m.activePage + 1 >= len(m.Pages)) {
						m.activePage = 0;
					} else {
						m.activePage++;
					}

				// Move to the next tab on the left
				case key.Matches(msg, m.guide.keys.Left):
					m.guide.lastKey = "←";
					if (m.activePage == 0) {
						m.activePage = len(m.Pages) - 1;
					} else {
						m.activePage--;
					}
				
				// Toggle short/full help
				case key.Matches(msg, m.guide.keys.Help):
					m.guide.help.ShowAll = !m.guide.help.ShowAll;
			}
	}

	// Return the updated model
	return m, nil
}

func (m model) View() string {
	// Initialize the main view string builder
	doc := strings.Builder{};

	// If the user is quitting, return the goodbye message
	if (m.guide.quitting) {
		doc.WriteString(goodbyeStyle.Render(goodbyeMessage));
		return docStyle.Render(doc.String());
	}

	// Render page tabs
	var renderedPages []string;
	for i, t := range m.Pages {
		page := t.Name;
		var style lipgloss.Style;

		if (i == m.activePage) {
			style = activePageStyle;
		} else {
			style = inactivePage;
		}

		// Render the active page
		page = style.Render(page);

		// Add a • to the end of the page if it's not the last page
		if (i != len(m.Pages) - 1) {
			page += separatorStyle.Render(" • ");
		}

		renderedPages = append(renderedPages, page);
	}
	pageRowBorder.BottomLeft = "├";
	pageRowBorder.BottomRight = "┤";
	pageRowStyle = pageRowStyle.Border(pageRowBorder, true);
	pageTabs := pageRowStyle.Render(strings.Join(renderedPages, ""));
	
	
	// Render current page content
	pageContent, _ := pageContentRenderer.Render(m.Pages[m.activePage].Content);
	pageContent = pageWindowStyle.Render(pageContent);

	// Render help
	helpContent := m.guide.help.View(m.guide.keys);
	helpContent = helpStyle.Render(helpContent);

	// Render the full view
	doc.WriteString(pageTabs);
	doc.WriteString("\n");
	doc.WriteString(pageContent);
	doc.WriteString("\n");
	doc.WriteString(helpContent);
	return docStyle.Render(doc.String());
}

func StartTUI() {
	p := tea.NewProgram(initialModel());
	if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}