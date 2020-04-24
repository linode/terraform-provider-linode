behavior "regexp_issue_labeler" "bug_label" {
    regexp = " \\[issue-type:bug-report\\] "
    labels = ["bug"]
}

behavior "regexp_issue_labeler" "enhancement_label" {
    regexp = " \\[issue-type:enhancement\\] "
    labels = ["enhancement"]
}
