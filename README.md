# PPMP - Personal Project Management Process

基于 SonarCloud + GitHub API 的项目代码质量管理流水线。自动拉取所有项目分析报告，归类统计问题，输出到 Obsidian 知识库，为 AI Agent 编码提供质量上下文。

## 核心理念

**从"事后修 Bug"变为"事前防 Bug"** — 通过跨项目问题统计，识别高频共性问题，在 Agent 编码阶段就规避。

## 工作流

```
SonarCloud (N repos)        GitHub API (repo metadata)
      │                              │
      ▼                              ▼
sonar-report (单项目报告)    GitHub工程规范检查
      │                              │
      └──────────┬───────────────────┘
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
                 │    ├── 周期报告.md        ← 本周/本月变化摘要
                 │    ├── GitHub工程规范.md  ← README/License/CI/Docker 检查
                 │    └── projects/          ← 各项目独立报告
                 │
                 └──► snapshots/YYYY-MM-DD/ ← 历史快照归档
```

## 快速开始

### 1. 编译

```bash
cd sonarcloud-report-skill/Script
go build -o sonar-report main.go
go build -o sonar-analyze analyze.go
```

### 2. 配置

在以下任一位置创建 `ppmp.json`（搜索顺序：`./` → `~/`）：

```json
{
  "org": "your-github-username",
  "output": "D:\\path\\to\\your\\obsidian\\vault\\PPMP",
  "github_token": "",
  "issues_limit": 500
}
```

| 字段 | 必填 | 说明 |
|------|------|------|
| `org` | ✅ | SonarCloud / GitHub 组织名 |
| `output` | ✅ | Obsidian 知识库输出目录 |
| `github_token` | ❌ | GitHub PAT（启用仓库规范检查） |
| `issues_limit` | ❌ | 每项目 Issue 上限（默认 500） |
| `base_url` | ❌ | SonarCloud 地址（默认 sonarcloud.io） |

### 3. 运行

```bash
# 零参数运行（读取 ppmp.json）
./sonar-analyze

# 或 CLI 参数覆盖
./sonar-analyze -org your-org -o /path/to/output

# 单项目报告
./sonar-report -project your-org_your-project -o report.md

# 列出所有项目
./sonar-report -org your-org -list
```

## 输出文件

| 文件 | 内容 | Agent 用途 |
|------|------|-----------|
| `项目质量总览.md` | 健康度排名 + 质量门 + 指标汇总 | 项目质量概览 |
| `问题归类统计.md` | 按规则聚合 Top 30 + 示例 | 了解问题分布 |
| `高频问题清单.md` | 跨项目共性问题 + 编码检查清单 | **编码前必读** |
| `质量待办.md` | BLOCKER/CRITICAL checkbox 清单 | 修复任务追踪 |
| `质量趋势.md` | 历史数据（每次 analyze 追加） | 质量变化趋势 |
| `增量变化.md` | 与上次运行对比（新增/已解决） | 发现退化 |
| `周期报告.md` | 本周/本月变化摘要 | 周期回顾 |
| `GitHub工程规范.md` | README/License/CI/Docker 检查 | 工程规范补全 |
| `projects/<key>.md` | 单项目详细报告 | 定位具体问题 |

## 特性

### 增量对比

每次运行自动与上次对比，输出 `增量变化.md`：
- 📈 新增 Issue（引入了新问题）
- 📉 已解决 Issue（修复了旧问题）
- 净变化趋势

### 待办状态持久化

`质量待办.md` 中手动勾选的 `[x]` 会在下次运行时保留。

### 健康度评分

`100 - (BLOCKER×10 + CRITICAL×5 + MAJOR×2 + MINOR×1)`

| 分段 | 含义 |
|------|------|
| 🟢 80-100 | 健康 |
| 🟡 60-79 | 需关注 |
| 🟠 40-59 | 需改进 |
| 🔴 0-39 | 严重 |

### GitHub 仓库规范检查

配置 `github_token` 后，自动检查每个仓库的工程规范：

| 检查项 | 说明 |
|--------|------|
| README.md | 项目说明文档 |
| LICENSE | 开源许可证 |
| .gitignore | Git 忽略规则 |
| CI config | GitHub Actions / GitLab CI / Jenkins |
| Lock file | go.sum / package-lock.json / Cargo.lock |
| Dockerfile | 容器化支持 |
| CONTRIBUTING.md | 贡献指南 |

每项 14 分，满分 100。

### 周期报告

从趋势数据自动生成本周/本月对比摘要（需积累 2+ 条数据）。

## Agent 集成

### 编码前

1. 读 `高频问题清单.md` → 提取 Checklist 作为约束
2. 读 `质量待办.md` → 检查该项目待修复项
3. 读 `projects/<key>.md` → 了解当前质量状态

### 编码后

重跑 `sonar-report -project <key>` 对比 Issue 变化。

## 即插即用指南

克隆本仓库后，按以下步骤操作：

```bash
# 1. 克隆
git clone https://github.com/Tangyd893/Personal-Project-Management-Process.git
cd Personal-Project-Management-Process/sonarcloud-report-skill/Script

# 2. 编译（需要 Go 1.21+）
go build -o sonar-report main.go
go build -o sonar-analyze analyze.go

# 3. 配置
# 在项目根目录或 $HOME 创建 ppmp.json
cp ../../ppmp.json ~/.ppmp.json
# 编辑 org 和 output 路径

# 4. 运行
./sonar-analyze
```

### 前置条件

| 条件 | 必需 | 说明 |
|------|------|------|
| Go 1.21+ | ✅ | 编译脚本用 |
| SonarCloud 账号 | ✅ | 项目需在 SonarCloud 上分析过 |
| GitHub PAT | ❌ | 启用仓库规范检查 |

### 常见问题

**Q: 运行报错 "org is required"**
A: 检查 `ppmp.json` 是否在当前目录或 `$HOME` 下，JSON 格式是否正确。

**Q: GitHub 检查全部 2/100**
A: `github_token` 未配置或 token 无效。在 SonarCloud 个人设置中生成 PAT。

**Q: 项目显示"未分析"**
A: 该项目从未在 SonarCloud 上运行过分析。需要先在 SonarCloud 中配置项目。

**Q: 输出目录中文文件名乱码**
A: 控制台显示问题，实际文件名正常。用文件管理器或 Obsidian 打开即可。

## 自动化（可选）

工具本身不包含调度，可按需配置：

- **Windows**: Task Scheduler
- **Linux/macOS**: crontab
- **CI/CD**: GitHub Actions

示例 crontab（每周一早 9 点）：
```cron
0 9 * * 1 /path/to/sonar-analyze
```

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
-org <key>           Organization key（或 ppmp.json）
-o <dir>             输出目录（或 ppmp.json）
-token <token>       SonarCloud Token
-github-token <t>    GitHub Token（启用仓库规范检查）
-issues-limit <n>    每项目 Issue 上限（默认 500）
```

## 示例输出

见 [`samples/`](samples/) 目录。

## 项目结构

```
PPMP/
├── README.md
├── .gitignore
├── ppmp.json                          # 配置文件模板
├── samples/                           # 示例输出
│   ├── 项目质量总览.md
│   ├── 问题归类统计.md
│   ├── 高频问题清单.md
│   ├── 质量待办.md
│   ├── 质量趋势.md
│   ├── 增量变化.md
│   ├── 周期报告.md
│   ├── GitHub工程规范.md
│   └── projects/
│       └── your-org_my-project-a.md
└── sonarcloud-report-skill/
    ├── SKILL.md                       # Agent Skill 文档
    └── Script/
        ├── go.mod
        ├── main.go                    # sonar-report
        └── analyze.go                 # sonar-analyze
```
