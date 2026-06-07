# PPMP - Personal Project Management Process

基于 SonarCloud API 的项目代码质量管理流水线。自动拉取所有项目分析报告，归类统计问题，输出到 Obsidian 知识库，为 AI Agent 编码提供质量上下文。

## 核心理念

**从"事后修 Bug"变为"事前防 Bug"** — 通过跨项目问题统计，识别高频共性问题，在 Agent 编码阶段就规避。

## 工作流

```
SonarCloud (N repos)
      │
      ▼
sonar-report (单项目报告)
      │
      ▼
sonar-analyze (批量分析 + 统计)
      │
      ├──► Obsidian 知识库
      │    ├── 项目质量总览.md    ← 健康度排名 + 各项目指标
      │    ├── 问题归类统计.md    ← 按规则聚合 Top N
      │    ├── 高频问题清单.md    ← 跨项目共性问题 + Agent Checklist
      │    ├── 质量待办.md        ← BLOCKER/CRITICAL 待修复清单
      │    ├── 质量趋势.md        ← 历史数据（自动追加）
      │    ├── 增量变化.md        ← 新增/已解决 Issue 对比
      │    └── projects/          ← 各项目独立报告
      │
      └──► snapshots/YYYY-MM-DD/ ← 历史快照归档
```

## 快速开始

```bash
# 编译
cd sonarcloud-report-skill/Script
go build -o sonar-report main.go
go build -o sonar-analyze analyze.go

# 方式 1: 配置文件（推荐）
cp ppmp.json ~/.ppmp.json   # 编辑 org 和 output
./sonar-analyze              # 零参数运行

# 方式 2: CLI 参数
./sonar-analyze -org your-org -o "/path/to/obsidian/PPMP"

# 单项目报告
./sonar-report -project your-org_your-project -o report.md

# 列出组织下所有项目
./sonar-report -org your-org -list
```

## 输出说明

| 文件 | 内容 | Agent 用途 |
|------|------|-----------|
| `项目质量总览.md` | 健康度排名 + 质量门 + 指标汇总 | 项目质量概览 |
| `问题归类统计.md` | 按规则聚合 Top 30 + 示例 | 了解问题分布 |
| `高频问题清单.md` | 跨项目共性问题 + 编码检查清单 | **编码前必读** |
| `质量待办.md` | BLOCKER/CRITICAL checkbox 清单（状态持久化） | 修复任务追踪 |
| `质量趋势.md` | 历史数据（每次 analyze 追加） | 质量变化趋势 |
| `增量变化.md` | 与上次运行对比（新增/已解决） | 发现退化 |
| `projects/<key>.md` | 单项目详细报告 | 定位具体问题 |

## 特性

### 配置文件

在 `./ppmp.json` 或 `~/.ppmp.json` 中保存常用参数，避免每次手打：

```json
{
  "org": "your-org",
  "output": "D:\\workspace\\MyMind\\PPMP",
  "token": "",
  "issues_limit": 500
}
```

CLI 参数优先于配置文件。

### 增量对比

每次运行自动与上次对比，输出 `增量变化.md`：
- 📈 新增 Issue（引入了新问题）
- 📉 已解决 Issue（修复了旧问题）
- 净变化趋势

### 待办状态持久化

`质量待办.md` 中手动勾选的 `[x]` 会在下次运行时保留，不会被覆盖。

### 健康度评分

0-100 分算法：`100 - (BLOCKER×10 + CRITICAL×5 + MAJOR×2 + MINOR×1)`

| 分段 | 含义 |
|------|------|
| 🟢 80-100 | 健康 |
| 🟡 60-79 | 需关注 |
| 🟠 40-59 | 需改进 |
| 🔴 0-39 | 严重 |
| ⚪ - | 未分析 |

## Agent 集成

### 编码前

Agent 在为某项目编写代码时，应先读取：
1. `高频问题清单.md` → 提取 Agent Coding Checklist 作为约束
2. `质量待办.md` → 检查该项目是否有待修复的 BLOCKER/CRITICAL
3. `projects/<key>.md` → 了解项目当前质量状态

### 编码后

提交代码后，可重新运行 `sonar-report -project <key>` 对比 Issue 数变化。

## 自动化（可选）

工具本身不包含自动化调度，但可按需配置：

- **Windows**: Task Scheduler 定时执行 `sonar-analyze`
- **Linux/macOS**: crontab 定时执行
- **CI/CD**: GitHub Actions / GitLab CI 中集成

## 参数参考

### sonar-report

```
-project <key>       Project key
-org <key>           Organization key (配合 -list)
-list                列出组织下所有项目
-token <token>       SonarCloud Token（私有项目）
-o <file>            输出文件
-format json|markdown
-issues-limit <n>    Issue 上限（默认 100）
```

### sonar-analyze

```
-org <key>           Organization key（必填，或在 ppmp.json 中配置）
-o <dir>             输出目录（必填，或在 ppmp.json 中配置）
-token <token>       SonarCloud Token
-issues-limit <n>    每项目 Issue 上限（默认 500）
```

## 示例输出

见 [`samples/`](samples/) 目录，包含各文件的模板示例。

## 项目结构

```
PPMP/
├── README.md
├── .gitignore
├── ppmp.json                          # 配置文件模板
├── samples/                           # 示例输出
│   ├── 项目质量总览.md
│   ├── 高频问题清单.md
│   ├── 质量待办.md
│   ├── 质量趋势.md
│   └── projects/
│       └── your-org_my-project-a.md
└── sonarcloud-report-skill/
    ├── SKILL.md                       # Agent Skill 文档
    └── Script/
        ├── go.mod
        ├── main.go                    # sonar-report 源码
        └── analyze.go                 # sonar-analyze 源码
```
