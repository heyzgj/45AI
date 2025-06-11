### 整体改造思路

> **只谈前端**：所有文案、色彩、字体、间距、动效、组件形态——一句话归纳：**统一语义 + 留白呼吸 + 轻量动画**。下面每条都是你可以直接写进 `uni.scss`、`.vue` 或组件库里的可落地规范。

---

## 1. Design Tokens（一次写好，全局调用）

| Token                | 取值                             | 用途              |
| -------------------- | ------------------------------ | --------------- |
| `--brand-pink`       | `#FF6B84`                      | 主按钮 / 选中态 / 进度条 |
| `--brand-pink-light` | `#FFECEE`                      | 点击波纹、禁用态背景      |
| `--bg-page`          | `#FFFFFF`                      | 页面背景            |
| `--bg-card`          | `#FFFFFF`                      | 卡片、弹窗           |
| `--fg-title`         | `#1F1F1F`                      | 主要文字            |
| `--fg-sub`           | `#666666`                      | 副标题、描述          |
| `--radius-md`        | `12rpx`                        | 卡片 & 按钮圆角       |
| `--shadow-elev3`     | `0 4rpx 12rpx rgba(0,0,0,.04)` | 卡片投影            |
| `--grid`             | `4rpx`                         | 水平 4 的倍数栅格      |

> **UniBest 接入**：在 `uni.scss` 写 `:root {}`，组件里用 `var(--brand-pink)`。

---

## 2. 字体 / 行高 / 间距

| 文档级别    | 字号(rpx) | 行高  | 场景             |
| ------- | ------- | --- | -------------- |
| H1      | 40      | 1.2 | 顶级标题（“AI 写真馆”） |
| H2      | 32      | 1.3 | 卡片标题 / 模版名     |
| Body    | 28      | 1.5 | 正文文案           |
| Caption | 24      | 1.4 | 描述、标签、时间       |

* **左右边距**：默认 `32rpx`。
* **卡片间距**：瀑布流 `16rpx`。
* **按钮高度**：`96rpx`，圆角 `48rpx`。

---

## 3. 首页 (模板商城)

### 3.1 Banner

* **删除前后对比需求** → 单张「品牌宣传」海报即可，尺寸 3:4，高度 `75vw`。
* 指示点居中，选中点用 `--brand-pink`。

### 3.2 功能宫格

```vue
<Grid4 rounded shadow-elev3 gap="24rpx">
  <GridItem icon="camera" text="生成"/>
  <GridItem icon="template" text="模板"/>
  <GridItem icon="portrait" text="写真馆"/>
  <GridItem icon="user" text="我的"/>
</Grid4>
```

* 图标全部换 **IconPark line** 24×24rpx，单色 `#FF6B84`。

### 3.3 Tab & Tag

```vue
<Tabs underline
      :active-bar-style="{background: 'var(--brand-pink)', height: '6rpx', borderRadius: '3rpx', transition: '300ms cubic-bezier(.4,0,.2,1)'}">
  <Tab title="推荐"/>
  <Tab title="最新"/>
  <Tab title="免费"/>
  <Tab title="收藏"/>
</Tabs>
```

* 下划线宽度随文字宽度自适应；切换时使用 *spring* 缓动。
* 二级标签改为 **横向可滚 Chips**，选中背景 `--brand-pink-light` + 字体 `--brand-pink`。

### 3.4 模板卡片

```vue
<Card shadow-elev3 radius="12rpx">
  <Image ratio="3:4" src="cover.jpg"/>
  <View class="p-24">
     <Text class="h2 single-line">Dusty Rose Dream</Text>
     <View class="flex-between mt-8">
        <Text class="brand">15 胶卷</Text>
     </View>
  </View>
</Card>
```

* 仅保留 **模版名 + 积分**，其他全部删。

---

## 4. 模板详情

| 元素       | 设计                                   |
| -------- | ------------------------------------ |
| 返回按钮     | 左上 `<` 24×24rpx，透明背景，点击区 44×44rpx    |
| 积分 + CTA | 横排<br>`Text(15 胶卷)` & `Button(立即生成)` |
| 描述       | 一行灰字 “柔雾粉调 · 梦幻灯效”                   |

---

## 5. AI 写真馆流程

### 5.1 上传页

```html
<Card dashed radius="12rpx" class="upload-box center">
  <Icon type="plus" size="48" color="var(--brand-pink)"/>
  <Text class="body">上传正面自拍</Text>
  <Caption>光线充足 · 无滤镜 · 露全脸</Caption>
</Card>
```

* 整个 Card 可点，长按弹出「拍照 / 相册」。

### 5.2 生成中 (Modal)

```vue
<Modal transparent blur>
  <Lottie file="loading.json" size="120rpx"/>
  <H2 class="mt-24">AI 正在绘制</H2>
  <Caption>预计 10-15 秒</Caption>
  <Progress :value="progress" height="8rpx" color="var(--brand-pink)"/>
</Modal>
```

### 5.3 生成结果

* **单图**：上方大图，底下三个按钮

  * 主色按钮 `保存到相册`
  * 线框按钮 `再生成一张`
  * 灰态按钮 `分享（即将上线）`
* **多图 (>10)**：进入 `PhotoPicker` 组件

  * 横向缩略条 + 主预览图
  * 底部浮动按钮 `保存选中 (X)` / `全选` / `再生成`

---

## 6. 「我的」+ 购买胶卷

### 6.1 Profile

* 头像上传后自动裁成 1:1 `radiusMd`。
* 昵称行尾放铅笔图标，可 inline 编辑。

### 6.2 最近活动

```vue
<ActivityItem icon="palette" title="使用 'Cyberpunk' 模版" sub="今日" point="-10"/>
```

* 三行即可：图标 / 模版名 / -积分。

### 6.3 交易记录空态

```html
<Illustration src="empty.svg"/>
<Body class="mt-24">暂无交易记录</Body>
<Button ghost brand class="mt-16">去生成美照</Button>
```

### 6.4 购买胶卷

* **卡片**：高 400rpx，信息层级

  ```
  120 胶卷
  ¥19.9          性价比首选
  均价 ¥0.16/张
  ```
* 选中卡片外框 `--brand-pink`; 未选灰框。
* CTA 灰态 `请选择套餐` → 选中后高亮 `立即购买`。

---

## 7. 动效小贴士

| 交互     | 动效描述                  | 实现                                                           |
| ------ | --------------------- | ------------------------------------------------------------ |
| Tab 切换 | 下划线滑动 + 8% 弹性回弹       | CSS `transition: transform .3s cubic-bezier(.34,1.56,.64,1)` |
| 卡片点击   | `scale(0.96)` 短压，松手回弹 | `:active` + `transition .1s`                                 |
| 生成完成   | Confetti 彩屑 600ms     | Lottie / canvas                                              |

---

## 8. 样本文案（全部中文化）

| 场景        | 文案                         |
| --------- | -------------------------- |
| 首页 Banner | “一键生成 10+ 专属美照”            |
| CTA(生成)   | “马上变美照”                    |
| 积分不足      | “胶卷不足，去充值？”                |
| 上传提醒      | “上传正面自拍\n光线充足 · 无滤镜 · 露全脸” |
| 保存成功      | “已保存相册，快去发朋友圈！”            |

---

## 9. 交付顺序（纯前端）

1. **token.scss**：写所有颜色/圆角/阴影。
2. **Tabs / Chips 组件**：下划线动画 + 选中态。
3. **Card 组件**：统一阴影/圆角/占位图。
4. **Upload / Loading / Result 组件**：按写真馆流程串起来。
5. **购买胶卷页**：卡片选中逻辑 + 价格细节。
6. 全站 **文案替换**(上表) & Emoji → 正规 IconPark。

照此执行，你的页面将从「开发者 Demo 感」→「可商用级」至少 **10×** 体面。剩下就看真实图片素材上线与否了。祝你改完即显高级感！
