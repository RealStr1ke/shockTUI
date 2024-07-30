package tui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	help "github.com/charmbracelet/bubbles/help"
	key "github.com/charmbracelet/bubbles/key"
	viewport "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	glamour "github.com/charmbracelet/glamour"
	lipgloss "github.com/charmbracelet/lipgloss"
);

type model struct {
	Pages			[]page
	activePage		int

	Themes			[]string
	activeTheme		int

	guide			guide
	width, height	int
	pageview 		pageview
};

type page struct {
	Name	string
	Content	string
	Order	int
};

type pageview struct {
	viewport viewport.Model
	ready	bool
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
	Theme 		key.Binding
	Quit 		key.Binding
};

var (
	highlightColor = lipgloss.AdaptiveColor{Light: "#8839EF", Dark: "#CBA6F7"};
	inactiveColor = lipgloss.AdaptiveColor{Light: "#313244", Dark: "#6C7086"};

	docStyle = lipgloss.NewStyle();
	pageRowTitleStyle = lipgloss.NewStyle().Align(lipgloss.Left).Foreground(lipgloss.Color("#F38BA8")).Bold(true).PaddingLeft(2);
	themeTitleStyle = lipgloss.NewStyle().Align(lipgloss.Right).Foreground(lipgloss.Color("#F38BA8")).Bold(true);
	pageRowStyle = lipgloss.NewStyle().Align(lipgloss.Left);
	activePageStyle = lipgloss.NewStyle().Bold(true).Foreground(highlightColor)
	inactivePageStyle = lipgloss.NewStyle().Foreground(inactiveColor)
	separatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF"));
	helpStyle = lipgloss.NewStyle();

	pageRowTitle = "str1ke  ";
);

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
	Theme: key.NewBinding(
		key.WithKeys("t", "ctrl+tab"),
		key.WithHelp("t", "Change theme"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "Quit"),
	),
};

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Help, k.Theme};
}
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right},
		{k.Quit, k.Help},
		{k.Theme},
	};
}

func (p pageview) headerView(title string) string {
	title = fmt.Sprintf("──[ %s ]──", lipgloss.NewStyle().Foreground(lipgloss.Color("#FAB387")).Bold(true).Render(title));
	line := strings.Repeat("─", max(0, p.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (p pageview) footerView(themeTitle string) string {
	info := fmt.Sprintf("──[ %3.f%% ]──", p.viewport.ScrollPercent()*100)
	theme := fmt.Sprintf("──[ %s ]──", themeTitleStyle.Render(themeTitle));
	line := strings.Repeat("─", max(0, p.viewport.Width - lipgloss.Width(info) - lipgloss.Width(theme)))
	return lipgloss.JoinHorizontal(lipgloss.Center, theme, line, info)
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
		Themes:      themes,
		activeTheme: 0,

		guide: 		 guide{
			keys:    keys,
			help:    help.New(),
			lastKey: "",
		},

		width:		 0,
		height:		 0,
	}
}

func (m model) updateStyleSizes() {
	m.guide.help.Width = m.width;
	docStyle = docStyle.Width(m.width);
	docStyle = docStyle.Height(m.height);
	pageRowStyle = pageRowStyle.Width(m.width - len(m.Themes[m.activeTheme]));
}

func renderMarkdown(m model, content string) string {
	themePath := filepath.Join("assets", "themes", m.Themes[m.activeTheme], "glamour.json");
	pageContentRenderer, _ := glamour.NewTermRenderer(glamour.WithPreservedNewLines(), glamour.WithWordWrap(m.width - 4), glamour.WithStylePath(themePath));
	pageContent, _ := pageContentRenderer.Render(content);
	return pageContent;
}

func (m model) Init() tea.Cmd {
	return nil;
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// Handles msg types
	switch msg := msg.(type) {
		// Handle window size msgs
		case tea.WindowSizeMsg:
			m.width = msg.Width;
			m.height = msg.Height;
			m.updateStyleSizes();
			headerHeight := lipgloss.Height(m.pageview.headerView(m.Pages[m.activePage].Name));
			footerHeight := lipgloss.Height(m.pageview.footerView(m.Themes[m.activeTheme]));
			verticalMarginHeight := headerHeight + footerHeight + 3;
			if (!m.pageview.ready) {
				m.pageview.viewport = viewport.New(msg.Width, msg.Height - verticalMarginHeight);
				m.pageview.viewport.YPosition = headerHeight;
				m.pageview.viewport.HighPerformanceRendering = false;
				m.pageview.ready = true;
				m.pageview.viewport.YPosition = headerHeight + 1;
			} else {
				m.pageview.viewport.Width = msg.Width;
				m.pageview.viewport.Height = msg.Height - verticalMarginHeight;
			}

			m.pageview.viewport.SetContent(renderMarkdown(m, m.Pages[m.activePage].Content));
			// cmds = append(cmds, viewport.Sync(m.pageview.viewport));
			m.pageview.viewport, cmd = m.pageview.viewport.Update(msg);
			cmds = append(cmds, cmd);
			cmds = append(cmds, tea.ClearScreen)
			return m, tea.Batch(cmds...);

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
				
				// Change theme
				case key.Matches(msg, m.guide.keys.Theme):
					m.guide.lastKey = "t";
					if (m.activeTheme + 1 >= len(m.Themes)) {
						m.activeTheme = 0;
					} else {
						m.activeTheme++;
					}
			}

			m.pageview.viewport.SetContent(renderMarkdown(m, m.Pages[m.activePage].Content));
			// cmds = append(cmds, viewport.Sync(m.pageview.viewport));
			m.pageview.viewport, cmd = m.pageview.viewport.Update(msg);
			cmds = append(cmds, cmd);
			return m, tea.Batch(cmds...);
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
	pageTabs := strings.Join(renderedPages, "");
	pageTabsTitle := pageRowTitleStyle.Render(pageRowTitle);
	
	// Render current page content
	// viewport.Sync(m.pageview.viewport);
	pageContent := m.pageview.headerView(m.Pages[m.activePage].Name) + "\n" + m.pageview.viewport.View() + "\n" + m.pageview.footerView(m.Themes[m.activeTheme]);

	// Render help
	helpContent := m.guide.help.View(m.guide.keys);
	helpContent = helpStyle.Render(helpContent);

	// Render the full view
	doc.WriteString(pageRowStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, pageTabsTitle, pageTabs)));
	doc.WriteString("\n");
	doc.WriteString(pageContent);
	doc.WriteString("\n");
	doc.WriteString(helpContent);
	return docStyle.Render(doc.String());
}

func StartTUI() {
	p := tea.NewProgram(
		initialModel(), 
		tea.WithAltScreen(), 
		tea.WithMouseCellMotion(),
	);
	if _, err := p.Run(); (err != nil) {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}