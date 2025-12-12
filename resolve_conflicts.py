#!/usr/bin/env python3
"""Script to resolve merge conflicts by keeping the 'fix: slow build' version (removing extra blank lines)"""

import re
from pathlib import Path

# List of files with conflicts
files_with_conflicts = [
    "apps/backend/common/validatorutils/validatorutils_test.go",
    "apps/backend/core-api/internal/handler/event/generate_certificate_image_example.go.example",
    "apps/backend/core-api/internal/handler/eventconfig/check_certificate_mint_readiness_response.go",
    "apps/backend/core-api/internal/repositories/postgres/event_certificate_font_family.go",
    "apps/backend/core-api/internal/usecase/cyptoutils/ethereum_test.go",
    "apps/backend/core-api/internal/usecase/cyptoutils/rpc_test.go",
    "apps/backend/core-api/internal/usecase/cyptoutils/utils_test.go",
    "apps/backend/core-api/internal/usecase/eventconfig/check_certificate_mint_readiness.go",
    "apps/backend/core-api/internal/usecase/eventconfig/event_certificate_font_family.go",
    "apps/backend/core-api/internal/usecase/onboard/register_sign_message_test.go",
    "apps/backend/services/auth/auth_roles_test.go",
    "apps/backend/services/auth/auth_test.go",
    "apps/backend/services/auth/jwt_roles_test.go",
    "apps/backend/services/auth/jwt_test.go",
    "apps/backend/services/oauth/oauth_test.go",
    "packages/database/migrations/000009_add_font_config_to_certificate_configs.down.sql",
    "packages/database/migrations/000012_make_event_name_positions_nullable.up.sql",
]


def resolve_conflict(content):
    """Remove conflict markers, keeping the 'theirs' version (after =======)"""
    # Pattern to match conflict blocks
    # <<<<<<< HEAD\n...\n=======\n...\n>>>>>>> commit_hash
    pattern = r"<<<<<<< HEAD\n(.*?)\n=======\n(.*?)\n>>>>>>> [^\n]+\n?"

    def replace_conflict(match):
        # Keep the second part (after =======) which is the incoming change
        return match.group(2) if match.group(2) else ""

    resolved = re.sub(pattern, replace_conflict, content, flags=re.DOTALL)
    return resolved


# Process each file
base_path = Path("/Users/supratouchsuwatno/Desktop/decm")
resolved_count = 0

for file_path in files_with_conflicts:
    full_path = base_path / file_path
    if full_path.exists():
        try:
            # Read the file
            with open(full_path, "r", encoding="utf-8") as f:
                content = f.read()

            # Check if there are conflicts
            if "<<<<<<< HEAD" in content:
                # Resolve conflicts
                resolved_content = resolve_conflict(content)

                # Write back
                with open(full_path, "w", encoding="utf-8") as f:
                    f.write(resolved_content)

                print(f"✓ Resolved: {file_path}")
                resolved_count += 1
            else:
                print(f"- No conflicts: {file_path}")
        except Exception as e:
            print(f"✗ Error processing {file_path}: {e}")
    else:
        print(f"✗ File not found: {file_path}")

print(f"\n{resolved_count} files resolved successfully!")
