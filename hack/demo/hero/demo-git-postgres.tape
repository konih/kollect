# Hero demo — Git + Postgres (~2–3 min). Run: task demo-hero-record-git-postgres
Output docs/assets/demo/hero-git-postgres.mp4
Output docs/assets/demo/hero-git-postgres.gif

Set Shell "bash"
Set FontSize 12
Set Width 1200
Set Height 680
Set Theme "Catppuccin Mocha"
Set TypingSpeed 40ms
Set PlaybackSpeed 1.0

# Title
Type "# kollect — Git audit trail + Postgres queries"
Sleep 2s
Ctrl+L

# Bootstrap (extended variant)
Type "task demo-hero-up-postgres"
Enter
Sleep 2s
Hide
Sleep 10s
Show
Sleep 1s
Ctrl+L

# Show dual-sink inventory
Type "kubectl get kinv,ksnap,kdb -n default"
Enter
Sleep 5s
Ctrl+L

# Git money shot (same drift story)
Type "kubectl set image deployment/shop api=shop:v2.1.0 -n demo"
Enter
Sleep 3s
Hide
Sleep 25s
Show
Type "git -C /tmp/kollect-hero-inventory pull && git -C /tmp/kollect-hero-inventory log -p -1"
Enter
Sleep 10s
Ctrl+L

# Postgres tangible moment
Type "kubectl exec -n kollect-system deploy/postgres -- psql -U kollect -d inventory -c \"SELECT payload->'attributes'->>'image' AS image, COUNT(*) FROM public.inventory_items GROUP BY 1 ORDER BY 1;\""
Enter
Sleep 12s
Ctrl+L

# Parallel multi-sink hint (honest: both refs on one inventory)
Type "kubectl get kollectinventory demo-inventory -o yaml | rg -n 'snapshotSinkRefs|databaseSinkRefs' -A1"
Enter
Sleep 8s
Ctrl+L

Type "# Git for drift · Postgres for dashboards · same inventory"
Sleep 3s
