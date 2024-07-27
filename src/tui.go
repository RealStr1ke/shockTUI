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
);


type model struct {
	Pages		[]page
	activePage	int
}

type page struct {
	Name	string
	Content	string
	Order	int
}

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
)

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
		Pages:       pages,
		activePage:  0,
	};
}

func (m model) Init() tea.Cmd {
	return nil;
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handles msg types
	switch msg := msg.(type) {
		// Handle key presses msgs
		case tea.KeyMsg:
			// Handles the actual key that was pressed
			switch msg.String() {
				// Quit the program
				case "ctrl+c", "q":
					return m, tea.Quit;
				
				// Move to the next tab on the right
				case "right", "l", "n", "tab":
					if (m.activePage + 1 >= len(m.Pages)) {
						m.activePage = 0;
					} else {
						m.activePage++;
					}

				// Move to the next tab on the left
				case "left", "h", "p":
					if (m.activePage == 0) {
						m.activePage = len(m.Pages) - 1;
					} else {
						m.activePage--;
					}
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
	
	// Render the full view
	doc.WriteString(pageTabs);
	doc.WriteString("\n");
	doc.WriteString(pageContent);
	return docStyle.Render(doc.String());
}

func StartTUI() {
	p := tea.NewProgram(initialModel());
	if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}