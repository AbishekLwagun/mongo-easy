import os
import sys
import subprocess

# Get the folder where this setup script is located.
install_dir = os.path.dirname(os.path.abspath(__file__))

# Check that the Mongo Easy executable exists.
exe_path = os.path.join(install_dir, "mongo-easy.exe")

if not os.path.exists(exe_path):
    print("Error: mongo-easy.exe was not found.")
    sys.exit(1)

print()
print("Mongo Easy Setup")
print()
print("Mongo Easy was found at:")
print(install_dir)
print()

answer = input("Add Mongo Easy to your user PATH? [Y/n]: ").strip().lower()

# Stop if the user does not want to change PATH.
if answer not in ("", "y", "yes"):
    print()
    print("Setup cancelled. Nothing was changed.")
    sys.exit(0)

# Get the current user PATH from Windows.
result = subprocess.run(
    [
        "powershell",
        "-NoProfile",
        "-Command",
        "[Environment]::GetEnvironmentVariable('Path', 'User')",
    ],
    capture_output=True,
    text=True,
)

if result.returncode != 0:
    print("Error: Could not read the user PATH.")
    sys.exit(1)

current_path = result.stdout.strip()

# Check if Mongo Easy is already in the user PATH.
path_entries = current_path.split(os.pathsep) if current_path else []

if install_dir in path_entries:
    print()
    print("Mongo Easy is already in your user PATH.")
    sys.exit(0)

# Add Mongo Easy to the user PATH.
new_path = current_path + os.pathsep + install_dir if current_path else install_dir

result = subprocess.run(
    [
        "powershell",
        "-NoProfile",
        "-Command",
        "[Environment]::SetEnvironmentVariable('Path', $env:NEW_PATH, 'User')",
    ],
    env={**os.environ, "NEW_PATH": new_path},
)

if result.returncode != 0:
    print()
    print("Error: Could not update the user PATH.")
    sys.exit(1)

print()
print("✔ Mongo Easy was added to your user PATH.")
print()
print("Close this PowerShell window and open a new one.")
print("Then run:")
print()
print("  mongo-easy --help")
print()