package tui

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

type model struct {
	Pages		[]page
	activePage	int
	guide		guide
	window 		window
	Themes		[]string
	activeTheme	int
};

type page struct {
	Name	string
	Content	string
	Order	int
};

type window struct {
	Width	int
	Height	int
};

type guide struct {
	keys		keyMap
	help		help.Model
	lastKey		string
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

var (
	highlightColor = lipgloss.AdaptiveColor{Light: "#8839EF", Dark: "#CBA6F7"};
	inactiveColor = lipgloss.AdaptiveColor{Light: "#313244", Dark: "#6C7086"};

	docStyle = lipgloss.NewStyle();
	activePageStyle = lipgloss.NewStyle().Bold(true).Foreground(highlightColor);
	inactivePageStyle = lipgloss.NewStyle().Foreground(inactiveColor);
	separatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF"));
	pageRowStyle = lipgloss.NewStyle().Padding(0, 1).Align(lipgloss.Center);
	pageWindowStyle = lipgloss.NewStyle().AlignVertical(lipgloss.Top);
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F38BA8"));
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

func check(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func getPages(dir string) ([]page, error) {
	files, err := os.ReadDir(dir)
	check(err, "Failed to read directory")

	pages := make(map[int]page)
	for _, file := range files {
		content, err := os.ReadFile(filepath.Join(dir, file.Name()))
		check(err, "Failed to read file")

		order, err := strconv.Atoi(file.Name()[:1])
		check(err, "Failed to convert order to integer")

		pages[order] = page{
			Name:    file.Name()[4 : len(file.Name())-3],
			Content: string(content),
			Order:   order,
		}
	}

	// Sort the pages by order
	var orderedPages []page
	for _, p := range pages {
		orderedPages = append(orderedPages, p)
	}
	sort.Slice(orderedPages, func(i, j int) bool {
		return orderedPages[i].Order < orderedPages[j].Order
	})

	return orderedPages, nil
}

func getThemes(dir string) ([]string, error) {
	files, err := os.ReadDir(dir);
	check(err, "Theme retrieval failed");

	var themes []string;
	for _, file := range files {
		themes = append(themes, file.Name());
	}
	sort.Strings(themes);

	return themes, nil;
}

func initialModel() model {
	pages, err := getPages("assets/pages");
	check(err, "Page retrieval failed");

	themes, err := getThemes("assets/themes");
	check(err, "Theme retrieval failed");

	return model{
		Pages:       pages,
		activePage:  0,
		guide: guide{
			keys:    keys,
			help:    help.New(),
			lastKey: "",
		},
		window: window{
			Width:  0,
			Height: 0,
		},
		Themes:      themes,
		activeTheme: 0,
	}
}

func updateStyleSizes(m model) {
	m.guide.help.Width = m.window.Width;
	docStyle = docStyle.Width(m.window.Width);
	pageRowStyle = pageRowStyle.Width(m.window.Width);
	pageWindowStyle = pageWindowStyle.Width(m.window.Width);
}

func renderMarkdown(m model, content string) string {
	themePath := filepath.Join("assets", "themes", m.Themes[m.activeTheme], "glamour.json");
	pageContentRenderer, _ := glamour.NewTermRenderer(glamour.WithPreservedNewLines(), glamour.WithWordWrap(m.window.Width - 4), glamour.WithStylePath(themePath));
	pageContent, _ := pageContentRenderer.Render(content);
	pageContent = pageWindowStyle.Render(pageContent);
	return pageContent;
}

func (m model) Init() tea.Cmd {
	return nil;
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handles msg types
	switch msg := msg.(type) {
		// Handle window size msgs
		case tea.WindowSizeMsg:
			m.window.Width = msg.Width;
			m.window.Height = msg.Height;
			updateStyleSizes(m);
			return m, tea.ClearScreen;

		// Handle key presses msgs
		case tea.KeyMsg:
			// Handles the actual key that was pressed
			switch {
				// Quit the program
				case key.Matches(msg, m.guide.keys.Quit):
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

	// Render page tabs
	var renderedPages []string;
	for i, t := range m.Pages {
		page := t.Name;
		var style lipgloss.Style;

		if (i == m.activePage) {
			style = activePageStyle;
		} else {
			style = inactivePageStyle;
		}

		// Render the active page
		page = style.Render(page);

		// Add a • to the end of the page if it's not the last page
		if (i != len(m.Pages) - 1) {
			page += separatorStyle.Render(" • ");
		}

		renderedPages = append(renderedPages, page);
	}
	pageTabs := pageRowStyle.Render(strings.Join(renderedPages, ""));
	
	
	// Render current page content
	pageContent := renderMarkdown(m, m.Pages[m.activePage].Content);

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
	p := tea.NewProgram(initialModel(), tea.WithAltScreen());
	if _, err := p.Run(); (err != nil) {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}