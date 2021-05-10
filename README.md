# TranslaX

一个轻量的命令行聚合翻译工具，支持谷歌、必应、搜狗、有道词典。

## 基本使用

1. 可选1（非交互）
运行时带参数，默认选择非交互式
```
translax [--from=zh] [--to=en]  [--text=我是翻译小能手]
```

2. 可选2（交互）
 
```
translax
```

## 高级
计划增加对文档（文本文档、PDF, docx...）翻译的支持

## 翻译引擎

| 引擎名称 | 说明 | 标准接口 | 国内直连 | 需要 Key 吗 | 免费 | 状态 |
|:---------|:--------|:--:|:--:|:--:|:--:|:--:|
|google|谷歌|Yes|Yes|No|Yes|可用|
|bing|必应|Yes|Yes|No|Yes|可用|
|sougou|搜狗|Yes|Yes|No|Yes|可用|
|youdao|有道|Yes|Yes|No|Yes|可用|

## LICENSE

