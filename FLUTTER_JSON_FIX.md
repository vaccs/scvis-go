# Fix: Flutter JSON Format Error for Single Analysis

## Problem

When using the single analysis feature (`single.tmpl`), the Flutter visualization app was receiving HTML instead of JSON, resulting in the error:

```
[SCVis] Error loading JSON: FormatException: SyntaxError: Unexpected token '<', "<!doctype "... is not valid JSON
```

## Root Cause

The `/scvis/project/getanalysis` endpoint was configured as a `protected` route, requiring authentication. When the unauthenticated single analysis flow tried to fetch the JSON, it was redirected to the login page (HTML), which the Flutter app tried to parse as JSON.

## Solution

Made the following changes to support both authenticated project analysis and unauthenticated single analysis:

### 1. Updated Route Configuration (`routes.go`)

Changed `/scvis/project/getanalysis` from `protected` to `dynamic` middleware:

```go
// Before:
mux.Handle("GET /scvis/project/getanalysis", protected.ThenFunc(app.projectGetAnalysis))

// After:
mux.Handle("GET /scvis/project/getanalysis", dynamic.ThenFunc(app.projectGetAnalysis))
```

This allows the endpoint to be accessed without authentication.

### 2. Enhanced Handler (`handlers.go`)

Updated `projectGetAnalysis` to handle both cases:

**For Single Analysis (projectID = -1):**
- Reads JSON from temporary directory stored in session
- Returns JSON directly to Flutter app
- Cleans up temp files and session data after serving
- No authentication required

**For Project Analysis (projectID > 0):**
- Checks authentication before proceeding
- Redirects to login if not authenticated
- Fetches analysis from database
- Returns JSON from stored project

### 3. Added Cookie Support

Updated `singleAnalyzePost` to set the `projectID` cookie:

```go
http.SetCookie(w, &http.Cookie{
    Name:  "projectID",
    Value: "-1",
    Path:  "/scvis/",
})
```

This signals to the Flutter app and backend that this is a single analysis request.

### 4. Added Cleanup

The single analysis handler now properly cleans up:
- Removes session data after serving JSON
- Deletes temporary directory in background goroutine
- Prevents disk space accumulation from abandoned analyses

## Testing

After these changes:

1. Build and restart the application:
   ```bash
   go build ./cmd/web
   sudo systemctl restart scvisual
   ```

2. Test single analysis flow:
   - Go to `https://localhost:8082/scvis/`
   - Upload C files
   - Click "Analyze"
   - Flutter visualization should load without JSON errors

3. Test authenticated project analysis:
   - Login to the application
   - Create/open a project
   - Run analysis
   - Flutter visualization should still work as before

## Flow Diagram

### Single Analysis Flow
```
User uploads files → singleAnalyzePost
                   ↓
          Sets projectID=-1 cookie
          Stores tempDir in session
          Saves analysis.json
                   ↓
          Redirect to /scvis/visual
                   ↓
          Flutter app loads
                   ↓
          Fetches /scvis/project/getanalysis
                   ↓
          Sees projectID=-1 cookie
          Reads from temp directory
          Returns JSON
          Cleans up
```

### Authenticated Project Analysis Flow
```
User analyzes project → projectAnalyzePost
                      ↓
            Sets projectID=X cookie
            Saves to database
                      ↓
            Redirect to /scvis/visual
                      ↓
            Flutter app loads
                      ↓
            Fetches /scvis/project/getanalysis
                      ↓
            Checks authentication
            Reads from database
            Returns JSON
```

## Files Modified

1. `cmd/web/handlers.go`
   - Enhanced `projectGetAnalysis` with authentication check for non-single analysis
   - Updated `singleAnalyzePost` to set projectID cookie
   - Added background cleanup of temp directories

2. `cmd/web/routes.go`
   - Changed `projectGetAnalysis` from `protected` to `dynamic` middleware
   - Added `singleGetAnalysis` route (backup, not needed after main fix)

## Related Issues

This fix resolves:
- Flutter JSON parsing errors in single analysis
- Temp directory cleanup issues
- Authentication flow for mixed authenticated/unauthenticated endpoints
