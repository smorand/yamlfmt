#!/bin/bash
set -e

# Update main from upstream (origin)
echo "Fetching and updating main from origin..."
git fetch origin
git checkout main
git pull origin main

# Rebase smo branch on main
echo "Rebasing smo on main..."
git checkout smo
git rebase main

echo "Done. Run 'git push smo smo --force-with-lease' to push changes."
