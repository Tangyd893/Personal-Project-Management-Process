# SonarCloud 分析报告

- **项目**: SnapShop
- **Project Key**: `Tangyd893_SnapShop`
- **最近分析时间**: 2026-05-29T09:31:39+0000
- **报告生成时间**: 2026-06-07 08:00:55

---

## 🚦 质量门状态

**状态: ✅ OK**

| 条件 | 阈值 | 实际值 | 状态 |
|------|------|--------|------|
| new_reliability_rating | GT 1 | 1 | ✅ |
| new_security_rating | GT 1 | 1 | ✅ |
| new_maintainability_rating | GT 1 | 1 | ✅ |
| new_duplicated_lines_density | GT 3 | 0.0 | ✅ |
| new_security_hotspots_reviewed | LT 100 | 100.0 | ✅ |

## 📊 核心指标

| 指标 | 值 | 参考标准 |
|------|-----|----------|
| 代码行数 | 4151 行 |  |
| 测试覆盖率 | N/A | ≥ 80% |
| 重复代码率 | 1.3 % | ≤ 5% |
| Bug 数 | 0 个 | 0 |
| 安全漏洞 | 14 个 | 0 |
| Code Smell | 146 个 |  |
| 安全热点 | 2 个 |  |
| 技术债务 | 814 分钟 |  |
| 认知复杂度 | 104 |  |

## 🐛 未解决 Issue

共 **160** 个未解决 Issue（显示前 100 个）

### 🔴 BLOCKER (14)

| 文件 | 行号 | 规则 | 描述 |
|------|------|------|------|
| `backend/snapshop-inventory/src/main/resources/application.yml` | 5 | `java:S6437` | Revoke and change this secret, as it is compromised. |
| `backend/snapshop-payment/src/main/resources/application.yml` | 6 | `java:S6437` | Revoke and change this secret, as it is compromised. |
| `backend/snapshop-payment/src/main/resources/application.yml` | 11 | `java:S6437` | Revoke and change this secret, as it is compromised. |
| `backend/snapshop-user/src/main/resources/application.yml` | 5 | `java:S6437` | Revoke and change this secret, as it is compromised. |
| `backend/snapshop-admin/src/main/resources/application.yml` | 8 | `java:S6437` | Revoke and change this secret, as it is compromised. |
| `backend/snapshop-auth/src/main/resources/application.yml` | 8 | `java:S6437` | Revoke and change this secret, as it is compromised. |
| `backend/snapshop-product/src/main/resources/application.yml` | 8 | `java:S6437` | Revoke and change this secret, as it is compromised. |
| `backend/snapshop-order/src/main/resources/application.yml` | 8 | `java:S6437` | Revoke and change this secret, as it is compromised. |
| `backend/snapshop-order/src/main/resources/application.yml` | 14 | `java:S6437` | Revoke and change this secret, as it is compromised. |
| `backend/snapshop-seckill/src/main/resources/application.yml` | 9 | `java:S6437` | Revoke and change this secret, as it is compromised. |
| `backend/snapshop-seckill/src/main/resources/application.yml` | 15 | `java:S6437` | Revoke and change this secret, as it is compromised. |
| `docker/mysql/init/02-data.sql` | 7 | `secrets:S8215` | Make sure this bcrypt password hash gets revoked, changed, and removed from the code. |
| `docker/mysql/init/02-data.sql` | 8 | `secrets:S8215` | Make sure this bcrypt password hash gets revoked, changed, and removed from the code. |
| `docker/mysql/init/02-data.sql` | 43 | `secrets:S8215` | Make sure this bcrypt password hash gets revoked, changed, and removed from the code. |

### 🟠 CRITICAL (5)

| 文件 | 行号 | 规则 | 描述 |
|------|------|------|------|
| `docker/mysql/init/01-schema.sql` | 77 | `plsql:S1192` | Define a constant instead of duplicating this literal 6 times. |
| `docker/mysql/init/01-schema.sql` | 15 | `plsql:S1192` | Define a constant instead of duplicating this literal 4 times. |
| `docker/mysql/init/02-data.sql` | 7 | `plsql:S1192` | Define a constant instead of duplicating this literal 4 times. |
| `docker/mysql/init/02-data.sql` | 17 | `plsql:S1192` | Define a constant instead of duplicating this literal 5 times. |
| `docker/mysql/init/02-data.sql` | 7 | `plsql:S1192` | Define a constant instead of duplicating this literal 3 times. |

### 🟡 MAJOR (81)

| 文件 | 行号 | 规则 | 描述 |
|------|------|------|------|
| `testing/api/smoke-test.sh` | 287 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 287 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 226 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 227 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 227 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 240 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 242 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 242 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 256 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 256 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 272 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 287 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 51 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/smoke-test.sh` | 54 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/smoke-test.sh` | 58 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/smoke-test.sh` | 62 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/smoke-test.sh` | 63 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/smoke-test.sh` | 67 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/smoke-test.sh` | 75 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 88 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/smoke-test.sh` | 92 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/smoke-test.sh` | 101 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/smoke-test.sh` | 116 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 154 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 173 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 175 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 175 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 189 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 192 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 192 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 206 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 208 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 208 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 214 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 214 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 225 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 225 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 226 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 290 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 293 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 293 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 293 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 306 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 306 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 310 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 312 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 312 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 315 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 315 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 331 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 331 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 336 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 348 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 359 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 361 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 361 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/smoke-test.sh` | 404 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `frontend/snapshop-web/src/views/LoginPage.vue` | 59 | `css:S7924` | Text does not meet the minimal contrast requirement with its background. |
| `frontend/snapshop-web/src/views/LoginPage.vue` | 6 | `Web:S6853` | A form label must be associated with a control and have accessible text. |
| `frontend/snapshop-web/src/views/LoginPage.vue` | 10 | `Web:S6853` | A form label must be associated with a control and have accessible text. |
| `frontend/snapshop-web/src/views/SeckillDetail.vue` | 77 | `typescript:S3358` | Extract this nested ternary operation into an independent statement. |
| `frontend/snapshop-web/src/views/SeckillDetail.vue` | 113 | `css:S7924` | Text does not meet the minimal contrast requirement with its background. |
| `testing/api/payment-api-test.sh` | 136 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/payment-api-test.sh` | 136 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/payment-api-test.sh` | 35 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/payment-api-test.sh` | 39 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/payment-api-test.sh` | 40 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/payment-api-test.sh` | 44 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/payment-api-test.sh` | 49 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/payment-api-test.sh` | 53 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/payment-api-test.sh` | 62 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/payment-api-test.sh` | 76 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/payment-api-test.sh` | 122 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/payment-api-test.sh` | 122 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/payment-api-test.sh` | 30 | `shelldre:S7679` | Assign this positional parameter to a local variable. |
| `testing/api/payment-api-test.sh` | 167 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/payment-api-test.sh` | 167 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/payment-api-test.sh` | 204 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/payment-api-test.sh` | 204 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/payment-api-test.sh` | 215 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |
| `testing/api/payment-api-test.sh` | 215 | `shelldre:S7688` | Use '[[' instead of '[' for conditional tests. The '[[' construct is safer and more feature-rich. |

## 📈 分析历史

| 时间 | 版本 |
|------|------|
| 2026-05-29T09:31:39+0000 | not provided |
| 2026-05-29T08:36:35+0000 | not provided |
| 2026-05-28T09:20:38+0000 | not provided |

---

_数据来源: https://sonarcloud.io/dashboard?id=Tangyd893_SnapShop_
