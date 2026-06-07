package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ── CLI flags ──────────────────────────────────────────────

var (
	projectKey  string
	orgKey      string
	token       string
	outputFile  string
	baseURL     string
	issuesLimit int
	format      string
	listMode    bool
)

func init() {
	flag.StringVar(&projectKey, "project", "", "SonarCloud project key (e.g. Tangyd893_SnapShop)")
	flag.StringVar(&orgKey, "org", "", "SonarCloud organization key (e.g. tangyd893)")
	flag.StringVar(&token, "token", "", "SonarCloud User Token (required for private projects)")
	flag.StringVar(&outputFile, "o", "", "Output file path (default: stdout)")
	flag.StringVar(&baseURL, "base-url", "https://sonarcloud.io", "SonarCloud base URL")
	flag.IntVar(&issuesLimit, "issues-limit", 100, "Max issues to fetch")
	flag.StringVar(&format, "format", "markdown", "Output format: markdown or json")
	flag.BoolVar(&listMode, "list", false, "List all projects in the organization")
}

// ── Shared types ───────────────────────────────────────────

type ComponentMeasure struct {
	Metric string `json:"metric"`
	Value  string `json:"value"`
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
	Status    string `json:"status"`
}

type AnalysisItem struct {
	Key            string `json:"key"`
	Date           string `json:"date"`
	ProjectVersion string `json:"projectVersion,omitempty"`
	Events         []struct {
		Category string `json:"category"`
		Name     string `json:"name"`
	} `json:"events,omitempty"`
}

type ProjectItem struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	Qualifier string `json:"qualifier"`
	Visibility string `json:"visibility"`
	LastAnalysisDate string `json:"lastAnalysisDate,omitempty"`
}

// ── API response types ─────────────────────────────────────

type MeasuresResponse struct {
	Component struct {
		Key       string             `json:"key"`
		Name      string             `json:"name"`
		Qualifier string             `json:"qualifier"`
		Measures  []ComponentMeasure `json:"measures"`
	} `json:"component"`
}

type QualityGateResponse struct {
	ProjectStatus struct {
		Status     string `json:"status"`
		Conditions []struct {
			Status         string `json:"status"`
			MetricKey      string `json:"metricKey"`
			Comparator     string `json:"comparator"`
			PeriodIndex    int    `json:"periodIndex"`
			ErrorThreshold string `json:"errorThreshold"`
			ActualValue    string `json:"actualValue"`
		} `json:"conditions"`
	} `json:"projectStatus"`
}

type IssuesResponse struct {
	Total  int         `json:"total"`
	Issues []IssueItem `json:"issues"`
}

type AnalysesResponse struct {
	Analyses []AnalysisItem `json:"analyses"`
	Paging   struct {
		Total  int `json:"total"`
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"paging"`
}

type ProjectsResponse struct {
	Components []ProjectItem `json:"components"`
	Paging     struct {
		Total  int `json:"total"`
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"paging"`
}

// ── HTTP helper ────────────────────────────────────────────

func apiGet(path string, params url.Values) ([]byte, error) {
	u := baseURL + path
	if params != nil {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	if token != "" {
		req.SetBasicAuth(token, "")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("api %s returned %d: %s", path, resp.StatusCode, string(body))
	}
	return body, nil
}

// ── Project discovery ──────────────────────────────────────

func fetchProjects(org string) ([]ProjectItem, error) {
	var allProjects []ProjectItem
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
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, fmt.Errorf("unmarshal projects: %w", err)
		}
		allProjects = append(allProjects, resp.Components...)
		if len(allProjects) >= resp.Paging.Total {
			break
		}
		page++
	}
	return allProjects, nil
}

func guessOrgFromProject(project string) string {
	// Project key format is typically org_repo, e.g. Tangyd893_SnapShop
	if idx := strings.Index(project, "_"); idx > 0 {
		return strings.ToLower(project[:idx])
	}
	return strings.ToLower(project)
}

// ── Report data fetchers ───────────────────────────────────

const metricKeys = "ncloc,coverage,duplicated_lines_density,bugs,vulnerabilities,code_smells,security_hotspots,sqale_index,cognitive_complexity"

func fetchMeasures() (*MeasuresResponse, error) {
	params := url.Values{
		"component":  {projectKey},
		"metricKeys": {metricKeys},
	}
	data, err := apiGet("/api/measures/component", params)
	if err != nil {
		return nil, err
	}
	var resp MeasuresResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal measures: %w", err)
	}
	return &resp, nil
}

func fetchQualityGate() (*QualityGateResponse, error) {
	params := url.Values{"projectKey": {projectKey}}
	data, err := apiGet("/api/qualitygates/project_status", params)
	if err != nil {
		return nil, err
	}
	var resp QualityGateResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal qualitygate: %w", err)
	}
	return &resp, nil
}

func fetchIssues() (*IssuesResponse, error) {
	params := url.Values{
		"componentKeys": {projectKey},
		"ps":            {fmt.Sprintf("%d", issuesLimit)},
		"resolved":      {"false"},
		"s":             {"SEVERITY"},
		"asc":           {"false"},
	}
	data, err := apiGet("/api/issues/search", params)
	if err != nil {
		return nil, err
	}
	var resp IssuesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal issues: %w", err)
	}
	return &resp, nil
}

func fetchAnalyses() (*AnalysesResponse, error) {
	params := url.Values{
		"project": {projectKey},
		"ps":      {"10"},
	}
	data, err := apiGet("/api/project_analyses/search", params)
	if err != nil {
		return nil, err
	}
	var resp AnalysesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal analyses: %w", err)
	}
	return &resp, nil
}

// ── Metric display helpers ─────────────────────────────────

type metricInfo struct {
	Key   string
	Label string
	Unit  string
	Good  string
}

var metricDisplay = []metricInfo{
	{"ncloc", "代码行数", "行", ""},
	{"coverage", "测试覆盖率", "%", "≥ 80%"},
	{"duplicated_lines_density", "重复代码率", "%", "≤ 5%"},
	{"bugs", "Bug 数", "个", "0"},
	{"vulnerabilities", "安全漏洞", "个", "0"},
	{"code_smells", "Code Smell", "个", ""},
	{"security_hotspots", "安全热点", "个", ""},
	{"sqale_index", "技术债务", "分钟", ""},
	{"cognitive_complexity", "认知复杂度", "", ""},
}

func getMeasure(measures []ComponentMeasure, key string) string {
	for _, m := range measures {
		if m.Metric == key {
			return m.Value
		}
	}
	return "N/A"
}

func severityIcon(s string) string {
	switch strings.ToUpper(s) {
	case "BLOCKER":
		return "🔴"
	case "CRITICAL":
		return "🟠"
	case "MAJOR":
		return "🟡"
	case "MINOR":
		return "🔵"
	case "INFO":
		return "⚪"
	default:
		return "❓"
	}
}

func statusIcon(s string) string {
	switch strings.ToUpper(s) {
	case "OK", "PASSED":
		return "✅"
	case "ERROR", "FAILED":
		return "❌"
	case "WARN":
		return "⚠️"
	default:
		return "❓"
	}
}

// ── Project list formatters ────────────────────────────────

func formatProjectListMarkdown(org string, projects []ProjectItem) string {
	var b strings.Builder

	b.WriteString("# SonarCloud 项目列表\n\n")
	b.WriteString(fmt.Sprintf("- **组织**: %s\n", org))
	b.WriteString(fmt.Sprintf("- **项目总数**: %d\n", len(projects)))
	b.WriteString(fmt.Sprintf("- **生成时间**: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	b.WriteString("| # | Project Key | 项目名 | 可见性 | 最近分析 |\n")
	b.WriteString("|---|-------------|--------|--------|----------|\n")
	for i, p := range projects {
		lastAnalysis := p.LastAnalysisDate
		if lastAnalysis == "" {
			lastAnalysis = "_从未分析_"
		}
		b.WriteString(fmt.Sprintf("| %d | `%s` | %s | %s | %s |\n",
			i+1, p.Key, p.Name, p.Visibility, lastAnalysis))
	}
	b.WriteString("\n---\n\n")
	b.WriteString(fmt.Sprintf("_使用 `-project <key>` 获取任意项目的详细报告_\n"))

	return b.String()
}

func formatProjectListJSON(org string, projects []ProjectItem) string {
	report := struct {
		Organization string        `json:"organization"`
		Count        int           `json:"count"`
		GeneratedAt  string        `json:"generated_at"`
		Projects     []ProjectItem `json:"projects"`
	}{
		Organization: org,
		Count:        len(projects),
		GeneratedAt:  time.Now().Format(time.RFC3339),
		Projects:     projects,
	}
	data, _ := json.MarshalIndent(report, "", "  ")
	return string(data)
}

// ── Full report formatters ─────────────────────────────────

func formatMarkdown(m *MeasuresResponse, qg *QualityGateResponse, issues *IssuesResponse, analyses *AnalysesResponse) string {
	var b strings.Builder

	// Header
	b.WriteString("# SonarCloud 分析报告\n\n")
	b.WriteString(fmt.Sprintf("- **项目**: %s\n", m.Component.Name))
	b.WriteString(fmt.Sprintf("- **Project Key**: `%s`\n", m.Component.Key))
	if len(analyses.Analyses) > 0 {
		b.WriteString(fmt.Sprintf("- **最近分析时间**: %s\n", analyses.Analyses[0].Date))
	}
	b.WriteString(fmt.Sprintf("- **报告生成时间**: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// Quality Gate
	b.WriteString("---\n\n## 🚦 质量门状态\n\n")
	if qg != nil {
		b.WriteString(fmt.Sprintf("**状态: %s %s**\n\n", statusIcon(qg.ProjectStatus.Status), qg.ProjectStatus.Status))
		if len(qg.ProjectStatus.Conditions) > 0 {
			b.WriteString("| 条件 | 阈值 | 实际值 | 状态 |\n")
			b.WriteString("|------|------|--------|------|\n")
			for _, c := range qg.ProjectStatus.Conditions {
				b.WriteString(fmt.Sprintf("| %s | %s %s | %s | %s |\n",
					c.MetricKey, c.Comparator, c.ErrorThreshold, c.ActualValue, statusIcon(c.Status)))
			}
			b.WriteString("\n")
		}
	}

	// Core metrics
	b.WriteString("## 📊 核心指标\n\n")
	b.WriteString("| 指标 | 值 | 参考标准 |\n")
	b.WriteString("|------|-----|----------|\n")
	for _, mi := range metricDisplay {
		val := getMeasure(m.Component.Measures, mi.Key)
		display := val
		if mi.Unit != "" && val != "N/A" {
			display = val + " " + mi.Unit
		}
		b.WriteString(fmt.Sprintf("| %s | %s | %s |\n", mi.Label, display, mi.Good))
	}
	b.WriteString("\n")

	// Issues
	b.WriteString("## 🐛 未解决 Issue\n\n")
	b.WriteString(fmt.Sprintf("共 **%d** 个未解决 Issue", issues.Total))
	if issues.Total > issuesLimit {
		b.WriteString(fmt.Sprintf("（显示前 %d 个）", issuesLimit))
	}
	b.WriteString("\n\n")

	if len(issues.Issues) > 0 {
		groups := map[string][]IssueItem{}
		for _, iss := range issues.Issues {
			groups[iss.Severity] = append(groups[iss.Severity], iss)
		}

		order := []string{"BLOCKER", "CRITICAL", "MAJOR", "MINOR", "INFO"}
		for _, sev := range order {
			items, ok := groups[sev]
			if !ok {
				continue
			}
			b.WriteString(fmt.Sprintf("### %s %s (%d)\n\n", severityIcon(sev), sev, len(items)))
			b.WriteString("| 文件 | 行号 | 规则 | 描述 |\n")
			b.WriteString("|------|------|------|------|\n")
			for _, iss := range items {
				comp := iss.Component
				if idx := strings.LastIndex(comp, ":"); idx >= 0 {
					comp = comp[idx+1:]
				}
				line := fmt.Sprintf("%d", iss.Line)
				if iss.Line == 0 {
					line = "-"
				}
				b.WriteString(fmt.Sprintf("| `%s` | %s | `%s` | %s |\n", comp, line, iss.Rule, iss.Message))
			}
			b.WriteString("\n")
		}
	} else {
		b.WriteString("_无未解决 Issue_ ✅\n\n")
	}

	// Analysis history
	b.WriteString("## 📈 分析历史\n\n")
	if len(analyses.Analyses) > 0 {
		b.WriteString("| 时间 | 版本 |\n")
		b.WriteString("|------|------|\n")
		for _, a := range analyses.Analyses {
			ver := a.ProjectVersion
			if ver == "" {
				ver = "-"
			}
			b.WriteString(fmt.Sprintf("| %s | %s |\n", a.Date, ver))
		}
		b.WriteString("\n")
	} else {
		b.WriteString("_无分析历史记录_\n\n")
	}

	// Footer
	b.WriteString("---\n\n")
	b.WriteString(fmt.Sprintf("_数据来源: %s/dashboard?id=%s_\n", baseURL, projectKey))

	return b.String()
}

func formatJSON(m *MeasuresResponse, qg *QualityGateResponse, issues *IssuesResponse, analyses *AnalysesResponse) string {
	measures := make(map[string]string)
	for _, mm := range m.Component.Measures {
		measures[mm.Metric] = mm.Value
	}

	report := struct {
		Project     string              `json:"project"`
		GeneratedAt string              `json:"generated_at"`
		QualityGate *QualityGateResponse `json:"quality_gate"`
		Measures    map[string]string    `json:"measures"`
		IssuesCount int                  `json:"issues_count"`
		Issues      []IssueItem          `json:"issues,omitempty"`
		Analyses    []AnalysisItem       `json:"analyses,omitempty"`
	}{
		Project:     m.Component.Key,
		GeneratedAt: time.Now().Format(time.RFC3339),
		QualityGate: qg,
		Measures:    measures,
		IssuesCount: issues.Total,
		Issues:      issues.Issues,
		Analyses:    analyses.Analyses,
	}

	data, _ := json.MarshalIndent(report, "", "  ")
	return string(data)
}

// ── Output helper ──────────────────────────────────────────

func output(content string) {
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to write file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "📄 Report saved to: %s\n", outputFile)
	} else {
		fmt.Print(content)
	}
}

// ── Main ───────────────────────────────────────────────────

func main() {
	flag.Parse()

	// ── List mode: discover projects ──
	if listMode {
		org := orgKey
		if org == "" && projectKey != "" {
			org = guessOrgFromProject(projectKey)
		}
		if org == "" {
			fmt.Fprintln(os.Stderr, "Error: -org or -project is required with -list")
			flag.Usage()
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "🔍 Listing projects in organization: %s\n", org)
		projects, err := fetchProjects(org)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to list projects: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "✅ Found %d projects\n", len(projects))

		var report string
		switch format {
		case "json":
			report = formatProjectListJSON(org, projects)
		default:
			report = formatProjectListMarkdown(org, projects)
		}
		output(report)
		return
	}

	// ── Report mode ──
	if projectKey == "" {
		fmt.Fprintln(os.Stderr, "Error: -project is required")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "🔍 Fetching SonarCloud data for project: %s\n", projectKey)

	type result struct {
		measures *MeasuresResponse
		qg       *QualityGateResponse
		issues   *IssuesResponse
		analyses *AnalysesResponse
		err      error
		label    string
	}

	ch := make(chan result, 4)

	go func() {
		m, err := fetchMeasures()
		ch <- result{measures: m, err: err, label: "measures"}
	}()
	go func() {
		qg, err := fetchQualityGate()
		ch <- result{qg: qg, err: err, label: "quality_gate"}
	}()
	go func() {
		iss, err := fetchIssues()
		ch <- result{issues: iss, err: err, label: "issues"}
	}()
	go func() {
		a, err := fetchAnalyses()
		ch <- result{analyses: a, err: err, label: "analyses"}
	}()

	var (
		measures *MeasuresResponse
		qg       *QualityGateResponse
		issues   *IssuesResponse
		analyses *AnalysesResponse
	)

	for i := 0; i < 4; i++ {
		r := <-ch
		if r.err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to fetch %s: %v\n", r.label, r.err)
			os.Exit(1)
		}
		switch r.label {
		case "measures":
			measures = r.measures
		case "quality_gate":
			qg = r.qg
		case "issues":
			issues = r.issues
		case "analyses":
			analyses = r.analyses
		}
	}

	fmt.Fprintf(os.Stderr, "✅ Data fetched successfully\n")

	var report string
	switch format {
	case "json":
		report = formatJSON(measures, qg, issues, analyses)
	default:
		report = formatMarkdown(measures, qg, issues, analyses)
	}
	output(report)
}
