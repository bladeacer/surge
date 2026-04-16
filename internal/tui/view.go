package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/SurgeDM/Surge/internal/config"
	"github.com/SurgeDM/Surge/internal/tui/colors"
	"github.com/SurgeDM/Surge/internal/tui/components"
	"github.com/SurgeDM/Surge/internal/utils"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Viewport layout
const maxUIDuration = 30 * 24 * time.Hour

// formatDurationForUI formats a duration as a human-readable clock string.
// Returns "M:SS" for sub-hour durations, "H:MM:SS" for multi-hour, "Xd Yh" for days.
func formatDurationForUI(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	if d >= maxUIDuration {
		return "\u221e"
	}

	totalSec := int(d.Seconds())

	if totalSec >= 86400 {
		days := totalSec / 86400
		hours := (totalSec % 86400) / 3600
		if hours > 0 {
			return fmt.Sprintf("%dd %dh", days, hours)
		}
		return fmt.Sprintf("%dd", days)
	}

	hours := totalSec / 3600
	mins := (totalSec % 3600) / 60
	secs := totalSec % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, mins, secs)
	}
	return fmt.Sprintf("%d:%02d", mins, secs)
}

// renderModalWithOverlay renders a modal centered on screen with a dark overlay effect
func (m RootModel) renderModalWithOverlay(modal string) string {
	// Place modal centered with dark gray background fill for overlay effect
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal,
		lipgloss.WithWhitespaceChars(" "), // Changed from "░" to avoid terminal rendering glitches
		lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Background(lipgloss.Color("236"))),
	)
}

func (m RootModel) wrapView(content string) tea.View {
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m RootModel) View() tea.View {
	if m.width == 0 {
		return m.wrapView("Loading...")
	}

	// Terminal too small to render any meaningful layout
	if m.width < MinTermWidth || m.height < MinTermHeight {
		msg := lipgloss.NewStyle().Foreground(colors.Cyan).Render(fmt.Sprintf("Terminal too small (min: %d×%d)", MinTermWidth, MinTermHeight))
		return m.wrapView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, msg))
	}

	if m.shuttingDown {
		modal := components.ConfirmationModal{
			Title:       "Shutting Down",
			Message:     "Pausing downloads and saving resume state...",
			Detail:      "Please wait",
			Keys:        components.NoKeys{},
			Help:        m.help,
			BorderColor: colors.Cyan,
			Width:       60,
			Height:      10,
		}
		modal.Width, modal.Height = GetDynamicModalDimensions(m.width, m.height, 40, 6, 60, 10)
		box := modal.RenderWithBtopBox(renderBtopBox, PaneTitleStyle)
		return m.wrapView(m.renderModalWithOverlay(box))
	}

	// === Handle Modal States First ===
	// These overlays sit on top of the dashboard or replace it

	if m.state == InputState {
		modal := components.AddDownloadModal{
			Title:           "Add Download",
			Inputs:          []textinput.Model{m.inputs[0], m.inputs[1], m.inputs[2], m.inputs[3]},
			Labels:          []string{"URL:", "Mirrors:", "Path:", "Filename:"},
			FocusedInput:    m.focusedInput,
			BrowseHintIndex: 2,
			Help:            m.help,
			HelpKeys:        m.keys.Input,
<<<<<<< HEAD
			BorderColor:     colors.NeonPink,
=======
			BorderColor:     colors.Pink,
			Width:           80,
			Height:          11,
>>>>>>> 29b8130 (feat: Removed Neon prefix for colours)
		}
		// Resolve dynamic dimensions
		w, _ := GetDynamicModalDimensions(m.width, m.height, 46, 8, 80, 0)
		modal.Width = w
		h := lipgloss.Height(modal.View()) + BoxStyle.GetVerticalFrameSize()
		_, modal.Height = GetDynamicModalDimensions(m.width, m.height, 46, 8, w, h)

		box := modal.RenderWithBtopBox(renderBtopBox, PaneTitleStyle)
		return m.wrapView(m.renderModalWithOverlay(box))
	}

	if m.state == FilePickerState {
		// Create a local copy to avoid modifying model during view (though View takes value receiver m)
		fp := m.filepicker
		picker := components.NewFilePickerModal(
			" Select Directory ",
			&fp,
			m.help,
			m.keys.FilePicker,
			colors.Pink,
		)
		// Resolve dynamic dimensions
		w, h := GetDynamicModalDimensions(m.width, m.height, 60, 10, 90, 20)
		picker.Width = w
		picker.Height = h

		box := picker.RenderWithBtopBox(renderBtopBox, PaneTitleStyle)
		return m.wrapView(m.renderModalWithOverlay(box))
	}

	if m.state == SettingsState {
		return m.wrapView(m.viewSettings())
	}

	if m.state == CategoryManagerState {
		return m.wrapView(m.viewCategoryManager())
	}

	if m.state == DuplicateWarningState {
		modal := components.ConfirmationModal{
			Title:       "\u26a0 Duplicate Detected",
			Message:     "A download with this URL already exists",
			Detail:      truncateString(m.duplicateInfo, 50),
			Keys:        m.keys.Duplicate,
			Help:        m.help,
<<<<<<< HEAD
			BorderColor: colors.NeonPink,
=======
			BorderColor: colors.Pink,
			Width:       60,
			Height:      10,
>>>>>>> 29b8130 (feat: Removed Neon prefix for colours)
		}
		// Resolve dynamic dimensions
		w, _ := GetDynamicModalDimensions(m.width, m.height, 40, 6, 60, 0)
		modal.Width = w
		// ConfirmationModal's internal height calculation depends on width (for help wrap)
		// but since it's a fixed-width confirmation message, we can approximate or call View()
		// Note: ConfirmationModal renders itself into the height passed,
		// so we need a reasonable estimate for 'h'.
		h := 10 // typical height for confirmation
		if m.duplicateInfo != "" {
			h = 11
		}
		_, modal.Height = GetDynamicModalDimensions(m.width, m.height, 40, 6, w, h)

		box := modal.RenderWithBtopBox(renderBtopBox, PaneTitleStyle)
		return m.wrapView(m.renderModalWithOverlay(box))
	}

	if m.state == ExtensionConfirmationState {
		extInputs := []textinput.Model{m.inputs[2], m.inputs[3]}
		focused := 0
		if m.focusedInput == 3 {
			focused = 1
		}
		modal := components.AddDownloadModal{
			Title:           "Extension Download",
			Inputs:          extInputs,
			Labels:          []string{"Path:", "Filename:"},
			FocusedInput:    focused,
			ShowURL:         true,
			URL:             truncateString(m.pendingURL, 68),
			BrowseHintIndex: 0,
			Help:            m.help,
			HelpKeys:        m.keys.Extension,
<<<<<<< HEAD
			BorderColor:     colors.NeonCyan,
=======
			BorderColor:     colors.Cyan,
			Width:           86,
			Height:          13,
>>>>>>> 29b8130 (feat: Removed Neon prefix for colours)
		}
		// Resolve dynamic dimensions
		w, _ := GetDynamicModalDimensions(m.width, m.height, 60, 10, 86, 0)
		modal.Width = w
		h := lipgloss.Height(modal.View()) + BoxStyle.GetVerticalFrameSize()
		_, modal.Height = GetDynamicModalDimensions(m.width, m.height, 60, 10, w, h)

		box := modal.RenderWithBtopBox(renderBtopBox, PaneTitleStyle)
		return m.wrapView(m.renderModalWithOverlay(box))
	}

	if m.state == BatchFilePickerState {
		fp := m.filepicker
		picker := components.NewFilePickerModal(
			" Select URL File (.txt) ",
			&fp,
			m.help,
			m.keys.FilePicker,
			colors.Cyan,
		)
		// Resolve dynamic dimensions
		w, h := GetDynamicModalDimensions(m.width, m.height, 60, 10, 90, 20)
		picker.Width = w
		picker.Height = h

		box := picker.RenderWithBtopBox(renderBtopBox, PaneTitleStyle)
		return m.wrapView(m.renderModalWithOverlay(box))
	}

	if m.state == BatchConfirmState {
		urlCount := len(m.pendingBatchURLs)
		modal := components.ConfirmationModal{
			Title:       "Batch Import",
			Message:     fmt.Sprintf("Add %d downloads?", urlCount),
			Detail:      truncateString(m.batchFilePath, 50),
			Keys:        m.keys.BatchConfirm,
			Help:        m.help,
<<<<<<< HEAD
			BorderColor: colors.NeonCyan,
=======
			BorderColor: colors.Cyan,
			Width:       60,
			Height:      10,
>>>>>>> 29b8130 (feat: Removed Neon prefix for colours)
		}
		// Resolve dynamic dimensions
		w, _ := GetDynamicModalDimensions(m.width, m.height, 40, 6, 60, 0)
		modal.Width = w
		h := 10 // typical height for confirmation
		_, modal.Height = GetDynamicModalDimensions(m.width, m.height, 40, 6, w, h)

		box := modal.RenderWithBtopBox(renderBtopBox, PaneTitleStyle)
		return m.wrapView(m.renderModalWithOverlay(box))
	}

	if m.state == QuitConfirmState {
		return m.wrapView(m.renderModalWithOverlay(m.viewQuitConfirm()))
	}

	if m.state == UpdateAvailableState && m.UpdateInfo != nil {
		modal := components.ConfirmationModal{
			Title:       "\u2b06 Update Available",
			Message:     fmt.Sprintf("A new version of Surge is available: %s", m.UpdateInfo.LatestVersion),
			Detail:      fmt.Sprintf("Current: %s", m.UpdateInfo.CurrentVersion),
			Keys:        m.keys.Update,
			Help:        m.help,
<<<<<<< HEAD
			BorderColor: colors.NeonCyan,
=======
			BorderColor: colors.Cyan,
			Width:       60,
			Height:      12,
>>>>>>> 29b8130 (feat: Removed Neon prefix for colours)
		}
		// Resolve dynamic dimensions
		w, _ := GetDynamicModalDimensions(m.width, m.height, 50, 8, 60, 0)
		modal.Width = w
		h := 12 // typical height for update prompt
		_, modal.Height = GetDynamicModalDimensions(m.width, m.height, 50, 8, w, h)

		box := modal.RenderWithBtopBox(renderBtopBox, PaneTitleStyle)
		return m.wrapView(m.renderModalWithOverlay(box))
	}

	if m.state == URLUpdateState {
		modal := components.AddDownloadModal{
			Title:           "Refresh URL",
			Inputs:          []textinput.Model{m.urlUpdateInput},
			Labels:          []string{"New URL:"},
			FocusedInput:    0,
			BrowseHintIndex: -1, // No browse hint needed
			Help:            m.help,
			HelpKeys:        m.keys.Input,
<<<<<<< HEAD
			BorderColor:     colors.NeonPink,
=======
			BorderColor:     colors.Pink,
			Width:           80,
			Height:          8,
>>>>>>> 29b8130 (feat: Removed Neon prefix for colours)
		}
		// Resolve dynamic dimensions
		w, _ := GetDynamicModalDimensions(m.width, m.height, 46, 6, 80, 0)
		modal.Width = w
		h := lipgloss.Height(modal.View()) + BoxStyle.GetVerticalFrameSize()
		_, modal.Height = GetDynamicModalDimensions(m.width, m.height, 46, 6, w, h)

		box := modal.RenderWithBtopBox(renderBtopBox, PaneTitleStyle)
		return m.wrapView(m.renderModalWithOverlay(box))
	}

	if m.state == HelpModalState {
		w, h := GetDynamicModalDimensions(m.width, m.height, 40, 10, PopupWidth, 22)
		modal := components.HelpModal{
			Title:       "Keyboard Shortcuts",
			HelpKeys:    m.keys.Dashboard,
			Help:        m.help,
<<<<<<< HEAD
			BorderColor: colors.NeonCyan,
			Width:       w,
			Height:      h,
=======
			BorderColor: colors.Cyan,
			Width:       modalW,
			Height:      modalH,
>>>>>>> 29b8130 (feat: Removed Neon prefix for colours)
		}
		box := modal.RenderWithBtopBox(renderBtopBox, PaneTitleStyle)
		return m.wrapView(m.renderModalWithOverlay(box))
	}

	// === MAIN DASHBOARD LAYOUT ===
	layout := CalculateDashboardLayout(m.width, m.height)

	// Footer - keybindings on left, version on bottom-right
	helpText := m.help.View(m.keys.Dashboard)
	versionBlue := colors.ThemeColor("#005cc5", "#58a6ff")
	versionText := lipgloss.NewStyle().Foreground(versionBlue).Render(fmt.Sprintf("v%s", m.CurrentVersion))

	// Hide help text at very narrow widths — version is more important
	var footerContent string
	if layout.AvailableWidth < 60 {
		footerContent = versionText
	} else {
		leftFooterWidth := layout.AvailableWidth - lipgloss.Width(versionText)
		if leftFooterWidth < 0 {
			leftFooterWidth = 0
		}
		footerContent = lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Width(leftFooterWidth).Render(helpText),
			versionText,
		)
	}
	footer := footerContent

	// Pre-calculate data needed for sub-renders
	stats := m.ComputeViewStats()
	selected := m.GetSelectedDownload()

<<<<<<< HEAD
=======
	detailWidth := rightWidth - PaneStyle.GetHorizontalFrameSize()
	if detailWidth < 0 {
		detailWidth = 0
	}

	if selected != nil {
		detailContent = renderFocusedDetails(selected, detailWidth, m.spinner.View())
	} else {
		// Default Placeholder
		detailContent = lipgloss.Place(detailWidth, 8, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(colors.Cyan).Render("No Download Selected"))
	}

	// Exact height from content + borders
	detailHeight := lipgloss.Height(detailContent) + BoxStyle.GetVerticalFrameSize()

	// Calculate Available Height for Rest
	remainingHeight := availableHeight - detailHeight
	if remainingHeight < 0 {
		remainingHeight = 0
	}

	// Calculate Chunk Map Needs
	chunkMapHeight := 0
	chunkMapNeeded := 0
	showChunkMap := false

	// Pre-fetch bitmap data if available
>>>>>>> 29b8130 (feat: Removed Neon prefix for colours)
	var bitmap []byte
	var bitmapWidth int
	var totalSize, chunkSize int64
	var chunkProgress []int64
	if selected != nil && selected.state != nil {
		bitmap, bitmapWidth, totalSize, chunkSize, chunkProgress = selected.state.GetBitmap()
	}

	// Pre-compute details content to avoid double-computation and width mismatches
	var detailContent string
	detailWidth := layout.RightWidth
	if layout.HideRightColumn {
		detailWidth = layout.LeftWidth
	}
	if selected != nil {
		detailContent = renderFocusedDetails(selected, detailWidth-components.BorderFrameWidth, m.spinner.View())
	} else {
		detailContent = renderEmptyMessage(detailWidth-components.BorderFrameWidth, layout.DetailHeight-components.BorderFrameHeight, "No download selected")
	}

<<<<<<< HEAD
	// Render Components
	logoColumn := m.renderHeaderBox(layout.LogoWidth, layout.HeaderHeight)
	logBox := m.renderLogBox(layout.LogWidth, layout.HeaderHeight)
=======
	// Recalculate Graph Area for rendering usage later
	// graphHeight is now set vertically.

	// --- SECTION 1: HEADER & LOGO (Top Left) + LOG BOX (Top Right) ---
	logoText := `
   _______  ___________ ____ 
  / ___/ / / / ___/ __ '/ _ \
 (__  ) /_/ / /  / /_/ /  __/
/____/\__,_/_/   \__, /\___/ 
                /____/       `

	// Calculate stats for tab bar
	stats := m.ComputeViewStats()
	active := stats.ActiveCount
	queued := stats.QueuedCount
	downloaded := stats.DownloadedCount

	// Logo takes ~45% of header width
	logoWidth := int(float64(leftWidth) * LogoWidthRatio)
	logWidth := leftWidth - logoWidth - BoxStyle.GetHorizontalFrameSize() // Rest for log box

	if logoWidth < 4 {
		logoWidth = 4 // Minimum for server box content
	}
	if logWidth < 4 {
		logWidth = 4 // Minimum for viewport
	}

	// Server info vars
	greenDot := lipgloss.NewStyle().Foreground(colors.StateDownloading).Render("●")
	host := m.ServerHost
	if host == "" {
		host = "127.0.0.1"
	}
	serverAddr := fmt.Sprintf("%s:%d", host, m.ServerPort)

	var statusLine string
	if m.IsRemote {
		statusLine = lipgloss.NewStyle().Foreground(colors.Cyan).Bold(true).Render(" Connected to " + serverAddr)
	} else {
		statusLine = lipgloss.NewStyle().Foreground(colors.Cyan).Bold(true).Render(" Serving at " + serverAddr)
	}

	serverContentWidth := logoWidth - (BoxStyle.GetHorizontalFrameSize() * 2)
	if serverContentWidth < 0 {
		serverContentWidth = 0
	}
	serverPortContent := lipgloss.NewStyle().
		Width(serverContentWidth).
		Align(lipgloss.Center).
		Render(greenDot + statusLine)
	serverBoxHeight := lipgloss.Height(serverPortContent) + 2
	if serverBoxHeight < 3 {
		serverBoxHeight = 3
	}

	// Render logo column (or just server info when too narrow)
	var logoColumn string
	if hideLogo {
		logoColumn = renderBtopBox("", PaneTitleStyle.Render(" Server "), serverPortContent, logoWidth, serverBoxHeight, colors.Gray)
	} else {
		var logoContent string
		if m.logoCache != "" {
			logoContent = m.logoCache
		} else {
			gradientLogo := ApplyGradient(logoText, colors.Pink, colors.Magenta)
			m.logoCache = lipgloss.NewStyle().Render(gradientLogo)
			logoContent = m.logoCache
		}

		logoBoxHeight := headerHeight - serverBoxHeight
		if logoBoxHeight < 1 {
			logoBoxHeight = 1
		}
		logoBox := lipgloss.Place(logoWidth, logoBoxHeight, lipgloss.Center, lipgloss.Center, logoContent)
		serverBox := renderBtopBox("", PaneTitleStyle.Render(" Server "), serverPortContent, logoWidth, serverBoxHeight, colors.Gray)
		logoColumn = lipgloss.JoinVertical(lipgloss.Left, logoBox, serverBox)
	}

	// Render log viewport
	vpWidth := logWidth - (BoxStyle.GetHorizontalFrameSize() * 2)
	if vpWidth < 0 {
		vpWidth = 0
	}
	vpHeight := headerHeight - (BoxStyle.GetVerticalFrameSize() * 2)
	if vpHeight < 1 {
		vpHeight = 1
	}
	m.logViewport.SetWidth(vpWidth)
	m.logViewport.SetHeight(vpHeight)
	logContent := m.logViewport.View()

	// Use different border color when focused
	logBorderColor := colors.Gray
	if m.logFocused {
		logBorderColor = colors.Pink
	}
	logBox := renderBtopBox(PaneTitleStyle.Render(" Activity Log "), "", logContent, logWidth, headerHeight, logBorderColor)

	// Combine logo column and log box horizontally
>>>>>>> 29b8130 (feat: Removed Neon prefix for colours)
	headerBox := lipgloss.JoinHorizontal(lipgloss.Top, logoColumn, logBox)

	listBox := m.renderDownloadsBox(layout.LeftWidth, layout.ListHeight, stats)

	// Right column
	var rightColumn string
	if !layout.HideRightColumn {
		// Show chunk map only if we have actual data to visualize
		hasChunks := len(bitmap) > 0 && bitmapWidth > 0
		showActualChunkMap := layout.ShowChunkMap && hasChunks && selected != nil && !selected.done

		// If we reserved space for chunk map but aren't showing it, give it to details
		if !showActualChunkMap && layout.ShowChunkMap {
			layout.DetailHeight += layout.ChunkMapHeight
		}

		graphBox := m.renderGraphBox(layout.RightWidth, layout.GraphHeight, stats)
		detailBox := m.renderDetailsBox(layout.RightWidth, layout.DetailHeight, detailContent)

<<<<<<< HEAD
=======
	// Calculate Available Height for the Graph
	graphContentHeight := graphHeight - BoxStyle.GetVerticalFrameSize() - LayoutGapStyle.GetVerticalFrameSize() - 2 // remaining padding
	if graphContentHeight < 3 {
		graphContentHeight = 3
	}

	// Stats box width inside the Network Activity box
	statsBoxWidth := GraphStatsWidth

	// Graph width calculation: hide stats box when too narrow
	buildAxisLines := func(height int, axisStyle lipgloss.Style) []string {
		label := func(v float64) string {
			if v <= 0 {
				return "0 MB/s"
			}
			return fmt.Sprintf("%.1f MB/s", v)
		}

		axisLines := make([]string, height)
		for i := range axisLines {
			axisLines[i] = axisStyle.Render("")
		}

		type axisMark struct {
			num int
			den int
		}

		marks := []axisMark{
			{num: 1, den: 1},
			{num: 1, den: 2},
			{num: 0, den: 1},
		}
		if height >= 9 {
			marks = []axisMark{
				{num: 1, den: 1},
				{num: 4, den: 5},
				{num: 3, den: 5},
				{num: 2, den: 5},
				{num: 1, den: 5},
				{num: 0, den: 1},
			}
		}

		for _, mark := range marks {
			row := 0
			if height > 1 {
				row = ((mark.den-mark.num)*(height-1) + mark.den/2) / mark.den
			}
			value := maxSpeed * float64(mark.num) / float64(mark.den)
			axisLines[row] = axisStyle.Render(label(value))
		}

		return axisLines
	}
	var graphWithAxis string
	if hideGraphStats {
		// No stats box — graph gets almost full width
		graphAreaWidth, axisWidth := GetGraphAreaDimensions(rightWidth, true)

		graphVisual := renderMultiLineGraph(graphData, graphAreaWidth, graphContentHeight, maxSpeed, nil)

		// Y-axis labels
		axisStyle := lipgloss.NewStyle().Width(axisWidth).Foreground(colors.Cyan).Align(lipgloss.Right)
		axisLines := buildAxisLines(graphContentHeight, axisStyle)
		axisColumn := lipgloss.NewStyle().
			Height(graphContentHeight).
			Align(lipgloss.Right).
			Render(strings.Join(axisLines, "\n"))

		graphWithAxis = lipgloss.JoinHorizontal(lipgloss.Top,
			graphVisual,
			axisColumn,
		)
	} else {
		// Get current speed and calculate total downloaded
		currentSpeed := 0.0
		if len(m.SpeedHistory) > 0 {
			currentSpeed = m.SpeedHistory[len(m.SpeedHistory)-1]
		}

		// Calculate total downloaded across all downloads
		totalDownloaded := stats.TotalDownloaded

		// Create stats content (left side inside box)
		speedMbps := currentSpeed * 8
		topMbps := topSpeed * 8

		valueStyle := lipgloss.NewStyle().Foreground(colors.Cyan).Bold(true)
		labelStyleStats := lipgloss.NewStyle().Foreground(colors.LightGray)
		dimStyle := lipgloss.NewStyle().Foreground(colors.Gray)

		statsContent := lipgloss.JoinVertical(lipgloss.Left,
			fmt.Sprintf("%s %s", valueStyle.Render("▼"), valueStyle.Render(fmt.Sprintf("%.2f MB/s", currentSpeed))),
			dimStyle.Render(fmt.Sprintf("  (%.0f Mbps)", speedMbps)),
			"",
			fmt.Sprintf("%s %s", labelStyleStats.Render("Top:"), valueStyle.Render(fmt.Sprintf("%.2f", topSpeed))),
			dimStyle.Render(fmt.Sprintf("  (%.0f Mbps)", topMbps)),
			"",
			fmt.Sprintf("%s %s", labelStyleStats.Render("Total:"), valueStyle.Render(utils.ConvertBytesToHumanReadable(totalDownloaded))),
		)

		// Style stats with a border box
		statsBoxStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.Gray).
			Padding(0, 1).
			Width(statsBoxWidth).
			Height(graphContentHeight)
		statsBox := statsBoxStyle.Render(statsContent)

		// Graph takes remaining width after stats box
		graphAreaWidth, axisWidth := GetGraphAreaDimensions(rightWidth, false)

		graphVisual := renderMultiLineGraph(graphData, graphAreaWidth, graphContentHeight, maxSpeed, nil)

		// Create Y-axis (right side of graph)
		axisStyle := lipgloss.NewStyle().Width(axisWidth).Foreground(colors.Cyan).Align(lipgloss.Right)
		axisLines := buildAxisLines(graphContentHeight, axisStyle)

		axisColumn := lipgloss.NewStyle().
			Height(graphContentHeight).
			Align(lipgloss.Right).
			Render(strings.Join(axisLines, "\n"))

		graphWithAxis = lipgloss.JoinHorizontal(lipgloss.Top,
			statsBox,
			graphVisual,
			axisColumn,
		)
	}

	// Add top and bottom padding inside the Network Activity box
	graphWithPadding := lipgloss.JoinVertical(lipgloss.Left,
		"", // Top padding
		graphWithAxis,
		"", // Bottom padding
	)

	// Render single network activity box containing stats + graph
	graphBox := renderBtopBox(PaneTitleStyle.Render(" Network Activity "), "", graphWithPadding, rightWidth, graphHeight, colors.Cyan)

	// Don't include graph box when too small to render
	renderGraphBox := graphHeight >= minGraphHeight

	// --- SECTION 3: DOWNLOAD LIST (Bottom Left) ---
	// Tab Bar
	tabBar := renderTabs(m.activeTab, active, queued, downloaded)

	// Search bar (shown when search is active or has a query)
	var leftTitle string
	if m.searchActive || m.searchQuery != "" {
		searchIcon := lipgloss.NewStyle().Foreground(colors.Cyan).Render("> ")
		var searchDisplay string
		if m.searchActive {
			searchDisplay = m.searchInput.View() +
				lipgloss.NewStyle().Foreground(colors.Gray).Render(" [esc exit]")
		} else {
			// Show query with clear hint
			searchDisplay = lipgloss.NewStyle().Foreground(colors.Pink).Render(m.searchQuery) +
				lipgloss.NewStyle().Foreground(colors.Gray).Render(" [f to clear]")
		}
		// Pad the search bar to look like a title block
		leftTitle = " " + lipgloss.JoinHorizontal(lipgloss.Left, searchIcon, searchDisplay) + " "
	}

	// Render the bubbles list or centered empty message
	var listContent string
	if len(m.list.Items()) == 0 {
		listContentHeight := listHeight - BoxStyle.GetVerticalFrameSize() - ModalPaddingStyle.GetVerticalFrameSize()

		listContentWidth := leftWidth - (BoxStyle.GetHorizontalFrameSize() * 4)
		if listContentWidth < 0 {
			listContentWidth = 0
		}

		if m.searchQuery != "" {
			listContent = lipgloss.Place(listContentWidth, listContentHeight, lipgloss.Center, lipgloss.Center,
				lipgloss.NewStyle().Foreground(colors.Cyan).Render("No matching downloads"))
		} else {
			listContent = lipgloss.Place(listContentWidth, listContentHeight, lipgloss.Center, lipgloss.Center,
				lipgloss.NewStyle().Foreground(colors.Cyan).Render("No downloads"))
		}
	} else {
		// ensure list fills the height
		m.list.SetHeight(listHeight - BoxStyle.GetVerticalFrameSize() - ModalPaddingStyle.GetVerticalFrameSize()) // adjust for padding/tabs
		listContent = m.list.View()
	}

	// Build list inner content - No search bar inside
	listInnerContent := lipgloss.JoinVertical(lipgloss.Left, tabBar, listContent)
	listInner := lipgloss.NewStyle().Padding(1, 2).Render(listInnerContent)

	// Determine border color for downloads box based on focus
	downloadsBorderColor := colors.Pink
	if m.logFocused {
		downloadsBorderColor = colors.Gray
	}
	listBox := renderBtopBox(leftTitle, PaneTitleStyle.Render(" Downloads "), listInner, leftWidth, listHeight, downloadsBorderColor)

	// --- SECTION 4: DETAILS PANE (Middle Right) ---
	// detailContent and selected are already calculated in the layout section

	detailBox := renderBtopBox("", PaneTitleStyle.Render(" File Details "), detailContent, rightWidth, detailHeight, colors.Gray)

	// --- SECTION 5: CHUNK MAP PANE (Bottom Right) ---
	var chunkBox string
	if showChunkMap {
		var chunkContent string
		// Bitmap data already fetched above
		if len(bitmap) > 0 {
			// New chunk map component
			// Calculate target rows based on available height (minus borders)
			targetRows := chunkMapHeight - 2
			if targetRows < 3 {
				targetRows = 3 // Minimum 3 rows
			}
			if targetRows > 5 {
				targetRows = 5 // Maximum 5 rows for compact look
			}
			chunkMapPadding := lipgloss.NewStyle().Padding(0, 2)
			chunkMapWidth := rightWidth - BoxStyle.GetHorizontalFrameSize() - chunkMapPadding.GetHorizontalFrameSize()
			if chunkMapWidth < 4 {
				chunkMapWidth = 4
			}
			chunkMap := components.NewChunkMapModel(bitmap, bitmapWidth, chunkMapWidth, targetRows, selected.paused, totalSize, chunkSize, chunkProgress)
			chunkContent = chunkMapPadding.Render(chunkMap.View()) // No bottom padding

			// If no chunks (not initialized or small file), show message
			if bitmapWidth == 0 {
				msg := "Chunk visualization not available"

				placeholderWidth := rightWidth - BoxStyle.GetHorizontalFrameSize()
				if placeholderWidth < 0 {
					placeholderWidth = 0
				}

				chunkContent = lipgloss.Place(placeholderWidth, chunkMapHeight-2, lipgloss.Center, lipgloss.Center,
					lipgloss.NewStyle().Foreground(colors.Gray).Render(msg))
			}
		}

		chunkBox = renderBtopBox("", PaneTitleStyle.Render(" Chunk Map "), chunkContent, rightWidth, chunkMapHeight, colors.Gray)
	}

	// --- ASSEMBLY ---

	var body string
	if hideRightColumn {
		// Terminal too narrow for two-column layout — list-only mode
		body = lipgloss.JoinVertical(lipgloss.Left, headerBox, listBox)
	} else {
		// Left Column
		leftColumn := lipgloss.JoinVertical(lipgloss.Left, headerBox, listBox)

		// Right Column (Graph + Detail + Chunk)
>>>>>>> 29b8130 (feat: Removed Neon prefix for colours)
		var rightParts []string
		if layout.GraphHeight >= layout.MinGraphHeight {
			rightParts = append(rightParts, graphBox)
		}
		rightParts = append(rightParts, detailBox)

		if showActualChunkMap {
			chunkBox := m.renderChunkMapBox(layout.RightWidth, layout.ChunkMapHeight, selected, bitmap, bitmapWidth, totalSize, chunkSize, chunkProgress)
			rightParts = append(rightParts, chunkBox)
		}
		rightColumn = lipgloss.JoinVertical(lipgloss.Left, rightParts...)
	}

	// Assembly
	var body string
	if layout.HideRightColumn {
		if layout.VerticalLayout {
			detailBox := m.renderDetailsBox(layout.LeftWidth, layout.DetailHeight, detailContent)
			body = lipgloss.JoinVertical(lipgloss.Left, headerBox, listBox, detailBox)
		} else {

			body = lipgloss.JoinVertical(lipgloss.Left, headerBox, listBox)
		}
	} else {

		leftColumn := lipgloss.JoinVertical(lipgloss.Left, headerBox, listBox)
		body = lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, rightColumn)
	}

	body = lipgloss.NewStyle().
		Width(layout.AvailableWidth).
		Height(layout.AvailableHeight).
		MaxWidth(layout.AvailableWidth).
		MaxHeight(layout.AvailableHeight).
		Render(body)

	fullView := lipgloss.JoinVertical(lipgloss.Left, body, footer)
	// Place content into available space, then wrap with WindowStyle margins
	return m.wrapView(lipgloss.Place(layout.AvailableWidth, m.height, lipgloss.Center, lipgloss.Top, fullView))
}

// Helper to render the detailed info pane
func renderFocusedDetails(d *DownloadModel, w int, spinnerView string) string {
	pct := 0.0
	if d.Total > 0 {
		pct = float64(d.Downloaded) / float64(d.Total)
	}

	// Consistent content width for centering
	contentWidth := w - (components.BorderFrameWidth * 2)
	if contentWidth < 0 {
		contentWidth = 0
	}

	// Separator Style
	divider := lipgloss.NewStyle().
		Foreground(colors.Gray).
		Width(contentWidth).
		Render("\n" + strings.Repeat("\u2500", contentWidth) + "\n")

	// Padding Style for sections
	sectionStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Padding(0, 1)

	// --- 1. Status Section ---
	statusStr := getDownloadStatus(d, spinnerView)
	statusStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colors.Gray).
		Width(contentWidth).
		Align(lipgloss.Center)

	statusBox := statusStyle.Render(statusStr)

	// --- 2. File Information Section ---
	displayFilename := d.Filename
	if displayFilename == "" || displayFilename == "Queued" {
		displayFilename = d.URL
	}

	displayPath := d.Destination
	if displayPath == "" {
		displayPath = d.URL
	}

	// Calculate inner width accounting for sectionStyle padding (0, 1)
	innerWidth := contentWidth - components.BorderFrameWidth
	if innerWidth < 0 {
		innerWidth = 0
	}
	valueWidth := innerWidth - 12
	if valueWidth < 5 {
		valueWidth = 5 // Minimum reasonable width
	}

	fileInfoContent := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Left, StatsLabelStyle.Render("URL: "), StatsValueStyle.Render(truncateMiddle(d.URL, valueWidth))),
		lipgloss.JoinHorizontal(lipgloss.Left, StatsLabelStyle.Render("File: "), StatsValueStyle.Render(truncateString(displayFilename, valueWidth))),
		lipgloss.JoinHorizontal(lipgloss.Left, StatsLabelStyle.Render("Path: "), StatsValueStyle.Render(truncateMiddle(displayPath, valueWidth))),
		lipgloss.JoinHorizontal(lipgloss.Left, StatsLabelStyle.Render("ID:   "), lipgloss.NewStyle().Foreground(colors.LightGray).Render(truncateString(d.ID, valueWidth))),
	)
	fileSection := sectionStyle.Render(fileInfoContent)

	// --- 3. Progress Section ---
	labelStr := "Progress: "
	progLabelStyle := lipgloss.NewStyle().Foreground(colors.NeonCyan)

	var progContent string
	if contentWidth > 45 { // Enough space for "Progress: " (10) + some bar + padding
		// Horizontal layout: Progress: [████████      ]
		maxProgWidth := contentWidth - lipgloss.Width(labelStr) - components.SingleLineHeight
		if maxProgWidth < 10 {
			maxProgWidth = 10
		}
		d.progress.SetWidth(maxProgWidth)
		progView := d.progress.ViewAs(pct)
		progContent = lipgloss.JoinHorizontal(lipgloss.Center, progLabelStyle.Render(labelStr), progView)
	} else {
		// Vertical layout for narrow terminals:
		// Progress:
		// [███████]
		maxProgWidth := contentWidth
		if maxProgWidth < 10 {
			maxProgWidth = 10 // Still clamp to a readable minimum, but we'll allow wrapping if term is REALLY tiny
		}
		// If contentWidth is actually smaller than 10, we must NOT exceed it to avoid "broken" look
		if maxProgWidth > contentWidth && contentWidth > 5 {
			maxProgWidth = contentWidth
		}

		d.progress.SetWidth(maxProgWidth)
		progView := d.progress.ViewAs(pct)

		centeredLabel := lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center).Render(progLabelStyle.Render(labelStr))
		centeredBar := lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center).Render(progView)
		progContent = lipgloss.JoinVertical(lipgloss.Left, centeredLabel, centeredBar)
	}

<<<<<<< HEAD
	progSection := lipgloss.NewStyle().Width(contentWidth).Render(progContent)
=======
	progLabel := lipgloss.NewStyle().Foreground(colors.Cyan).Render("Progress: ")
	progContent := lipgloss.JoinVertical(lipgloss.Left, progLabel, progView)

	// Progress bar has its own width handling usually, but let's wrap it to be sure
	progSection := lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center).Render(progContent)
>>>>>>> 29b8130 (feat: Removed Neon prefix for colours)

	// --- 4. Stats Grid Section ---
	var speedStr, etaStr, sizeStr, timeStr string
	// TUI owns elapsed time: compute from StartTime for active downloads,
	// use frozen d.Elapsed for completed downloads.
	var elapsed time.Duration
	if d.done {
		elapsed = d.Elapsed
	} else if d.Elapsed > 0 {
		elapsed = d.Elapsed
	} else if !d.StartTime.IsZero() {
		elapsed = time.Since(d.StartTime)
	}

	// Size
	if d.done {
		sizeStr = utils.ConvertBytesToHumanReadable(d.Total)
	} else {
		sizeStr = fmt.Sprintf("%s / %s", utils.ConvertBytesToHumanReadable(d.Downloaded), utils.ConvertBytesToHumanReadable(d.Total))
	}

	// Speed & ETA
	if d.done {
		if elapsed.Seconds() >= 1 {
			avgSpeed := float64(d.Total) / float64(int(elapsed.Seconds()))
			speedStr = fmt.Sprintf("%.2f MB/s (Avg)", avgSpeed/float64(config.MB))
		} else if d.Speed > 0 {
			speedStr = fmt.Sprintf("%.2f MB/s (Avg)", d.Speed/float64(config.MB))
		} else if elapsed.Seconds() > 0 {
			avgSpeed := float64(d.Total) / elapsed.Seconds()
			speedStr = fmt.Sprintf("%.2f MB/s (Avg)", avgSpeed/float64(config.MB))
		} else {
			speedStr = "N/A"
		}
		etaStr = "Done"
	} else if d.resuming {
		speedStr = "Resuming..."
		etaStr = "..."
	} else if d.paused || d.Speed == 0 {
		speedStr = "Paused"
		etaStr = "\u221e"
	} else {
		speedStr = fmt.Sprintf("%.2f MB/s", d.Speed/float64(config.MB))
		if d.Total > 0 {
			remaining := d.Total - d.Downloaded
			etaSeconds := float64(remaining) / d.Speed
			// Clamp ETA to 24 hours max to prevent bonkers values
			const maxETASeconds = 24 * 60 * 60
			if etaSeconds > maxETASeconds || etaSeconds < 0 {
				etaStr = "\u221e"
			} else {
				etaDuration := time.Duration(etaSeconds) * time.Second
				// EMA smooth ETA to prevent jitter from speed fluctuations
				if d.lastETA > 0 {
					const etaAlpha = 0.3
					etaDuration = time.Duration(etaAlpha*float64(etaDuration) + (1-etaAlpha)*float64(d.lastETA))
				}
				d.lastETA = etaDuration
				etaStr = formatDurationForUI(etaDuration)
			}
		} else {
			etaStr = "\u221e"
		}
	}

	timeStr = formatDurationForUI(elapsed)

	// Stats Layout
	colWidth := (contentWidth - (components.BorderFrameWidth * 2)) / 2
	leftColItems := []string{
		lipgloss.JoinHorizontal(lipgloss.Left, StatsLabelStyle.Width(7).Render("Size:"), StatsValueStyle.Render(sizeStr)),
		lipgloss.JoinHorizontal(lipgloss.Left, StatsLabelStyle.Width(7).Render("Speed:"), StatsValueStyle.Render(speedStr)),
	}
	isActive := !d.done && !d.paused && !d.pausing && d.Speed > 0
	if isActive {
		conns := d.Connections
		if conns == 0 {
			conns = 1 // Single-connection download (range requests not supported)
		}
		connStr := fmt.Sprintf("%d", conns)
		leftColItems = append(leftColItems, lipgloss.JoinHorizontal(lipgloss.Left, StatsLabelStyle.Width(7).Render("Conns:"), StatsValueStyle.Render(connStr)))
	}
	leftCol := lipgloss.JoinVertical(lipgloss.Left, leftColItems...)
	rightCol := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Left, StatsLabelStyle.Width(7).Render("Time:"), StatsValueStyle.Render(timeStr)),
		lipgloss.JoinHorizontal(lipgloss.Left, StatsLabelStyle.Width(7).Render("ETA:"), StatsValueStyle.Render(etaStr)),
	)

	statsContent := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(colWidth).Render(leftCol),
		lipgloss.NewStyle().Width(colWidth).Render(rightCol),
	)
	statsSection := sectionStyle.Render(statsContent)

	// --- 5. Mirrors Section ---
	var mirrorSection string
	if d.state != nil && len(d.state.GetMirrors()) > 0 {
		activeCount := 0
		errorCount := 0
		total := len(d.state.GetMirrors())
		for _, m := range d.state.GetMirrors() {
			if m.Active {
				activeCount++
			}
			if m.Error {
				errorCount++
			}
		}
		// More prominent Mirrors display
		mirrorLabel := StatsLabelStyle.Render("Mirrors")
		mirrorStats := lipgloss.NewStyle().Foreground(colors.LightGray).Render(fmt.Sprintf("%d Active / %d Total (%d Errors)", activeCount, total, errorCount))

		mirrorSection = sectionStyle.Render(lipgloss.JoinVertical(lipgloss.Left, mirrorLabel, mirrorStats))
	}

	// --- 6. Error Section ---
	var errorSection string
	if d.err != nil {
		errorSection = sectionStyle.
			Render(lipgloss.NewStyle().Foreground(colors.StateError).Render("Error: " + d.err.Error()))
	}

	// Combine with Dividers
	// Use explicit calls to insert divider only where needed
	var parts []string

	parts = append(parts, statusBox)
	parts = append(parts, fileSection)
	parts = append(parts, divider)
	parts = append(parts, progSection)
	parts = append(parts, divider)
	parts = append(parts, statsSection)

	if mirrorSection != "" {
		parts = append(parts, divider)
		parts = append(parts, mirrorSection)
	}

	if errorSection != "" {
		parts = append(parts, divider)
		parts = append(parts, errorSection)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, parts...)

	return lipgloss.NewStyle().
		Padding(1, 2). // Outer padding
		Render(content)
}

func getDownloadStatus(d *DownloadModel, spinnerView string) string {
	if d.pausing {
		return lipgloss.NewStyle().Foreground(colors.StatePaused).Render(spinnerView + " Pausing...")
	}
	if d.resuming {
		return lipgloss.NewStyle().Foreground(colors.StateDownloading).Render(spinnerView + " Resuming...")
	}
	status := components.DetermineStatus(d.done, d.paused, d.err != nil, d.Speed, d.Downloaded)
	return status.RenderWithSpinner(spinnerView)
}

func (m RootModel) calcTotalSpeed() float64 {
	total := 0.0
	for _, d := range m.downloads {
		// Skip completed downloads
		if d.done {
			continue
		}
		total += d.Speed
	}
	return total / float64(config.MB)
}

func (m RootModel) ComputeViewStats() ViewStats {
	var stats ViewStats
	for _, d := range m.downloads {
		if d.done {
			stats.DownloadedCount++
		} else if !d.paused && !d.pausing && (d.Speed > 0 || d.Connections > 0 || d.resuming) {
			stats.ActiveCount++
		} else {
			stats.QueuedCount++
		}
		stats.TotalDownloaded += d.Downloaded
	}
	return stats
}

func truncateString(s string, i int) string {
	if i <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= i {
		return s
	}
	if i <= 1 {
		return "\u2026"
	}
	return lipgloss.NewStyle().MaxWidth(i-1).Render(s) + "\u2026"
}

func truncateMiddle(s string, i int) string {
	if i <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= i {
		return s
	}
	if i <= 5 {
		return truncateString(s, i)
	}

	runes := []rune(s)
	// We use i-1 because \u2026 is one character
	start := (i - 1) / 2
	end := i - 1 - start

	if start+end+1 > len(runes) {
		return s
	}

	return string(runes[:start]) + "\u2026" + string(runes[len(runes)-end:])
}

func renderTabs(activeTab, activeCount, queuedCount, doneCount int) string {
	tabs := []components.Tab{
		{Label: "Queued", Count: queuedCount},
		{Label: "Active", Count: activeCount},
		{Label: "Done", Count: doneCount},
	}
	return components.RenderTabBar(tabs, activeTab, ActiveTabStyle, TabStyle)
}

func (m RootModel) viewQuitConfirm() string {
	w, h := GetDynamicModalDimensions(m.width, m.height, 40, 8, 60, 10)
	innerWidth := w - (components.BorderFrameWidth * 2)

	messageStyle := lipgloss.NewStyle().
		Foreground(colors.White).
		Width(innerWidth).
		Align(lipgloss.Center)

	detailStyle := lipgloss.NewStyle().
		Foreground(colors.Magenta).
		Bold(true).
		Width(innerWidth).
		Align(lipgloss.Center)

	pad := "   "

	activeFirst := lipgloss.NewStyle().Foreground(colors.White).Background(colors.Pink).Bold(true).Underline(true)
	activeRest := lipgloss.NewStyle().Foreground(colors.White).Background(colors.Pink).Bold(true)
	activePad := lipgloss.NewStyle().Background(colors.Pink)

	inactiveFirst := lipgloss.NewStyle().Foreground(colors.LightGray).Background(lipgloss.Color("236")).Underline(true)
	inactiveRest := lipgloss.NewStyle().Foreground(colors.LightGray).Background(lipgloss.Color("236"))
	inactivePad := lipgloss.NewStyle().Background(lipgloss.Color("236"))

	renderBtn := func(padStyle, firstStyle, restStyle lipgloss.Style, first, rest string) string {
		return padStyle.Render(pad) + firstStyle.Render(first) + restStyle.Render(rest) + padStyle.Render(pad)
	}

	yesFirst, yesRest, yesPad := activeFirst, activeRest, activePad
	noFirst, noRest, noPad := inactiveFirst, inactiveRest, inactivePad
	if m.quitConfirmFocused == 1 {
		yesFirst, yesRest, yesPad = inactiveFirst, inactiveRest, inactivePad
		noFirst, noRest, noPad = activeFirst, activeRest, activePad
	}

	yesBtn := renderBtn(yesPad, yesFirst, yesRest, "Y", "ep!")
	noBtn := renderBtn(noPad, noFirst, noRest, "N", "ope")

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, yesBtn, "     ", noBtn)
	centeredButtons := lipgloss.NewStyle().Width(innerWidth).Align(lipgloss.Center).Render(buttons)

	stats := m.ComputeViewStats()
	detail := ""
	if stats.ActiveCount > 0 {
		detail = fmt.Sprintf("%d active download(s) will be paused", stats.ActiveCount)
	}

	helpStyle := lipgloss.NewStyle().Foreground(colors.Gray).Width(innerWidth).Align(lipgloss.Center)
	helpText := helpStyle.Render(m.help.View(m.keys.QuitConfirm))

	var lines []string
	lines = append(lines, messageStyle.Render("Are you sure you want to quit?"))
	if detail != "" {
		lines = append(lines, detailStyle.Render(detail))
	}
	lines = append(lines, "")
	lines = append(lines, "")
	lines = append(lines, centeredButtons)

	innerHeight := h - components.BorderFrameHeight
	contentHeight := lipgloss.Height(lipgloss.JoinVertical(lipgloss.Left, lines...))
	helpHeight := lipgloss.Height(helpText)
	spacing := innerHeight - contentHeight - helpHeight
	if spacing < 0 {
		spacing = 0
	}
	for i := 0; i < spacing; i++ {
		lines = append(lines, "")
	}
	// Replace last line with help text if there was space, otherwise just append
	if len(lines) > 0 && spacing > 0 {
		lines[len(lines)-1] = helpText
	} else {
		lines = append(lines, helpText)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
<<<<<<< HEAD
	return renderBtopBox(PaneTitleStyle.Render(" Quit Surge "), "", content, w, h, colors.NeonPink)
=======
	return renderBtopBox(PaneTitleStyle.Render(" Quit Surge "), "", content, width, height, colors.Pink)
>>>>>>> 29b8130 (feat: Removed Neon prefix for colours)
}

// renderBtopBox creates a btop-style box with title embedded in the top border
// Supports left and right titles (e.g., search on left, pane name on right)
// Accepts pre-styled title strings
// Example: ╭─ 🔍 Search... ─────────── Downloads ─╮
// Delegates to components.RenderBtopBox for the actual rendering
var renderBtopBox = components.RenderBtopBox
