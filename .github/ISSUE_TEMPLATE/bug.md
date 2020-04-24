---
name: Bug
about: For when you encounter an unexpected behavior.
---

<!-- Please keep this line, it helps us automate tagging issues. [issue-type:bug-report] -->

### Terraform Version

<!-- Please paste the output of `terraform -v` here -->

### Effected Resource(s)

<!-- Please list the effected resources and data sources here. -->

- linode_XXXXXX

### Terraform Configuration Files

<!--
Please put your Terraform configuration files in terraform code blocks like the one below.

If your configuration is too complex to be put into code blocks, consider using a Github Gist or a service like Dropbox and share a link to a zip file containing your configuration.
-->

```terraform
# ...
```

### Debug Output

<!--
Please provide a link to a Github Gist containing the complete debug output. Please don't paste the output in the issue.

To obtain the debug output, run `terraform apply` with the environment variable `TF_LOG=DEBUG`, and optionally, `LINODE_DEBUG=1`. Please ensure you do not upload sensitive information.
-->

### Panic Output

<!--
If Terraform produced a panic, please provide a link to the Github Gist containing the output of the `crash.log`.
-->

### Expected Behavior

<!-- What should have happened? -->

### Actual Behavior

<!-- What actually happened? -->

### Steps to Reproduce

<!--
Please list the steps required to repoduce this issue.

Be sure to mention anything atypical about your accounts/setup.
-->
