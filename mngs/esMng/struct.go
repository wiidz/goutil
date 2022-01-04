package esMng

type SearchType string

// 这几种好像不支持
const MatchPhrase = "match_phrase" // 精确匹配所有同时包含
const MultiMatch = "multi_match"   // 两个字段进行匹配，其中一个字段有这个文档就满足
const Term = "term"                // 完全匹配，即不进行分词器分析，文档中必须包含整个搜索的词汇

// 这几种支持
const BestFields = "best_fields"
const Boolean = "boolean"
const MostFields = "most_fields" // 越多字段匹配的文档评分越高
const CrossFields = "cross_fields" // 分词词汇是分配到不同字段中
const Phrase = "phrase"
const PhrasePrefix = "phrase_prefix"
const BoolPrefix = "bool_prefix"