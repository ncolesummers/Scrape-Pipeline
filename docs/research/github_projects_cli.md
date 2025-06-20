# GitHub Projects v2 CLI How-To Guide

This guide provides essential commands for working with GitHub Projects v2 via the GitHub CLI. The project number for Scrape-Pipeline is **2**.

## Setup

```bash
# Install GitHub CLI
brew install gh  # macOS
# Check installation
gh --version

# Authenticate with project scope
gh auth login
# Verify project scope is included
gh auth status
```

## Core Project Commands

```bash
# List your projects
gh project list

# Create a new project
gh project create --title "Scrape Pipeline Tasks" --owner your-username

# View project details
gh project view 2 --owner your-username

# Copy project structure
gh project copy 2 --source-owner source-org --title "New Sprint"
```

## Working with Items

```bash
# List items in a project
gh project item-list 2 --owner your-username

# Add issue to project
gh project item-add 2 --owner your-username --url https://github.com/user/repo/issues/10

# Create draft issue in project
gh project item-create 2 --owner your-username --title "Implement rate limiting" --body "Add configurable rate limiting to scraper"

# Edit item fields
gh project item-edit 2 --owner your-username --id ITEM_ID --field-id STATUS_FIELD_ID --text "In Progress"
```

## Automation Examples

```bash
# GitHub Action to auto-add closed issues to "Done"
gh workflow run update-project.yml -f issue_number=15 -f column_name="Done"

# Batch update items via script
for issue in $(gh issue list --state open --json number --jq '.[].number'); do
  gh project item-add 2 --owner your-username --url "https://github.com/your-username/repo/issues/$issue"
done
```

## Field Management

```bash
# List fields in project
gh project field-list 2 --owner your-username

# Create custom field
gh project field-create 2 --owner your-username --name "Priority" --data-type SingleSelect --options "High,Medium,Low"
```

## JSON Output for Scripting

For automations and scripting, you can use the `--json` flag with most commands to get structured data:

```bash
# Get project items in JSON format
gh project item-list 2 --owner your-username --json

# Use jq to filter results
gh project item-list 2 --owner your-username --json | jq '.[] | select(.status=="In Progress")'
```