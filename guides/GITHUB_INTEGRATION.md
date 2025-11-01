# GitHub Integration Setup for SCVis

This guide explains how to set up GitHub integration for SCVis to allow administrators to submit issues directly to the GitHub repository.

## Overview

The GitHub integration feature allows administrators to:
- Submit issues directly to the GitHub repository from the SCVis web interface
- Track issue status locally and link to GitHub issues
- View and manage all submitted issues
- Automatically format issues with relevant context

## Prerequisites

1. **GitHub Repository**: You need a GitHub repository to submit issues to
2. **GitHub Personal Access Token**: Required for API authentication
3. **Repository Permissions**: The token must have permissions to create issues in the target repository

## Setup Instructions

### 1. Create a GitHub Personal Access Token

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Give your token a descriptive name (e.g., "SCVis Issue Submission")
4. Select the following scopes:
   - `repo` (Full control of private repositories) - if your repository is private
   - `public_repo` (Access public repositories) - if your repository is public
5. Click "Generate token"
6. **Copy the token immediately** - you won't be able to see it again

### 2. Configure Environment Variables

Add the following variables to your `.env` file:

```bash
# GitHub Integration
GITHUB_TOKEN=ghp_your_personal_access_token_here
GITHUB_OWNER=vaccs
GITHUB_REPO=scvis-go
```

**Parameters:**
- `GITHUB_TOKEN`: Your GitHub personal access token
- `GITHUB_OWNER`: The GitHub username or organization that owns the repository
- `GITHUB_REPO`: The repository name

### 3. Example Configuration

For the official SCVis repository:

```bash
# Development (.env)
GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
GITHUB_OWNER=vaccs
GITHUB_REPO=scvis-go
```

For production (`/opt/scvisual/.env`):

```bash
# Production
GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
GITHUB_OWNER=vaccs
GITHUB_REPO=scvis-go
```

### 4. Security Considerations

1. **Token Security**: 
   - Never commit your GitHub token to version control
   - Use environment variables or secure configuration files
   - Regularly rotate your tokens
   - Use tokens with minimal required permissions

2. **File Permissions**:
   ```bash
   # Set appropriate permissions for configuration files
   chmod 600 .env
   sudo chmod 640 /opt/scvisual/.env
   sudo chown root:www-data /opt/scvisual/.env
   ```

3. **Token Scope**: Only grant the minimum permissions needed (issue creation)

## Usage

### For Administrators

1. **Navigate to Issues**: Click "Issues" in the main navigation
2. **Create New Issue**: Click "Submit New Issue" 
3. **Fill in Details**:
   - Title: Clear, descriptive title
   - Type: Bug, Feature, Enhancement, Documentation, or Question
   - Priority: Low, Medium, High, or Critical
   - Description: Detailed description of the issue
4. **GitHub Submission**: Check "Submit to GitHub repository" to create the issue on GitHub
5. **Submit**: Click "Submit Issue"

### Issue Management

- **View All Issues**: See all issues submitted through the system
- **Track Status**: Monitor issue progress from pending to resolved
- **GitHub Links**: Direct links to GitHub issues when submitted
- **Status Updates**: Administrators can update issue status

## GitHub Issue Format

When submitted to GitHub, issues are automatically formatted with:

```markdown
**Issue submitted from SCVis application**

**Type:** Bug Report
**Priority:** High
**Submitted by User ID:** 123

**Description:**
[User's description here]

---
*This issue was automatically created from the SCVis web application.*
```

Issues are automatically labeled with:
- `scvis-issue` - Identifies issues from SCVis
- Issue type (e.g., `bug`, `feature`, `enhancement`)
- Priority level (e.g., `priority-high`, `priority-critical`)

## Troubleshooting

### Common Issues

1. **"GitHub token not configured"**
   - Check that `GITHUB_TOKEN` is set in your environment
   - Verify the token is not expired
   - Ensure the token has correct permissions

2. **"GitHub API error (status 401)"**
   - Token is invalid or expired
   - Generate a new token with proper permissions

3. **"GitHub API error (status 403)"**
   - Token lacks sufficient permissions
   - Repository may not exist or be accessible

4. **"GitHub API error (status 404)"**
   - Repository not found
   - Check `GITHUB_OWNER` and `GITHUB_REPO` settings
   - Verify token has access to the repository

### Testing the Connection

The application will indicate if GitHub is properly configured:
- ✅ Green indicator: "GitHub Connected" 
- ⚠️ Warning: "GitHub integration is not configured"

### Verification Steps

1. **Check Configuration**:
   ```bash
   # Verify environment variables are loaded
   echo $GITHUB_TOKEN
   echo $GITHUB_OWNER
   echo $GITHUB_REPO
   ```

2. **Test API Access**:
   ```bash
   # Test GitHub API access
   curl -H "Authorization: token $GITHUB_TOKEN" \
        -H "Accept: application/vnd.github.v3+json" \
        https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO
   ```

3. **Check Logs**: Monitor application logs for GitHub API errors

## Development

For development purposes, you can:

1. **Use a Test Repository**: Create a test repository for development
2. **Local Testing**: Test against a personal repository before production
3. **Mock Mode**: The system gracefully handles missing GitHub configuration

### Example Development Setup

```bash
# Development environment
GITHUB_TOKEN=ghp_development_token_here
GITHUB_OWNER=your-username
GITHUB_REPO=scvis-test
```

## Production Deployment

1. **Generate Production Token**: Create a token specifically for production
2. **Set Production Configuration**: Update `/opt/scvisual/.env`
3. **Test Integration**: Submit a test issue to verify functionality
4. **Monitor Usage**: Regular check logs for any API issues

## API Rate Limits

GitHub API has rate limits:
- **Authenticated requests**: 5,000 per hour
- **Issue creation**: Should be well within limits for normal usage

The application includes appropriate error handling for rate limit scenarios.

## Support

If you encounter issues with GitHub integration:

1. Check the troubleshooting section above
2. Verify your GitHub token and permissions
3. Review application logs for specific error messages
4. Submit an issue through the SCVis interface (if GitHub integration is working)
5. Contact the development team

## Security Notes

- **Never expose tokens**: Ensure tokens are not logged or displayed in error messages
- **Regular rotation**: Change tokens periodically for security
- **Monitor usage**: GitHub provides token usage analytics
- **Audit permissions**: Regularly review what repositories your tokens can access
