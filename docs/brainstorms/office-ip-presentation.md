---
date: 2026-09-20
topic: office-ip-presentation
---

# Office IP Presentation

## Summary

让 `FindPlace` 支持一个可配置的 IP → 标签映射。配置中的项既可以是单 IP（如 `1.2.3.4:Beijing`），也可以是 CIDR（如 `10.0.0.0/24:OfficeNet`）。命中时直接返回标签；未命中时保持现有 `ipip.FindCity` 行为。

## Problem Frame

`pkg/web/template.go` 的 `FindPlace` 走 `ipip.FindCity` 做 IP 归属地查询，对员工居家或出差网络没问题，但对办公室静态 IP 段只能拿到城市粒度，没法区分北京和苏州两个办公点。管理员没有途径告诉系统"这些是办公室"，所以也无法在状态页面上把办公室连接从外部连接里突出出来。

## Requirements

- R1. 在配置中新增一项办公室 IP 映射，支持单 IP 和 CIDR 两种形式，使用冒号分隔 IP 与标签、英文逗号分隔多条目。
- R2. `FindPlace` 在收到映射中命中的 IP 时，直接返回对应标签，跳过 `ipip.FindCity`。
- R3. `FindPlace` 在 IP 未命中映射时，保持现有行为不变（`ipip.FindCity` 结果，缺失时返回 `[未知地区]`）。
- R4. 配置项缺省或未设置时，行为与当前一致。

## Acceptance Examples

- AE1. **Covers R1, R2.** Given 配置含 `1.2.3.4:Beijing`，when `FindPlace("1.2.3.4")`，then 返回 `"Beijing"`。
- AE2. **Covers R1, R2.** Given 配置含 `10.0.0.0/24:OfficeNet`，when `FindPlace("10.0.0.5")`，then 返回 `"OfficeNet"`。
- AE3. **Covers R2, R3.** Given 配置含 `1.2.3.4:Beijing`，when `FindPlace("8.8.8.8")`，then 走 `ipip.FindCity("8.8.8.8")` 的现有结果。
- AE4. **Covers R4.** Given 配置项未设置，when 任意 `FindPlace` 调用，then 与升级前行为完全一致。

## Success Criteria

- 管理员能在不重启服务的情况下不增加复杂度（沿用现有 envconfig 加载模式），通过单一配置项声明多个办公室 IP 段及其显示标签。
- 状态页面上，命中办公室映射的连接以标签呈现，便于一眼区分；非命中行外观与升级前无差。

## Scope Boundaries

- 不动 `IsOfficeIP`（`pkg/web/template.go` 的 TODO stub），即便它逻辑上相关。
- 不引入 CIDR 之外的匹配方式（如地域区间、自定义通配）。
- 不做配置热更新（沿用 init 加载语义）。
- 不基于 `ipip` 反查结果回填未配置 IP 的标签。

## Key Decisions

- 命中时只返回标签、不附加 `ipip` 结果：用户明确选择，避免"Beijing/北京北京"这种冗余。
- 格式采用 `ip:label` 列表，与现有 envconfig map 类型语法对齐，复用仓库配置惯例。
- 匹配粒度覆盖单 IP 和 CIDR：办公室既有公网固定 IP，也有 NAT 后的内网段，需要都支持。

## Dependencies / Assumptions

- 假设当前 `FindPlace` 调用方不需要区分"命中映射"和"命中 ipip"的标记信息 — 如果后续模板想做差异化样式，需另行扩展返回类型。
