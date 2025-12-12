#!/usr/bin/env python3
"""Script to resolve merge conflicts"""

from pathlib import Path

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

base_path = Path("/Users/supratouchsuwatno/Desktop/decm")
resolved_count = 0

for file_path in files_with_conflicts:
    full_path = base_path / file_path
    if full_path.exists():
        try:
            with open(full_path, "r", encoding="utf-8") as f:
                lines = f.readlines()

            # Process line by line
            new_lines = []
            i = 0
            while i < len(lines):
                line = lines[i]

                if line.startswith("<<<<<<< HEAD"):
                    # Found conflict start
                    # Skip lines until we find '======='
                    i += 1
                    while i < len(lines) and not lines[i].startswith("======="):
                        i += 1

                    # Skip '=======' line
                    if i < len(lines):
                        i += 1

                    # Keep lines until '>>>>>>>'
                    while i < len(lines) and not lines[i].startswith(">>>>>>>"):
                        new_lines.append(lines[i])
                        i += 1

                    # Skip '>>>>>>>' line
                    if i < len(lines):
                        i += 1
                else:
                    new_lines.append(line)
                    i += 1

            # Write back if conflicts were found
            if "<<<<<<< HEAD" in "".join(lines):
                with open(full_path, "w", encoding="utf-8") as f:
                    f.writelines(new_lines)
                print(f"✓ Resolved: {file_path}")
                resolved_count += 1
            else:
                print(f"- No conflicts: {file_path}")
        except Exception as e:
            print(f"✗ Error: {file_path}: {e}")
    else:
        print(f"✗ Not found: {file_path}")

print(f"\n{resolved_count} files resolved!")
