package dto

type SubjectResponse struct {
    ID          uint64 `json:"id"`
    SubjectType string `json:"subjectType"`
    LoginID     string `json:"loginId"`
    ExternalID  string `json:"externalId"`
    Status      int    `json:"status"`
}

type ListSubjectsRequest struct {
    Page        int    `form:"page"`
    PageSize    int    `form:"pageSize"`
    SubjectType string `form:"subjectType"`
    Keyword     string `form:"keyword"`
}


