# 任务调度框架流程分析

## 一、核心组件

### 1. 任务调度器 (Sched)
- **位置**: `aio-scheduler/cmd/service/task/sche.go`
- **功能**: 定时轮询待执行任务并启动执行

### 2. 任务执行器 (Exec)
- **位置**: `aio-scheduler/cmd/service/task/sche.go`
- **功能**: 执行任务的所有阶段和步骤

### 3. 用户操作接口
- **Stage操作**: `POST /aioux/stages/:id/operation`
- **Step操作**: `POST /aioux/steps/:id/operation`

## 二、状态定义

### Task状态
- `created`: 待启动
- `started`: 进行中
- `finished`: 已完成
- `failed`: 失败
- `terminated`: 已终止

### Stage状态
- `created`: 待执行
- `started`: 执行中
- `finished`: 已完成
- `failed`: 失败
- `wait-finished`: 等待完成（需要用户操作）
- `terminated`: 已终止

### Step状态
- `created`: 待执行
- `started`: 执行中
- `finished`: 已完成
- `failed`: 失败
- `useraction_blocked`: 用户操作阻塞（手动模式）
- `terminated`: 已终止

### Step模式
- `StepManual`: 手动模式（需要用户操作）
- `StepFailedStop`: 失败停止
- `StepFailedContinue`: 失败继续

## 三、调度流程

### 3.1 任务调度循环 (Sched)

```
启动调度器 (每秒轮询一次)
    ↓
查找状态为 created 的任务 (IdelTask)
    ↓
找到任务？
    ├─ 否 → 继续轮询
    └─ 是 → 更新任务状态为 started
            ↓
        启动 goroutine 执行任务
            ├─ preExec: 预处理（清理、初始化）
            └─ Exec: 执行任务
```

### 3.2 任务执行流程 (Exec)

```
开始执行任务
    ↓
遍历所有 Stages
    ↓
Stage 状态检查
    ├─ 不是 created → 跳过
    └─ 是 created → 更新 Stage 状态为 started
                    ↓
                遍历 Stage 的所有 Steps
                    ↓
                Step 状态检查
                    ├─ 不是 created → 跳过
                    └─ 是 created → 检查 Step 模式
                                    ↓
                                Step 模式检查
                                    ├─ StepManual → 更新为 useraction_blocked，返回等待用户操作
                                    └─ 其他模式 → 更新 Step 状态为 started
                                                    ↓
                                                获取 Step Handler (StepFactor)
                                                    ↓
                                                Handler 存在？
                                                    ├─ 否 → 标记 Step/Stage/Task 为 failed，返回
                                                    └─ 是 → 执行 Handler.Do()
                                                                ↓
                                                            执行成功？
                                                                ├─ 否 → 根据 Step.Mode 处理
                                                                │       ├─ StepFailedStop → 标记失败，终止任务
                                                                │       └─ StepFailedContinue → 标记 Step 失败，继续执行
                                                                └─ 是 → 执行 Handler.Check()
                                                                            ↓
                                                                        Check 通过？
                                                                            ├─ 否 → 根据 Step.Mode 处理
                                                                            │       ├─ StepFailedStop → 标记失败，终止任务
                                                                            │       └─ StepFailedContinue → 标记 Step 失败，继续执行
                                                                            └─ 是 → 标记 Step 为 finished
                                                                                    ↓
                                                                                更新实例状态
                                                                                    ↓
                                                                                继续下一个 Step
                    ↓
                所有 Steps 执行完成
                    ↓
                检查 Stage 模式
                    ├─ IsAutoStage → 标记 Stage 为 finished，继续下一个 Stage
                    ├─ StageFailedContinue → 检查是否有失败的 Step
                    │                       ├─ 有 → 标记 Stage 为 failed
                    │                       └─ 无 → 标记 Stage 为 finished
                    └─ 其他模式 → 检查是否为最后一个 Stage
                                    ├─ 是 → 标记 Stage 为 finished
                                    └─ 否 → 标记 Stage 为 wait-finished，返回等待用户操作
    ↓
所有 Stages 执行完成
    ↓
标记 Task 为 finished
```

### 3.3 用户操作流程

#### 3.3.1 Stage 操作

```
用户请求: POST /aioux/stages/:id/operation
    ↓
验证 Stage 状态为 wait-finished
    ↓
解析操作类型 (Action)
    ↓
操作类型
    ├─ continue → 标记 Stage 为 finished
    │             标记 Task 为 created (重新调度)
    │
    └─ abort → 标记 Stage 为 terminated
               标记 Task 为 terminated
               执行 TaskComplete
```

#### 3.3.2 Step 操作

```
用户请求: POST /aioux/steps/:id/operation
    ↓
验证 Step 状态为 failed 或 useraction_blocked
    ↓
解析操作类型 (Action)
    ↓
操作类型
    ├─ continue → 修改 Step 模式为 StepFailedStop
    │             (fallthrough 到 retry)
    │
    ├─ retry → 记录日志
    │         标记 Step 为 created
    │         标记 Stage 为 created
    │         标记 Task 为 created (重新调度)
    │
    ├─ ignore → 记录日志
    │           特殊 Step 类型 → 执行 TaskComplete
    │           标记 Step 为 finished
    │           标记 Stage 为 created
    │           标记 Task 为 created (重新调度)
    │
    └─ abort → 记录日志
               标记 Step 为 terminated
               标记 Stage 为 terminated
               标记 Task 为 terminated
               执行 TaskComplete
```

## 四、关键设计点

1. **定时轮询**: 调度器每秒轮询一次，查找待执行任务
2. **异步执行**: 每个任务在独立的 goroutine 中执行
3. **状态驱动**: 通过状态机控制任务、阶段、步骤的执行流程
4. **用户干预**: 支持手动模式和用户操作接口，允许在特定节点进行人工干预
5. **错误处理**: 根据 Step 和 Stage 的模式决定失败后的行为（停止或继续）
6. **预处理**: 任务执行前进行清理和初始化工作

## 五、流程图

详见下面的 Mermaid 流程图。

