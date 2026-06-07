package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ── CLI flags ──────────────────────────────────────────────

var (
	orgKey     string
	token      string
	outputDir  string
	baseURL    string
	issueLimit int
)

func init() {
	flag.StringVar(&orgKey, "org", "", "SonarCloud organization key")
	flag.StringVar(&token, "token", "", "SonarCloud User Token")
	flag.StringVar(&outputDir, "o", "", "Output directory (e.g. D:\\workspace\\MyMind\\PPMP)")
	flag.StringVar(&baseURL, "base-url", "https://sonarcloud.io", "SonarCloud base URL")
	flag.IntVar(&issueLimit, "issues-limit", 500, "Max issues per project")
}

// ── Config file ────────────────────────────────────────────

type Config struct {
	Org        string `json:"org"`
	Output     string `json:"output"`
	Token      string `json:"token,omitempty"`
	BaseURL    string `json:"base_url,omitempty"`
	IssueLimit int    `json:"issues_limit,omitempty"`
}

func loadConfig() Config {
	cfg := Config{}
	// Search order: ./ppmp.json, $OUTPUT_DIR/ppmp.json, $HOME/.ppmp.json
	candidates := []string{"ppmp.json"}
	if outputDir != "" {
		candidates = append(candidates, filepath.Join(outputDir, "ppmp.json"))
	}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, ".ppmp.json"))
	}
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		json.Unmarshal(data, &cfg)
		break
	}
	return cfg
}

func resolveConfig() {
	cfg := loadConfig()
	// CLI flags override config
	if orgKey == "" {
		orgKey = cfg.Org
	}
	if outputDir == "" {
		outputDir = cfg.Output
	}
	if token == "" {
		token = cfg.Token
	}
	if baseURL == "https://sonarcloud.io" && cfg.BaseURL != "" {
		baseURL = cfg.BaseURL
	}
	if issueLimit == 500 && cfg.IssueLimit > 0 {
		issueLimit = cfg.IssueLimit
	}
}

// ── Types ──────────────────────────────────────────────────

type ProjectItem struct {
	Key              string `json:"key"`
	Name             string `json:"name"`
	LastAnalysisDate string `json:"lastAnalysisDate,omitempty"`
}

type ProjectsResponse struct {
	Components []ProjectItem `json:"components"`
	Paging     struct {
		Total int `json:"total"`
	} `json:"paging"`
}

type IssueItem struct {
	Key       string `json:"key"`
	Rule      string `json:"rule"`
	Severity  string `json:"severity"`
	Component string `json:"component"`
	Line      int    `json:"line"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	Effort    string `json:"effort"`
}

type IssuesResponse struct {
	Total  int         `json:"total"`
	Issues []IssueItem `json:"issues"`
}

type ComponentMeasure struct {
	Metric string `json:"metric"`
	Value  string `json:"value"`
}

type MeasuresResponse struct {
	Component struct {
		Key      string             `json:"key"`
		Name     string             `json:"name"`
		Measures []ComponentMeasure `json:"measures"`
	} `json:"component"`
}

type QualityGateResponse struct {
	ProjectStatus struct {
		Status     string `json:"status"`
		Conditions []struct {
			Status         string `json:"status"`
			MetricKey      string `json:"metricKey"`
			Comparator     string `json:"comparator"`
			ErrorThreshold string `json:"errorThreshold"`
			ActualValue    string `json:"actualValue"`
		} `json:"conditions"`
	} `json:"projectStatus"`
}

// ── Aggregation types ──────────────────────────────────────

type ProjectReport struct {
	Project   ProjectItem
	Measures  *MeasuresResponse
	QG        *QualityGateResponse
	Issues    []IssueItem
	TotalOpen int
}

type RuleStat struct {
	Rule     string
	Severity string
	Type     string
	Count    int
	Projects map[string]bool
	Examples []IssueExample
}

type IssueExample struct {
	Project string
	File    string
	Line    int
	Message string
}

type ProjectSummary struct {
	Key             string
	Name            string
	QGStatus        string
	CodeLines       string
	Coverage        string
	Bugs            string
	Vulnerabilities string
	CodeSmells      string
	Debt            string
	TotalIssues     int
	HealthScore     int
	Analyzed        bool
}

type TrendEntry struct {
	Date        string
	TotalIssues int
	Blockers    int
	Criticals   int
	Majors      int
	QGPassed    int
	QGFailed    int
	NotAnalyzed int
}

type TodoItem struct {
	Project  string
	ProjName string
	Severity string
	Rule     string
	File     string
	Line     int
	Message  string
}

// ── HTTP helper ────────────────────────────────────────────

func apiGet(path string, params url.Values) ([]byte, error) {
	u := baseURL + path
	if params != nil {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.SetBasicAuth(token, "")
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("api %s returned %d: %s", path, resp.StatusCode, string(body))
	}
	return body, nil
}

// ── Data fetchers ──────────────────────────────────────────

func fetchProjects(org string) ([]ProjectItem, error) {
	var all []ProjectItem
	page := 1
	for {
		params := url.Values{
			"organization": {org},
			"qualifiers":   {"TRK"},
			"ps":           {"100"},
			"p":            {fmt.Sprintf("%d", page)},
		}
		data, err := apiGet("/api/components/search", params)
		if err != nil {
			return nil, err
		}
		var resp ProjectsResponse
		json.Unmarshal(data, &resp)
		all = append(all, resp.Components...)
		if len(all) >= resp.Paging.Total || len(resp.Components) == 0 {
			break
		}
		page++
	}
	return all, nil
}

func fetchProjectData(project ProjectItem) (*ProjectReport, error) {
	report := &ProjectReport{Project: project}

	params := url.Values{
		"component":  {project.Key},
		"metricKeys": {"ncloc,coverage,duplicated_lines_density,bugs,vulnerabilities,code_smells,security_hotspots,sqale_index,cognitive_complexity"},
	}
	data, err := apiGet("/api/measures/component", params)
	if err != nil {
		return nil, fmt.Errorf("measures: %w", err)
	}
	json.Unmarshal(data, &report.Measures)

	params = url.Values{"projectKey": {project.Key}}
	data, err = apiGet("/api/qualitygates/project_status", params)
	if err != nil {
		return nil, fmt.Errorf("qualitygate: %w", err)
	}
	json.Unmarshal(data, &report.QG)

	var allIssues []IssueItem
	issuePage := 1
	for {
		params = url.Values{
			"componentKeys": {project.Key},
			"ps":            {"500"},
			"p":             {fmt.Sprintf("%d", issuePage)},
			"resolved":      {"false"},
			"s":             {"SEVERITY"},
			"asc":           {"false"},
		}
		data, err = apiGet("/api/issues/search", params)
		if err != nil {
			return nil, fmt.Errorf("issues: %w", err)
		}
		var issuesResp IssuesResponse
		json.Unmarshal(data, &issuesResp)
		allIssues = append(allIssues, issuesResp.Issues...)
		if len(allIssues) >= issuesResp.Total || len(issuesResp.Issues) == 0 {
			break
		}
		issuePage++
		if len(allIssues) >= issueLimit {
			break
		}
	}
	report.Issues = allIssues
	report.TotalOpen = len(allIssues)

	return report, nil
}

// ── Helpers ────────────────────────────────────────────────

func getMeasure(measures []ComponentMeasure, key string) string {
	for _, m := range measures {
		if m.Metric == key {
			return m.Value
		}
	}
	return "N/A"
}

func shortFile(component string) string {
	if idx := strings.LastIndex(component, ":"); idx >= 0 {
		return component[idx+1:]
	}
	return component
}

func classifyByPattern(issues []IssueItem) map[string][]IssueItem {
	groups := map[string][]IssueItem{}
	for _, iss := range issues {
		f := shortFile(iss.Component)
		ext := ""
		if dot := strings.LastIndex(f, "."); dot >= 0 {
			ext = strings.ToLower(f[dot:])
		}
		switch {
		case ext == ".yml" || ext == ".yaml":
			groups["YAML config"] = append(groups["YAML config"], iss)
		case ext == ".sql":
			groups["SQL script"] = append(groups["SQL script"], iss)
		case ext == ".sh" || ext == ".bash":
			groups["Shell script"] = append(groups["Shell script"], iss)
		case ext == ".vue":
			groups["Vue component"] = append(groups["Vue component"], iss)
		case ext == ".ts" || ext == ".js":
			groups["TypeScript/JS"] = append(groups["TypeScript/JS"], iss)
		case ext == ".java":
			groups["Java source"] = append(groups["Java source"], iss)
		case ext == ".go":
			groups["Go source"] = append(groups["Go source"], iss)
		case ext == ".css" || ext == ".scss":
			groups["CSS/Style"] = append(groups["CSS/Style"], iss)
		case ext == ".xml":
			groups["XML config"] = append(groups["XML config"], iss)
		default:
			groups["Other"] = append(groups["Other"], iss)
		}
	}
	return groups
}

func calcHealthScore(issues []IssueItem, analyzed bool) int {
	if !analyzed {
		return -1
	}
	score := 100
	for _, iss := range issues {
		switch iss.Severity {
		case "BLOCKER":
			score -= 10
		case "CRITICAL":
			score -= 5
		case "MAJOR":
			score -= 2
		case "MINOR":
			score -= 1
		}
	}
	if score < 0 {
		score = 0
	}
	return score
}

func scoreIcon(score int) string {
	if score < 0 {
		return "⚪"
	}
	if score >= 80 {
		return "🟢"
	}
	if score >= 60 {
		return "🟡"
	}
	if score >= 40 {
		return "🟠"
	}
	return "🔴"
}

// ── Analysis ───────────────────────────────────────────────

func analyze(reports []ProjectReport) {
	fmt.Fprintf(os.Stderr, "\nAnalyzing %d projects...\n", len(reports))

	ruleMap := map[string]*RuleStat{}
	var totalIssues int
	severityTotals := map[string]int{}
	typeTotals := map[string]int{}

	for _, r := range reports {
		totalIssues += len(r.Issues)
		for _, iss := range r.Issues {
			severityTotals[iss.Severity]++
			typeTotals[iss.Type]++
			rs, ok := ruleMap[iss.Rule]
			if !ok {
				rs = &RuleStat{
					Rule:     iss.Rule,
					Severity: iss.Severity,
					Type:     iss.Type,
					Projects: map[string]bool{},
				}
				ruleMap[iss.Rule] = rs
			}
			rs.Count++
			rs.Projects[r.Project.Key] = true
			if len(rs.Examples) < 3 {
				rs.Examples = append(rs.Examples, IssueExample{
					Project: r.Project.Key,
					File:    shortFile(iss.Component),
					Line:    iss.Line,
					Message: iss.Message,
				})
			}
		}
	}

	var rules []*RuleStat
	for _, rs := range ruleMap {
		rules = append(rules, rs)
	}
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Count > rules[j].Count
	})

	var summaries []ProjectSummary
	var qgPassed, qgFailed, notAnalyzed int
	for _, r := range reports {
		analyzed := r.Project.LastAnalysisDate != ""
		s := ProjectSummary{
			Key:         r.Project.Key,
			Name:        r.Project.Name,
			TotalIssues: len(r.Issues),
			Analyzed:    analyzed,
			HealthScore: calcHealthScore(r.Issues, analyzed),
		}
		if r.QG != nil {
			s.QGStatus = r.QG.ProjectStatus.Status
		}
		if r.Measures != nil {
			s.CodeLines = getMeasure(r.Measures.Component.Measures, "ncloc")
			s.Coverage = getMeasure(r.Measures.Component.Measures, "coverage")
			s.Bugs = getMeasure(r.Measures.Component.Measures, "bugs")
			s.Vulnerabilities = getMeasure(r.Measures.Component.Measures, "vulnerabilities")
			s.CodeSmells = getMeasure(r.Measures.Component.Measures, "code_smells")
			s.Debt = getMeasure(r.Measures.Component.Measures, "sqale_index")
		}
		if !analyzed {
			notAnalyzed++
		} else if s.QGStatus == "OK" {
			qgPassed++
		} else {
			qgFailed++
		}
		summaries = append(summaries, s)
	}

	// ── Write to latest/ ──
	latestDir := filepath.Join(outputDir, "latest")
	projDir := filepath.Join(latestDir, "projects")
	os.MkdirAll(projDir, 0755)

	// ── Diff before overwriting latest/ ──
	generateDiff(latestDir, reports)

	generateOverview(latestDir, summaries, totalIssues, severityTotals, typeTotals)
	generateRuleStats(latestDir, rules, len(reports), totalIssues)
	generateHighFreqIssues(latestDir, rules, len(reports))
	generateTodo(latestDir, reports)
	for _, r := range reports {
		generateProjectDetail(projDir, r)
	}

	// ── Snapshot archive ──
	snapshotDir := filepath.Join(outputDir, "snapshots", time.Now().Format("2006-01-02"))
	snapshotProjDir := filepath.Join(snapshotDir, "projects")
	os.MkdirAll(snapshotProjDir, 0755)
	copyDir(latestDir, snapshotDir)

	// ── Trend ──
	entry := TrendEntry{
		Date:        time.Now().Format("2006-01-02"),
		TotalIssues: totalIssues,
		Blockers:    severityTotals["BLOCKER"],
		Criticals:   severityTotals["CRITICAL"],
		Majors:      severityTotals["MAJOR"],
		QGPassed:    qgPassed,
		QGFailed:    qgFailed,
		NotAnalyzed: notAnalyzed,
	}
	appendTrend(entry)

	fmt.Fprintf(os.Stderr, "\nDone. Output: %s\n", outputDir)
	fmt.Fprintf(os.Stderr, "  latest/        <- Agent reads this\n")
	fmt.Fprintf(os.Stderr, "  snapshots/%s/ <- archived\n", time.Now().Format("2006-01-02"))
	fmt.Fprintf(os.Stderr, "  quality-trend.md <- trend data\n")
}

// ── Copy helper ────────────────────────────────────────────

func copyDir(src, dst string) {
	entries, err := os.ReadDir(src)
	if err != nil {
		return
	}
	for _, e := range entries {
		srcPath := filepath.Join(src, e.Name())
		dstPath := filepath.Join(dst, e.Name())
		if e.IsDir() {
			os.MkdirAll(dstPath, 0755)
			copyDir(srcPath, dstPath)
			continue
		}
		data, err := os.ReadFile(srcPath)
		if err != nil {
			continue
		}
		os.WriteFile(dstPath, data, 0644)
	}
}

// ── Incremental diff ───────────────────────────────────────

func generateDiff(latestDir string, currentReports []ProjectReport) {
	// Build current issue set (key -> project)
	currentIssues := map[string]string{}
	for _, r := range currentReports {
		for _, iss := range r.Issues {
			currentIssues[iss.Key] = r.Project.Key
		}
	}

	// Load previous issues from latest/projects/*.md
	// We can't perfectly reconstruct issue keys from markdown,
	// so we use a snapshot JSON approach instead.
	snapshotFile := filepath.Join(outputDir, ".snapshot-issues.json")
	var prevIssues map[string]string // issue key -> project key
	data, err := os.ReadFile(snapshotFile)
	if err == nil {
		json.Unmarshal(data, &prevIssues)
	}

	// Save current snapshot for next run
	snapshotData, _ := json.MarshalIndent(currentIssues, "", "  ")
	os.WriteFile(snapshotFile, snapshotData, 0644)

	if prevIssues == nil {
		// First run, no diff
		return
	}

	// Compute diff
	var newIssues []struct {
		Key      string
		Project  string
	}
	var resolvedIssues []struct {
		Key      string
		Project  string
	}

	for key, proj := range currentIssues {
		if _, exists := prevIssues[key]; !exists {
			newIssues = append(newIssues, struct {
				Key     string
				Project string
			}{key, proj})
		}
	}
	for key, proj := range prevIssues {
		if _, exists := currentIssues[key]; !exists {
			resolvedIssues = append(resolvedIssues, struct {
				Key     string
				Project string
			}{key, proj})
		}
	}

	// Write diff report
	var b strings.Builder
	b.WriteString("# 增量变化\n\n")
	b.WriteString(fmt.Sprintf("- **对比日期**: %s vs 上次运行\n", time.Now().Format("2006-01-02 15:04")))
	b.WriteString(fmt.Sprintf("- **新增 Issue**: %d\n", len(newIssues)))
	b.WriteString(fmt.Sprintf("- **已解决 Issue**: %d\n\n", len(resolvedIssues)))

	if len(newIssues) == 0 && len(resolvedIssues) == 0 {
		b.WriteString("_无变化_\n\n")
	} else {
		// Net change
		net := len(newIssues) - len(resolvedIssues)
		icon := "➡️"
		if net > 0 {
			icon = "📈"
		} else if net < 0 {
			icon = "📉"
		}
		b.WriteString(fmt.Sprintf("**净变化**: %s %d\n\n", icon, net))

		if len(resolvedIssues) > 0 {
			b.WriteString(fmt.Sprintf("## 🟢 已解决 (%d)\n\n", len(resolvedIssues)))
			byProj := map[string][]string{}
			for _, ri := range resolvedIssues {
				byProj[ri.Project] = append(byProj[ri.Project], ri.Key)
			}
			for proj, keys := range byProj {
				b.WriteString(fmt.Sprintf("**%s**: %d 个\n", proj, len(keys)))
			}
			b.WriteString("\n")
		}

		if len(newIssues) > 0 {
			b.WriteString(fmt.Sprintf("## 🔴 新增 (%d)\n\n", len(newIssues)))
			byProj := map[string][]string{}
			for _, ni := range newIssues {
				byProj[ni.Project] = append(byProj[ni.Project], ni.Key)
			}
			for proj, keys := range byProj {
				b.WriteString(fmt.Sprintf("**%s**: %d 个\n", proj, len(keys)))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("---\n\nRelated:\n- [[项目质量总览]]\n- [[质量趋势]]\n")
	writeFile(filepath.Join(latestDir, "增量变化.md"), b.String())
}

// ── Trend ──────────────────────────────────────────────────

func appendTrend(entry TrendEntry) {
	trendFile := filepath.Join(outputDir, "质量趋势.md")

	var entries []TrendEntry
	data, err := os.ReadFile(trendFile)
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if !strings.HasPrefix(line, "|") || strings.HasPrefix(line, "| 日期") || strings.HasPrefix(line, "|---") {
				continue
			}
			parts := strings.Split(line, "|")
			if len(parts) < 9 {
				continue
			}
			te := TrendEntry{
				Date: strings.TrimSpace(parts[1]),
			}
			fmt.Sscanf(strings.TrimSpace(parts[2]), "%d", &te.TotalIssues)
			fmt.Sscanf(strings.TrimSpace(parts[3]), "%d", &te.Blockers)
			fmt.Sscanf(strings.TrimSpace(parts[4]), "%d", &te.Criticals)
			fmt.Sscanf(strings.TrimSpace(parts[5]), "%d", &te.Majors)
			fmt.Sscanf(strings.TrimSpace(parts[6]), "%d", &te.QGPassed)
			fmt.Sscanf(strings.TrimSpace(parts[7]), "%d", &te.QGFailed)
			fmt.Sscanf(strings.TrimSpace(parts[8]), "%d", &te.NotAnalyzed)
			entries = append(entries, te)
		}
	}

	today := entry.Date
	found := false
	for i, e := range entries {
		if e.Date == today {
			entries[i] = entry
			found = true
			break
		}
	}
	if !found {
		entries = append(entries, entry)
	}

	if len(entries) > 90 {
		entries = entries[len(entries)-90:]
	}

	var b strings.Builder
	b.WriteString("# 质量趋势\n\n")
	b.WriteString("> 每次运行 `sonar-analyze` 自动追加数据。保留最近 90 条记录。\n\n")
	b.WriteString(fmt.Sprintf("- **更新时间**: %s\n\n", time.Now().Format("2006-01-02 15:04")))
	b.WriteString("| 日期 | 总 Issue | Blocker | Critical | Major | QG通过 | QG失败 | 未分析 |\n")
	b.WriteString("|------|----------|---------|----------|-------|--------|--------|--------|\n")
	for _, e := range entries {
		b.WriteString(fmt.Sprintf("| %s | %d | %d | %d | %d | %d | %d | %d |\n",
			e.Date, e.TotalIssues, e.Blockers, e.Criticals, e.Majors, e.QGPassed, e.QGFailed, e.NotAnalyzed))
	}
	b.WriteString("\n---\n\nRelated:\n- [[项目质量总览]]\n- [[质量待办]]\n")

	writeFile(trendFile, b.String())
}

// ── File generators ────────────────────────────────────────

var sevIcon = map[string]string{
	"BLOCKER": "🔴", "CRITICAL": "🟠", "MAJOR": "🟡", "MINOR": "🔵", "INFO": "⚪",
}

func generateOverview(dir string, summaries []ProjectSummary, totalIssues int, sev map[string]int, typ map[string]int) {
	var b strings.Builder
	b.WriteString("# SonarCloud 项目质量总览\n\n")
	b.WriteString(fmt.Sprintf("- **组织**: %s\n", orgKey))
	b.WriteString(fmt.Sprintf("- **项目数**: %d\n", len(summaries)))
	b.WriteString(fmt.Sprintf("- **未解决 Issue 总计**: %d\n", totalIssues))
	b.WriteString(fmt.Sprintf("- **更新时间**: %s\n\n", time.Now().Format("2006-01-02 15:04")))

	// Health score ranking
	b.WriteString("## 健康度排名\n\n")
	sorted := make([]ProjectSummary, len(summaries))
	copy(sorted, summaries)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].HealthScore < 0 {
			return false
		}
		if sorted[j].HealthScore < 0 {
			return true
		}
		return sorted[i].HealthScore > sorted[j].HealthScore
	})
	b.WriteString("| 排名 | 项目 | 健康分 | 质量门 | Issue 数 | 状态 |\n")
	b.WriteString("|------|------|--------|--------|----------|------|\n")
	for i, s := range sorted {
		status := "已分析"
		if !s.Analyzed {
			status = "未分析"
		}
		scoreStr := fmt.Sprintf("%d", s.HealthScore)
		if s.HealthScore < 0 {
			scoreStr = "-"
		}
		qg := s.QGStatus
		if qg == "OK" {
			qg = "✅"
		} else if qg == "ERROR" {
			qg = "❌"
		} else {
			qg = "-"
		}
		b.WriteString(fmt.Sprintf("| %d | [[%s]] | %s %s | %s | %d | %s |\n",
			i+1, s.Name, scoreIcon(s.HealthScore), scoreStr, qg, s.TotalIssues, status))
	}
	b.WriteString("\n")

	// Project detail table
	b.WriteString("## 各项目概况\n\n")
	b.WriteString("| 项目 | 质量门 | 代码行 | 覆盖率 | Bug | 漏洞 | Smell | 技术债务 | Issue 数 |\n")
	b.WriteString("|------|--------|--------|--------|-----|------|-------|----------|----------|\n")
	for _, s := range summaries {
		qg := s.QGStatus
		if qg == "OK" {
			qg = "✅"
		} else if qg == "ERROR" {
			qg = "❌"
		} else {
			qg = "-"
		}
		cov := s.Coverage
		if cov == "N/A" {
			cov = "-"
		} else {
			cov += "%"
		}
		b.WriteString(fmt.Sprintf("| [[%s]] | %s | %s | %s | %s | %s | %s | %s min | %d |\n",
			s.Name, qg, s.CodeLines, cov, s.Bugs, s.Vulnerabilities, s.CodeSmells, s.Debt, s.TotalIssues))
	}
	b.WriteString("\n")

	// Severity breakdown
	b.WriteString("## 严重级别分布\n\n")
	b.WriteString("| 级别 | 数量 | 占比 |\n")
	b.WriteString("|------|------|------|\n")
	sevOrder := []string{"BLOCKER", "CRITICAL", "MAJOR", "MINOR", "INFO"}
	for _, s := range sevOrder {
		c := sev[s]
		if c == 0 {
			continue
		}
		pct := float64(c) / float64(totalIssues) * 100
		b.WriteString(fmt.Sprintf("| %s %s | %d | %.1f%% |\n", sevIcon[s], s, c, pct))
	}
	b.WriteString("\n")

	// Type breakdown
	b.WriteString("## 问题类型分布\n\n")
	b.WriteString("| 类型 | 数量 | 占比 |\n")
	b.WriteString("|------|------|------|\n")
	typeOrder := []string{"BUG", "VULNERABILITY", "CODE_SMELL", "SECURITY_HOTSPOT"}
	typeLabel := map[string]string{
		"BUG": "Bug", "VULNERABILITY": "Vulnerability",
		"CODE_SMELL": "Code Smell", "SECURITY_HOTSPOT": "Security Hotspot",
	}
	for _, t := range typeOrder {
		c := typ[t]
		if c == 0 {
			continue
		}
		pct := float64(c) / float64(totalIssues) * 100
		b.WriteString(fmt.Sprintf("| %s | %d | %.1f%% |\n", typeLabel[t], c, pct))
	}
	b.WriteString("\n---\n\nRelated:\n- [[问题归类统计]]\n- [[高频问题清单]]\n- [[质量待办]]\n- [[质量趋势]]\n")

	writeFile(filepath.Join(dir, "项目质量总览.md"), b.String())
}

func generateRuleStats(dir string, rules []*RuleStat, projectCount int, totalIssues int) {
	var b strings.Builder
	b.WriteString("# 问题归类统计\n\n")
	b.WriteString(fmt.Sprintf("- **统计规则数**: %d\n", len(rules)))
	b.WriteString(fmt.Sprintf("- **Issue 总计**: %d\n", totalIssues))
	b.WriteString(fmt.Sprintf("- **更新时间**: %s\n\n", time.Now().Format("2006-01-02 15:04")))

	b.WriteString("## 按规则统计（Top 30）\n\n")
	b.WriteString("| # | 规则 | 级别 | 类型 | 数量 | 涉及项目 | 占比 |\n")
	b.WriteString("|---|------|------|------|------|----------|------|\n")
	limit := 30
	if len(rules) < limit {
		limit = len(rules)
	}
	for i := 0; i < limit; i++ {
		r := rules[i]
		pct := float64(r.Count) / float64(totalIssues) * 100
		b.WriteString(fmt.Sprintf("| %d | `%s` | %s %s | %s | %d | %d/%d | %.1f%% |\n",
			i+1, r.Rule, sevIcon[r.Severity], r.Severity, r.Type, r.Count, len(r.Projects), projectCount, pct))
	}
	b.WriteString("\n")

	b.WriteString("## 规则详情\n\n")
	for i, r := range rules {
		if i >= 50 {
			break
		}
		b.WriteString(fmt.Sprintf("### %s %s (%d hits, %d projects)\n\n", sevIcon[r.Severity], r.Rule, r.Count, len(r.Projects)))
		b.WriteString(fmt.Sprintf("- **Severity**: %s\n", r.Severity))
		b.WriteString(fmt.Sprintf("- **Type**: %s\n", r.Type))

		affected := make([]string, 0, len(r.Projects))
		for p := range r.Projects {
			affected = append(affected, p)
		}
		sort.Strings(affected)
		b.WriteString(fmt.Sprintf("- **Projects**: %s\n", strings.Join(affected, ", ")))
		b.WriteString(fmt.Sprintf("- **Rule**: [%s](%s/coding_rules?q=%s)\n\n", r.Rule, baseURL, r.Rule))

		if len(r.Examples) > 0 {
			b.WriteString("**Examples:**\n\n")
			b.WriteString("| Project | File | Line | Message |\n")
			b.WriteString("|---------|------|------|--------|\n")
			for _, ex := range r.Examples {
				line := fmt.Sprintf("%d", ex.Line)
				if ex.Line == 0 {
					line = "-"
				}
				b.WriteString(fmt.Sprintf("| `%s` | `%s` | %s | %s |\n", ex.Project, ex.File, line, ex.Message))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("---\n\nRelated:\n- [[项目质量总览]]\n- [[高频问题清单]]\n")

	writeFile(filepath.Join(dir, "问题归类统计.md"), b.String())
}

func generateHighFreqIssues(dir string, rules []*RuleStat, projectCount int) {
	var b strings.Builder
	b.WriteString("# 高频问题清单\n\n")
	b.WriteString("> Cross-project recurring issues. Read this before writing code.\n\n")
	b.WriteString(fmt.Sprintf("- **Updated**: %s\n\n", time.Now().Format("2006-01-02 15:04")))

	var highFreq []*RuleStat
	for _, r := range rules {
		if len(r.Projects) >= 2 || r.Count >= 5 {
			highFreq = append(highFreq, r)
		}
	}

	if len(highFreq) == 0 {
		b.WriteString("_No frequent issues_\n")
		writeFile(filepath.Join(dir, "高频问题清单.md"), b.String())
		return
	}

	byType := map[string][]*RuleStat{}
	for _, r := range highFreq {
		byType[r.Type] = append(byType[r.Type], r)
	}

	typeOrder := []string{"BUG", "VULNERABILITY", "CODE_SMELL", "SECURITY_HOTSPOT"}
	typeLabel := map[string]string{
		"BUG": "Bug", "VULNERABILITY": "Vulnerability",
		"CODE_SMELL": "Code Smell", "SECURITY_HOTSPOT": "Security Hotspot",
	}

	guidelines := map[string]string{
		"java:S6437":       "Hardcoded secrets in config. Use env vars or Vault.",
		"secrets:S8215":    "Bcrypt hash in SQL init. Generate at app startup.",
		"plsql:S1192":      "Duplicated string literals in SQL. Extract to constants.",
		"shelldre:S7688":   "Use [[ instead of [ for shell tests.",
		"shelldre:S7679":   "Assign $1/$2 to local vars before use.",
		"css:S7924":        "Text contrast below WCAG AA (4.5:1).",
		"Web:S6853":        "Label must associate with input via for attribute.",
		"typescript:S3358": "Nested ternary. Use if/else instead.",
	}

	for _, t := range typeOrder {
		items := byType[t]
		if len(items) == 0 {
			continue
		}
		b.WriteString(fmt.Sprintf("## %s\n\n", typeLabel[t]))
		for _, r := range items {
			b.WriteString(fmt.Sprintf("### %s `%s` (%d hits, %d projects)\n\n", sevIcon[r.Severity], r.Rule, r.Count, len(r.Projects)))
			if guide, ok := guidelines[r.Rule]; ok {
				b.WriteString(fmt.Sprintf("> %s\n\n", guide))
			}
			b.WriteString(fmt.Sprintf("- **Rule**: [%s](%s/coding_rules?q=%s)\n", r.Rule, baseURL, r.Rule))
			affected := make([]string, 0, len(r.Projects))
			for p := range r.Projects {
				affected = append(affected, "`"+p+"`")
			}
			sort.Strings(affected)
			b.WriteString(fmt.Sprintf("- **Projects**: %s\n\n", strings.Join(affected, " ")))
		}
	}

	// Agent checklist
	b.WriteString("---\n\n## Agent Coding Checklist\n\n")
	b.WriteString("Check these when generating code:\n\n")
	b.WriteString("### Security\n")
	b.WriteString("- [ ] No hardcoded secrets/passwords\n")
	b.WriteString("- [ ] No plaintext hashes in SQL scripts\n")
	b.WriteString("- [ ] No SQL injection risks\n\n")
	b.WriteString("### Java / Spring Boot\n")
	b.WriteString("- [ ] `application.yml` uses `${}` placeholders\n")
	b.WriteString("- [ ] No excessive string literal duplication\n\n")
	b.WriteString("### Shell Scripts\n")
	b.WriteString("- [ ] Use `[[ ]]` for conditionals\n")
	b.WriteString("- [ ] Positional params assigned to local vars\n\n")
	b.WriteString("### Frontend\n")
	b.WriteString("- [ ] Text contrast meets WCAG AA\n")
	b.WriteString("- [ ] `<label>` associated with input\n")
	b.WriteString("- [ ] No nested ternary expressions\n\n")

	b.WriteString("---\n\nRelated:\n- [[项目质量总览]]\n- [[问题归类统计]]\n- [[质量待办]]\n")

	writeFile(filepath.Join(dir, "高频问题清单.md"), b.String())
}

// ── Todo list (with state persistence) ─────────────────────

func generateTodo(dir string, reports []ProjectReport) {
	todoFile := filepath.Join(dir, "质量待办.md")

	// Load existing checked items to preserve state
	checkedItems := map[string]bool{}
	data, err := os.ReadFile(todoFile)
	if err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := scanner.Text()
			// Match "- [x]" pattern
			if strings.Contains(line, "- [x] ") || strings.Contains(line, "- [X] ") {
				// Extract the unique part (file + rule + message)
				trimmed := strings.TrimSpace(line)
				// Remove the checkbox prefix
				idx := strings.Index(trimmed, "] ")
				if idx >= 0 {
					key := trimmed[idx+2:]
					checkedItems[key] = true
				}
			}
		}
	}

	var b strings.Builder
	b.WriteString("# 质量待办\n\n")
	b.WriteString("> Auto-generated from BLOCKER/CRITICAL issues. Check off when fixed, re-run analyze to refresh.\n")
	b.WriteString("> Checked items are preserved across runs.\n\n")
	b.WriteString(fmt.Sprintf("- **Updated**: %s\n\n", time.Now().Format("2006-01-02 15:04")))

	var noAnalysis []string
	var todos []TodoItem

	for _, r := range reports {
		if r.Project.LastAnalysisDate == "" {
			noAnalysis = append(noAnalysis, r.Project.Key)
			continue
		}
		for _, iss := range r.Issues {
			if iss.Severity == "BLOCKER" || iss.Severity == "CRITICAL" {
				todos = append(todos, TodoItem{
					Project:  r.Project.Key,
					ProjName: r.Project.Name,
					Severity: iss.Severity,
					Rule:     iss.Rule,
					File:     shortFile(iss.Component),
					Line:     iss.Line,
					Message:  iss.Message,
				})
			}
		}
	}

	if len(noAnalysis) > 0 {
		b.WriteString("## ⚠️ 未分析项目\n\n")
		for _, p := range noAnalysis {
			b.WriteString(fmt.Sprintf("- [ ] 运行分析: `%s`\n", p))
		}
		b.WriteString("\n")
	}

	blockers := filterTodos(todos, "BLOCKER")
	if len(blockers) > 0 {
		b.WriteString(fmt.Sprintf("## 🔴 BLOCKER (%d)\n\n", len(blockers)))
		byProj := groupTodosByProject(blockers)
		for proj, items := range byProj {
			b.WriteString(fmt.Sprintf("### %s\n\n", proj))
			for _, item := range items {
				line := fmt.Sprintf("L%d", item.Line)
				if item.Line == 0 {
					line = "-"
				}
				// Build the item text for checking persistence
				itemText := fmt.Sprintf("`%s` %s %s: %s", item.File, line, item.Rule, item.Message)
				checkbox := "- [ ] "
				if checkedItems[itemText] {
					checkbox = "- [x] "
				}
				b.WriteString(fmt.Sprintf("%s%s\n", checkbox, itemText))
			}
			b.WriteString("\n")
		}
	}

	criticals := filterTodos(todos, "CRITICAL")
	if len(criticals) > 0 {
		b.WriteString(fmt.Sprintf("## 🟠 CRITICAL (%d)\n\n", len(criticals)))
		byProj := groupTodosByProject(criticals)
		for proj, items := range byProj {
			b.WriteString(fmt.Sprintf("### %s\n\n", proj))
			for _, item := range items {
				line := fmt.Sprintf("L%d", item.Line)
				if item.Line == 0 {
					line = "-"
				}
				itemText := fmt.Sprintf("`%s` %s %s: %s", item.File, line, item.Rule, item.Message)
				checkbox := "- [ ] "
				if checkedItems[itemText] {
					checkbox = "- [x] "
				}
				b.WriteString(fmt.Sprintf("%s%s\n", checkbox, itemText))
			}
			b.WriteString("\n")
		}
	}

	if len(todos) == 0 && len(noAnalysis) == 0 {
		b.WriteString("_All clear! No BLOCKER or CRITICAL issues._\n\n")
	}

	b.WriteString("---\n\nRelated:\n- [[项目质量总览]]\n- [[高频问题清单]]\n- [[质量趋势]]\n- [[增量变化]]\n")

	writeFile(todoFile, b.String())
}

func filterTodos(todos []TodoItem, severity string) []TodoItem {
	var result []TodoItem
	for _, t := range todos {
		if t.Severity == severity {
			result = append(result, t)
		}
	}
	return result
}

func groupTodosByProject(todos []TodoItem) map[string][]TodoItem {
	m := map[string][]TodoItem{}
	for _, t := range todos {
		m[t.ProjName] = append(m[t.ProjName], t)
	}
	return m
}

// ── Project detail ─────────────────────────────────────────

func generateProjectDetail(dir string, r ProjectReport) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# %s\n\n", r.Project.Name))
	b.WriteString(fmt.Sprintf("- **Project Key**: `%s`\n", r.Project.Key))
	b.WriteString(fmt.Sprintf("- **Open Issues**: %d\n", r.TotalOpen))
	b.WriteString(fmt.Sprintf("- **Updated**: %s\n\n", time.Now().Format("2006-01-02 15:04")))

	if r.Project.LastAnalysisDate == "" {
		b.WriteString("> ⚠️ **This project has never been analyzed on SonarCloud.**\n\n")
	} else {
		b.WriteString(fmt.Sprintf("- **Last Analysis**: %s\n\n", r.Project.LastAnalysisDate))
	}

	analyzed := r.Project.LastAnalysisDate != ""
	score := calcHealthScore(r.Issues, analyzed)
	if analyzed {
		b.WriteString(fmt.Sprintf("## Health Score: %s %d/100\n\n", scoreIcon(score), score))
	}

	if r.QG != nil {
		status := r.QG.ProjectStatus.Status
		icon := "✅"
		if status == "ERROR" {
			icon = "❌"
		}
		b.WriteString(fmt.Sprintf("## Quality Gate: %s %s\n\n", icon, status))
	}

	if r.Measures != nil {
		b.WriteString("## Metrics\n\n")
		ms := r.Measures.Component.Measures
		b.WriteString("| Metric | Value |\n|--------|-------|\n")
		labels := map[string]string{
			"ncloc": "Lines of Code", "coverage": "Coverage", "duplicated_lines_density": "Duplication",
			"bugs": "Bugs", "vulnerabilities": "Vulnerabilities", "code_smells": "Smells",
			"sqale_index": "Debt(min)", "cognitive_complexity": "Cognitive Complexity",
		}
		for _, m := range ms {
			if l, ok := labels[m.Metric]; ok {
				b.WriteString(fmt.Sprintf("| %s | %s |\n", l, m.Value))
			}
		}
		b.WriteString("\n")
	}

	if len(r.Issues) > 0 {
		b.WriteString("## Issue Distribution\n\n")
		groups := classifyByPattern(r.Issues)
		type ps struct {
			name  string
			count int
		}
		var pss []ps
		for name, items := range groups {
			pss = append(pss, ps{name, len(items)})
		}
		sort.Slice(pss, func(i, j int) bool { return pss[i].count > pss[j].count })
		b.WriteString("| Category | Count |\n|----------|-------|\n")
		for _, p := range pss {
			b.WriteString(fmt.Sprintf("| %s | %d |\n", p.name, p.count))
		}
		b.WriteString("\n")

		b.WriteString("## Key Issues\n\n")
		b.WriteString("| Sev | File | Line | Rule | Message |\n")
		b.WriteString("|-----|------|------|------|--------|\n")
		limit := 20
		if len(r.Issues) < limit {
			limit = len(r.Issues)
		}
		for i := 0; i < limit; i++ {
			iss := r.Issues[i]
			f := shortFile(iss.Component)
			line := fmt.Sprintf("%d", iss.Line)
			if iss.Line == 0 {
				line = "-"
			}
			b.WriteString(fmt.Sprintf("| %s %s | `%s` | %s | `%s` | %s |\n",
				sevIcon[iss.Severity], iss.Severity, f, line, iss.Rule, iss.Message))
		}
		b.WriteString("\n")
	}

	writeFile(filepath.Join(dir, r.Project.Key+".md"), b.String())
}

// ── Utility ────────────────────────────────────────────────

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", path, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "  -> %s\n", path)
}

// ── Main ───────────────────────────────────────────────────

func main() {
	flag.Parse()

	// Load config file, CLI flags override
	resolveConfig()

	if orgKey == "" {
		fmt.Fprintln(os.Stderr, "Error: -org is required (or set in ppmp.json)")
		flag.Usage()
		os.Exit(1)
	}
	if outputDir == "" {
		fmt.Fprintln(os.Stderr, "Error: -o is required (or set in ppmp.json)")
		flag.Usage()
		os.Exit(1)
	}

	os.MkdirAll(outputDir, 0755)

	fmt.Fprintf(os.Stderr, "Discovering projects in org: %s\n", orgKey)
	projects, err := fetchProjects(orgKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list projects: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Found %d projects\n\n", len(projects))

	var reports []ProjectReport
	for i, p := range projects {
		fmt.Fprintf(os.Stderr, "[%d/%d] %s ...", i+1, len(projects), p.Key)
		r, err := fetchProjectData(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, " SKIP: %v\n", err)
			continue
		}
		reports = append(reports, *r)
		fmt.Fprintf(os.Stderr, " %d issues\n", len(r.Issues))
	}

	if len(reports) == 0 {
		fmt.Fprintln(os.Stderr, "No project data fetched")
		os.Exit(1)
	}

	analyze(reports)
}
