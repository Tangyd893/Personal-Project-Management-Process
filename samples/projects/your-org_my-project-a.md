# my-project-a

- **Project Key**: `your-org_my-project-a`
- **Open Issues**: 8
- **Updated**: 2026-06-07 08:00

- **Last Analysis**: 2026-06-07T06:30:00+0000

## Health Score: 🟢 85/100

## Quality Gate: ✅ OK

## Metrics

| Metric | Value |
|--------|-------|
| Lines of Code | 5000 |
| Coverage | 72.0 |
| Duplication | 3.2 |
| Bugs | 0 |
| Vulnerabilities | 1 |
| Smells | 7 |
| Debt(min) | 120 |
| Cognitive Complexity | 45 |

## Issue Distribution

| Category | Count |
|----------|-------|
| Java source | 5 |
| Shell script | 2 |
| YAML config | 1 |

## Key Issues

| Sev | File | Line | Rule | Message |
|-----|------|------|------|--------|
| 🟡 MAJOR | `src/main/java/com/example/Service.java` | 42 | `java:S1192` | Define a constant instead of duplicating this literal. |
| 🟡 MAJOR | `scripts/deploy.sh` | 15 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. |
| 🔵 MINOR | `src/main/resources/application.yml` | 3 | `yaml:S1135` | Track uses of "TODO" tags. |
