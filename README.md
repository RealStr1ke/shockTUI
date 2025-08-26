# shockTUI

A portfolio TUI built in Go, showcasing my personal information, projects, and achievements in a silly interactive and themeable format.

## Features

- **Interactive Navigation**: Navigate through different portfolio sections using keyboard shortcuts
- **Multiple Themes**: Switch between various pre-built themes (Dark, Dracula, Pink, Tokyo Night)
- **Markdown Rendering**: Beautiful markdown rendering with syntax highlighting
- **Responsive Design**: Adapts to terminal window size
- **Smooth Scrolling**: Viewport-based content scrolling
- **Help System**: Built-in help menu with keyboard shortcuts

## Screenshots

*Navigate through pages using arrow keys or `h`/`l`, switch themes with `t`, and toggle help with `?`*

## Usage

### Build from Source
```bash
git clone https://github.com/RealStr1ke/shockTUI.git
cd shockTUI
go mod download
go build -o shocktui
./shocktui
```

### Quick Run
```bash
go run main.go
```

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `←` / `h` | Previous page |
| `→` / `l` / `Tab` | Next page |
| `t` / `Ctrl+Tab` | Change theme |
| `?` | Toggle help |
| `q` / `Esc` / `Ctrl+C` | Quit |

## Dependencies

This project uses the following excellent libraries:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style and layout
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [go-toml](https://github.com/pelletier/go-toml) - TOML configuration parsing

## Customization

### Adding New Pages
1. Create a new markdown file in `assets/pages/` with the format `# - title.md`
2. The application will automatically detect and load the new page

### Creating Custom Themes
1. Create a new directory in `assets/themes/`
2. Add a `glamour.json` file for markdown styling
3. Add a `colors.toml` file for UI color configuration
4. The theme will be automatically available in the theme switcher

### Color Configuration
Edit the `colors.toml` file in any theme directory:
```toml
TitleColor = "#FF0000"
InactiveTabColor = "#666666"
ActiveTabColor = "#FFFFFF"
TabTitleColor = "#00FF00"
TabBorderColor = "#0000FF"
ThemeTitleColor = "#FFFF00"
PageProgressColor = "#FF00FF"
```

## Development

### Running in Development
```bash
go run main.go
```

### Building
```bash
go build -o shocktui
```
## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

*Built with ❤️ using Go and the Charm libraries*