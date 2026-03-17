package render

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aymantoumi/portlens/internal/conflict"
	"github.com/aymantoumi/portlens/internal/registry"
	"github.com/aymantoumi/portlens/internal/types"
)

func PrintConflictBlock(s conflict.Suggestion) {
	content := fmt.Sprintf("  conflict  ·  port %d\n  %s\n", s.Port, s.Description)
	if s.Alternative != 0 {
		content += fmt.Sprintf("  move to %s  —  portlens convention\n", styleGreen.Render(fmt.Sprintf("%d", s.Alternative)))
	}
	if s.FixCommand != "" {
		content += fmt.Sprintf("  fix: %s", stylePrimary.Render(s.FixCommand))
	}
	fmt.Println(content)
}

func PrintFreeResult(result *types.FreeResult) {
	fmt.Printf(" %s  %s\n", styleGreen.Render("▣"), styleWhite.Bold(true).Render(fmt.Sprintf("portlens free %d-%d", result.Range.Min, result.Range.Max)))
	fmt.Printf(" %s  %d ports scanned\n", styleGreen.Render("✓"), len(result.Available))

	if len(result.Available) == 0 {
		fmt.Println(styleMuted.Render("  no free ports found"))
	} else {
		printFreeTable(result.Available)
	}

	if len(result.Blocked) > 0 {
		fmt.Println()
		for _, b := range result.Blocked {
			fmt.Printf(" %s %d  %s\n", styleRed.Render("✗"), b.Port, styleMuted.Render(b.Reason))
		}
	}

	if len(result.Available) > 0 {
		ports := make([]string, len(result.Available))
		for i, p := range result.Available {
			ports[i] = fmt.Sprintf("%d", p)
		}
		fmt.Printf("\n %s  %s\n", styleMuted.Render("copy:"), styleBlue.Render(strings.Join(ports, " ")))
	}
}

func printFreeTable(ports []int) {
	width := 6 + 20 + 8 + 10
	border := fmt.Sprintf(" %s%s%s", "┌", strings.Repeat("─", width), "┐")
	fmt.Println(border)

	header := fmt.Sprintf(" %s %s %s %s", "│",
		styleDim.Width(6).Render("PORT"),
		styleDim.Width(20).Render("SERVICE"),
		styleDim.Width(8).Render("STATUS"),
	)
	fmt.Println(header)

	midBorder := fmt.Sprintf(" %s%s%s", "├", strings.Repeat("─", width), "┤")
	fmt.Println(midBorder)

	for i, p := range ports {
		reg := lookupLabel(p)
		portStr := styleGreen.Width(6).Render(fmt.Sprintf("%d", p))
		svcStr := stylePrimary.Width(20).Render(truncate(reg, 20))
		statusStr := styleGreen.Width(8).Render("free")
		fmt.Printf(" %s %s %s %s\n", "│", portStr, svcStr, statusStr)
		if i < len(ports)-1 {
			rowSep := fmt.Sprintf(" %s%s%s", "├", strings.Repeat("─", width), "┤")
			fmt.Println(rowSep)
		}
	}

	botBorder := fmt.Sprintf(" %s%s%s", "└", strings.Repeat("─", width), "┘")
	fmt.Println(botBorder)
}

func lookupLabel(port int) string {
	reg := registry.Lookup(port)
	if reg != nil {
		return reg.Name
	}
	return ""
}

func PrintSuggestTable(stack string, layout []types.RegistryEntry, active map[int]*types.PortEntry, conflicts []*types.ConflictPair) {
	fmt.Printf("%s portlens suggest --stack %s\n", styleGreen.Render("❯"), stack)
	fmt.Println()
	fmt.Printf("%s  layout for %s  ·  cross-checked against active ports\n", styleDim.Render("·"), styleBlue.Render(stack))
	fmt.Println()

	conflictPorts := map[int]bool{}
	for _, c := range conflicts {
		conflictPorts[c.Port] = true
	}

	categories := []string{}
	catSeen := map[string]bool{}
	for _, e := range layout {
		if !catSeen[e.Category] {
			categories = append(categories, e.Category)
			catSeen[e.Category] = true
		}
	}

	for _, cat := range categories {
		fmt.Printf("%s\n", styleDim.Render(strings.ToUpper(cat)))
		for _, e := range layout {
			if e.Category != cat {
				continue
			}
			portStyle := CategoryColor(cat)
			annotation := ""
			if conflictPorts[e.Port] {
				annotation = "  " + styleRed.Render("conflict")
			} else if _, inUse := active[e.Port]; inUse {
				annotation = "  " + styleAmber.Render("in use")
			}
			fmt.Printf("  %s  %s%s\n",
				portStyle.Render(fmt.Sprintf("%d", e.Port)),
				stylePrimary.Render(e.Name),
				annotation,
			)
		}
		fmt.Println()
	}
}

func PrintSuggestCompose(layout []types.RegistryEntry, active map[int]*types.PortEntry) {
	fmt.Println(styleMuted.Render("# portlens suggest  --export compose"))
	fmt.Println(styleMuted.Render("# paste into your docker-compose.yml service block"))
	fmt.Println()
	for _, e := range layout {
		if _, inUse := active[e.Port]; inUse {
			fmt.Printf("# %s (already in use)\n", e.Name)
			continue
		}
		fmt.Printf("      - \"%d:%d\"  # %s\n", e.Port, e.Port, e.Name)
	}
}

func PrintSuggestEnv(layout []types.RegistryEntry, active map[int]*types.PortEntry) {
	fmt.Println(styleMuted.Render("# portlens suggest  --export env"))
	for _, e := range layout {
		key := toEnvKey(e.Name)
		fmt.Printf("%s=%d\n", key, e.Port)
	}
}

func PrintSuggestJSON(layout []types.RegistryEntry, active map[int]*types.PortEntry) {
	type row struct {
		Port     int    `json:"port"`
		Name     string `json:"name"`
		Category string `json:"category"`
		InUse    bool   `json:"in_use"`
	}
	var out []row
	for _, e := range layout {
		_, inUse := active[e.Port]
		out = append(out, row{
			Port:     e.Port,
			Name:     e.Name,
			Category: e.Category,
			InUse:    inUse,
		})
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
}

func PrintFixSuggestion(pair *types.ConflictPair, alt int) {
	content := fmt.Sprintf("  conflict on port %d\n", pair.Port)
	if alt != 0 {
		content += fmt.Sprintf("  %s", styleAmber.Render("move to ")+styleGreen.Render(fmt.Sprintf("%d", alt)))
	} else {
		content += styleMuted.Render("no conventional alternative found in registry")
	}
	fmt.Println(styleConflictBlock.Render(content))
}

func PrintInfoLine(msg string) {
	content := "  " + styleBlue.Render("info") + "  " + styleMuted.Render(msg)
	fmt.Println(styleInfoBlock.Render(content))
}

func PrintErrorLine(msg string) {
	fmt.Fprintf(os.Stderr, "  %s  %s\n", styleRed.Render("error"), styleMuted.Render(msg))
}

func PrintWatchHeader(interval int) {
	fmt.Printf("%s portlens watch\n", styleGreen.Render("❯"))
	fmt.Println()
	fmt.Printf("%s  live  ·  refresh %ds  ·  press %s to quit\n",
		styleGreen.Render("●"), interval, stylePrimary.Render("Ctrl+C"))
	fmt.Println()

	fmt.Printf("%s\n", styleDim.Render("TIME       EV   PORT     SERVICE              SOURCE"))
	fmt.Println(styleDivider.Render(divider))
}

func PrintWatchEvent(ts, event string, port int, e *types.PortEntry) {
	evStyle := styleDim
	evChar := "·"
	switch event {
	case "+":
		evStyle = styleGreen
		evChar = "+"
	case "-":
		evStyle = styleAmber
		evChar = "-"
	case "!":
		evStyle = styleRed
		evChar = "!"
	}

	name := e.RegistryName
	if name == "" {
		name = e.ProcessName
	}
	if name == "" {
		name = e.ContainerName
	}
	if name == "" {
		name = "?"
	}

	tag := SourceTag(e)

	fmt.Printf("%s  %s  %s  %s  %s\n",
		styleDim.Width(10).Render(ts),
		evStyle.Width(4).Render(evChar),
		stylePrimary.Width(8).Render(fmt.Sprintf("%d", port)),
		stylePrimary.Width(22).Render(truncate(name, 22)),
		tag,
	)
}

func toEnvKey(name string) string {
	var out []byte
	for _, c := range name {
		switch {
		case c >= 'a' && c <= 'z':
			out = append(out, byte(c-32))
		case c >= 'A' && c <= 'Z':
			out = append(out, byte(c))
		default:
			out = append(out, '_')
		}
	}
	return string(out)
}
