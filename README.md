# ghlabel
Since all company projects must abide by the [Drud community guidelines]
(https://github.com/drud/community/blob/master/development/issue_workflow.md#labels),
it makes sense to automate some processes. Hence, Ghlabel.

Ghlabel is a tool that automatically standardizes GitHub issue labels across a user or organization's repositories.
A reference repository is used as the template for labels, and those labels are automatically copied to all or
a single repository.

## Usage
The tool currently has two functions previewing staged label changes and applying them.

Ghlabel runs in preview mode by default.:
```
./ghlabel --org=drud --ref=community
```
You can execute changes using -r.
```
./ghlabel --org=drud --ref=community -r
```