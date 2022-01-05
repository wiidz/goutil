package esMng

type SearchType string

// 这几种好像不支持
const MatchPhrase = "match_phrase" // 精确匹配所有同时包含
const MultiMatch = "multi_match"   // 两个字段进行匹配，其中一个字段有这个文档就满足
const Term = "term"                // 完全匹配，即不进行分词器分析，文档中必须包含整个搜索的词汇

// 这几种支持
const BestFields = "best_fields" // 【最佳字段】某一个field匹配尽可能多的关键词的doc优先返回来
const Boolean = "boolean"
const MostFields = "most_fields" // 【多数字段】尽可能返回更多field匹配某个关键词的doc，优先返回回来
const CrossFields = "cross_fields" // 【跨字段】分词词汇是分配到不同字段中
const Phrase = "phrase" // 必须包含一模一样的串，才会返回（包含短语的意思）
const PhrasePrefix = "phrase_prefix"
const BoolPrefix = "bool_prefix"