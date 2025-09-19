# GitHub Activity CLI  

A simple **command-line tool** written in Go that fetches and displays a user’s **recent GitHub activity** using the [GitHub Events API](https://docs.github.com/en/rest/activity/events).  

Url : [GitHub User Activity project on roadmap.sh](https://roadmap.sh/projects/github-user-activity).

## Features  
- Fetches the **latest public activity** of any GitHub user.  
- Filters events from the **last 48 hours**.  
- Displays detailed messages for different GitHub events:  
  - ✅ Issues (opened, closed, edited, etc.)  
  - ✅ Pull requests (opened, merged, reviewed, etc.)  
  - ✅ Stars ⭐  
  - ✅ Forks 🍴  
  - ✅ Pushes (commit count)  
  - ✅ Wiki edits, repo creation, deletion, and more.  
- Lightweight and easy to install as a global CLI tool.  

## Installation  
1. Clone this repository:  
   ```bash
   git clone https://github.com/<your-username>/github-activity.git
   cd github-activity
   ```
2. Build the binary:  
   ```bash
   go build -o github-activity
   ```
3. Move it into your \$PATH to use globally:  
   ```bash
   sudo mv github-activity /usr/local/bin/
   ```

## Usage  
```bash
github-activity <username>
```

### Example  
```bash
github-activity torvalds
```

