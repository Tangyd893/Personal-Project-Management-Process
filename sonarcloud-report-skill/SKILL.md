# PPMP - SonarCloud Quality Analysis Skill

通过 SonarCloud API 获取项目代码质量分析报告，支持单项目报告和跨项目批量统计分析。

## 工具

| 工具 | 用途 | 源码 |
|------|------|------|
| `sonar-report` | 单项目报告（指标 + Issue + 历史） | `Script/main.go` |
| `sonar-analyze` | 批量分析（跨项目统计 → Obsidian） | `Script/analyze.go` |

## 前置条件

- Go 1.21+
- SonarCloud 公开项目无需 Token；私有项目需提供

## 快速开始

```bash
# 1. 编译
cd sonarcloud-report-skill/Script
go build -o sonar-report.exe main.go
go build -o sonar-analyze.exe analyze.go

# 2. 发现项目
./sonar-report -org tangyd893 -list

# 3. 单项目报告
./sonar-report -project Tangyd893_SnapShop -o report.md

# 4. 批量分析 → Obsidian
./sonar-analyze -org tangyd893 -o "D:\workspace\MyMind\PPMP"
```

## Agent 使用流程

### 场景 A: 编码前查看高频问题

用户要求编写代码时，Agent 应先读取 `高频问题清单.md`，将 Agent Coding Checklist 作为生成代码的约束条件。

```
读取 D:\workspace\MyMind\PPMP\高频问题清单.md
→ 提取 "Agent Coding Checklist" 部分
→ 在生成代码时逐项检查
```

### 场景 B: 项目质量审查

用户询问某项目质量时：

```
读取 D:\workspace\MyMind\PPMP\projects\<project-key>.md
→ 展示质量门状态、核心指标、关键 Issue
→ 如需更新，运行 sonar-report 拉取最新数据
```

### 场景 C: 跨项目质量对比

```
读取 D:\workspace\MyMind\PPMP\项目质量总览.md
→ 对比各项目质量门、Issue 数、技术债务
→ 识别需要重点关注的项目
```

### 场景 D: 刷新分析数据

```bash
# 运行 sonar-analyze 重新生成所有报告
cd sonarcloud-report-skill/Script && ./sonar-analyze -org tangyd893 -o "D:\workspace\MyMind\PPMP"
```

## 参数说明

### sonar-report

| 参数 | 说明 |
|------|------|
| `-project <key>` | SonarCloud project key |
| `-org <key>` | Organization key（配合 `-list`） |
| `-list` | 列出组织下所有项目 |
| `-token <t>` | SonarCloud Token |
| `-o <file>` | 输出文件路径 |
| `-format json\|markdown` | 输出格式 |
| `-issues-limit <n>` | Issue 上限（默认 100） |

### sonar-analyze

| 参数 | 说明 |
|------|------|
| `-org <key>` | Organization key（必填） |
| `-o <dir>` | 输出目录（必填） |
| `-token <t>` | SonarCloud Token |
| `-issues-limit <n>` | 每项目 Issue 上限（默认 500） |

## SonarCloud API 端点

| 端点 | 用途 |
|------|------|
| `GET /api/components/search` | 项目列表 |
| `GET /api/measures/component` | 项目指标 |
| `GET /api/qualitygates/project_status` | 质量门 |
| `GET /api/issues/search` | Issue 列表 |
| `GET /api/project_analyses/search` | 分析历史 |

## Project Key 格式

SonarCloud project key 通常为 `组织名_仓库名`，如：
- `Tangyd893_SnapShop`
- `Tangyd893_ERP-Go`
- `Tangyd893_HIS-Go`

如不确定 key，用 `-list` 一次性列出所有项目。
