<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/SonarCloud-enabled-blue?style=for-the-badge&logo=sonarcloud&logoColor=white" alt="SonarCloud">
  <img src="https://img.shields.io/badge/Obsidian-output-7C3AED?style=for-the-badge&logo=obsidian&logoColor=white" alt="Obsidian">
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
</p>

<h1 align="center">🔧 PPMP</h1>
<h3 align="center">Personal Project Management Process</h3>
<p align="center">基于 SonarCloud + GitHub API 的项目代码质量管理流水线</p>
<p align="center">
  <b>自动拉取 → 跨项目统计 → 问题归类 → Obsidian 知识库 → Agent 编码上下文</b>
</p>

---

## 💡 核心理念

> **从"事后修 Bug"变为"事前防 Bug"**

通过跨项目问题统计，识别高频共性问题，在 AI Agent 编码阶段就规避。不是修 Bug，是让 Bug 不出现。

## 🏗️ 工作流

```
  SonarCloud (N repos)           GitHub API (repo metadata)
         │                                │
         ▼                                ▼
  📊 sonar-report                🔍 GitHub工程规范检查
  (单项目报告)                   (README/License/CI/Docker)
         │                                │
         └────────────┬───────────────────┘
                      ▼
             ⚙️ sonar-analyze
             (批量分析 + 统计)
                      │
         ┌────────────┼────────────┐
         ▼            ▼            ▼
   📁 Obsidian    📸 snapshots   📈 趋势追踪
   (知识库输出)   (历史归档)     (自动追加)
```

## 🚀 快速开始

### ✅ 前置条件

| 依赖 | 必需 | 获取方式 |
|:-----|:----:|---------|
| **Go 1.21+** | ✅ | [go.dev/dl](https://go.dev/dl/) |
| **SonarCloud 账号** | ✅ | [sonarcloud.io](https://sonarcloud.io) 注册，项目需已运行过分析 |
| **GitHub PAT** | ⭐ | [生成 Token](https://github.com/settings/tokens)（启用仓库规范检查） |

> 💡 **没有 SonarCloud？** 先去 [sonarcloud.io](https://sonarcloud.io) 创建账号，导入你的 GitHub 仓库，跑一次分析。本工具读取的是 SonarCloud 的分析结果。

### 📦 Step 1 — 克隆 & 编译

```bash
git clone https://github.com/Tangyd893/Personal-Project-Management-Process.git
cd Personal-Project-Management-Process/sonarcloud-report-skill/Script
go build -o sonar-report main.go
go build -o sonar-analyze analyze.go
```

### ⚙️ Step 2 — 配置

在 `$HOME` 下创建 `.ppmp.json`：

```bash
# Linux / macOS
cat > ~/.ppmp.json << 'EOF'
{
  "org": "your-github-username",
  "output": "/path/to/your/obsidian/vault/PPMP",
  "github_token": "ghp_xxxxxxxxxxxx",
  "issues_limit": 500
}
EOF

# Windows PowerShell
@"
{
  "org": "your-github-username",
  "output": "D:\\path\\to\\your\\obsidian\\vault\\PPMP",
  "github_token": "ghp_xxxxxxxxxxxx",
  "issues_limit": 500
}
"@ | Out-File -Encoding utf8 $env:USERPROFILE\.ppmp.json
```

> 🤖 **让 AI 帮你生成配置：** 把下面这段提示词发给 ChatGPT / Copilot / Cursor：
>
> ```
> 帮我生成一个 PPMP 工具的配置文件 .ppmp.json，我的信息如下：
> - GitHub 用户名：（填你的）
> - Obsidian 知识库路径：（填你的，如 D:\workspace\MyMind\PPMP）
> - GitHub Token：（填你的，格式 ghp_xxxx）
> - SonarCloud 组织名：（通常和 GitHub 用户名一致）
>
> 输出 JSON 格式，我直接保存到 $HOME/.ppmp.json
> ```

| 字段 | 必填 | 说明 |
|:-----|:----:|------|
| `org` | ✅ | SonarCloud / GitHub 组织名（通常就是你的 GitHub 用户名） |
| `output` | ✅ | Obsidian 知识库输出目录的**绝对路径** |
| `github_token` | ⭐ | GitHub PAT，启用仓库规范检查（可留空） |
| `issues_limit` | ❌ | 每项目 Issue 上限（默认 500） |
| `base_url` | ❌ | SonarCloud 地址（默认 `https://sonarcloud.io`） |

### ▶️ Step 3 — 运行

```bash
# 零参数运行（自动读取 ~/.ppmp.json）
./sonar-analyze

# 单项目报告
./sonar-report -project your-org_your-project -o report.md

# 列出所有项目
./sonar-report -org your-org -list
```

## 📂 输出文件

| 文件 | 说明 | Agent 用途 |
|:-----|------|:----------|
| 📊 `项目质量总览.md` | 健康度排名 + 质量门 + 指标汇总 | 全局概览 |
| 📋 `问题归类统计.md` | 按规则聚合 Top 30 + 示例 | 问题分布 |
| ⚡ `高频问题清单.md` | 跨项目共性问题 + 编码检查清单 | **编码前必读** |
| ✅ `质量待办.md` | BLOCKER/CRITICAL checkbox 清单 | 修复追踪 |
| 📈 `质量趋势.md` | 历史数据（每次运行自动追加） | 趋势分析 |
| 🔄 `增量变化.md` | 与上次运行对比（新增/已解决） | 发现退化 |
| 📅 `周期报告.md` | 本周/本月变化摘要 | 周期回顾 |
| 🏗️ `GitHub工程规范.md` | README/License/CI/Docker 检查 | 规范补全 |
| 📁 `projects/<key>.md` | 单项目详细报告 | 定位问题 |

## ✨ 特性

### 🔄 增量对比

每次运行自动与上次快照对比：

```
📈 新增 Issue: 3 个（引入了新问题）
📉 已解决 Issue: 5 个（修复了旧问题）
➡️ 净变化: -2（质量改善）
```

### ✅ 待办状态持久化

`质量待办.md` 中手动勾选的 `[x]` 会在下次运行时保留。修复一个，勾一个，不丢失。

### 🏥 健康度评分

```
score = 100 - (BLOCKER×10 + CRITICAL×5 + MAJOR×2 + MINOR×1)
```

| 🟢 80-100 | 🟡 60-79 | 🟠 40-59 | 🔴 0-39 |
|:---------:|:--------:|:--------:|:-------:|
| 健康 | 需关注 | 需改进 | 严重 |

### 🏗️ GitHub 仓库规范检查

配置 `github_token` 后，自动扫描每个仓库：

| 检查项 | 分值 | 说明 |
|:------|:----:|------|
| 📄 README.md | 14 | 项目说明文档 |
| 📜 LICENSE | 14 | 开源许可证 |
| 🚫 .gitignore | 14 | Git 忽略规则 |
| ⚙️ CI config | 14 | GitHub Actions / GitLab CI |
| 🔒 Lock file | 14 | go.sum / package-lock.json |
| 🐳 Dockerfile | 14 | 容器化支持 |
| 🤝 CONTRIBUTING | 14 | 贡献指南 |

### 📅 周期报告

从趋势数据自动生成本周/本月对比摘要（需积累 2+ 条数据后自动生效）。

## 🤖 Agent 集成

### 编码前（3 步走）

```
1️⃣  读 高频问题清单.md → 提取 Agent Coding Checklist 作为约束
2️⃣  读 质量待办.md     → 检查该项目待修复的 BLOCKER/CRITICAL
3️⃣  读 projects/<key>.md → 了解项目当前质量状态
```

### 编码后

重跑 `sonar-report -project <key>` 对比 Issue 变化。

## ❓ 常见问题

| 问题 | 解决方案 |
|:-----|---------|
| `org is required` | 检查 `~/.ppmp.json` 是否存在，JSON 格式是否正确 |
| GitHub 检查全部低分 | `github_token` 未配置或 token 无效 |
| 项目显示"未分析" | 该仓库从未在 SonarCloud 上运行过分析，需先配置 |
| 输出中文文件名乱码 | 控制台显示问题，实际文件正常，用 Obsidian 打开即可 |
| `go: command not found` | 未安装 Go，去 [go.dev/dl](https://go.dev/dl/) 下载 |

## ⏰ 自动化（可选）

工具本身不包含调度，可按需配置定时执行：

<details>
<summary>🪟 Windows — Task Scheduler</summary>

```powershell
# 创建定时任务：每周一早 9 点
schtasks /create /tn "PPMP-Analyze" /tr "C:\path\to\sonar-analyze.exe" /sc weekly /d MON /st 09:00
```
</details>

<details>
<summary>🐧 Linux / macOS — crontab</summary>

```cron
# 每周一早 9 点
0 9 * * 1 /path/to/sonar-analyze
```
</details>

<details>
<summary>🔄 CI/CD — GitHub Actions</summary>

```yaml
# .github/workflows/ppmp.yml
name: PPMP Analysis
on:
  schedule:
    - cron: '0 1 * * 1'  # UTC, 周一早 9 点 CST
jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - run: cd sonarcloud-report-skill/Script && go build -o sonar-analyze analyze.go
      - run: ./sonar-analyze -org ${{ secrets.SONAR_ORG }} -o ./output -token ${{ secrets.SONAR_TOKEN }}
```
</details>

## 📖 参数参考

<details>
<summary>sonar-report（单项目报告）</summary>

```
-project <key>       Project key（如 Tangyd893_SnapShop）
-org <key>           Organization key（配合 -list）
-list                列出组织下所有项目
-token <token>       SonarCloud Token（私有项目必填）
-o <file>            输出文件路径
-format json|markdown  输出格式
-issues-limit <n>    Issue 上限（默认 100）
```
</details>

<details>
<summary>sonar-analyze（批量分析）</summary>

```
-org <key>           Organization key（或 ppmp.json）
-o <dir>             输出目录（或 ppmp.json）
-token <token>       SonarCloud Token（私有项目）
-github-token <t>    GitHub Token（启用仓库规范检查）
-issues-limit <n>    每项目 Issue 上限（默认 500）
```
</details>

## 📁 项目结构

```
PPMP/
├── 📖 README.md
├── 🚫 .gitignore
├── ⚙️ ppmp.json                          # 配置模板（不入库）
├── 📦 samples/                           # 示例输出
│   ├── 📊 项目质量总览.md
│   ├── 📋 问题归类统计.md
│   ├── ⚡ 高频问题清单.md
│   ├── ✅ 质量待办.md
│   ├── 📈 质量趋势.md
│   ├── 🔄 增量变化.md
│   ├── 📅 周期报告.md
│   ├── 🏗️ GitHub工程规范.md
│   └── 📁 projects/
│       └── your-org_my-project-a.md
└── 🔧 sonarcloud-report-skill/
    ├── 📝 SKILL.md
    └── 💻 Script/
        ├── go.mod
        ├── main.go                      # sonar-report
        └── analyze.go                   # sonar-analyze
```

---

<p align="center">
  <b>⭐ 觉得有用？点个 Star 支持一下！</b>
</p>
