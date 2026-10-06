#!/usr/bin/env python3
"""Remove orphaned function bodies left by sed removal of function declarations."""
import os
import re

PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

def fix_file(filepath):
    with open(filepath) as f:
        lines = f.readlines()
    
    # Find orphaned code blocks (lines that start with tab and aren't inside a function)
    result = []
    skip_until_brace = False
    brace_count = 0
    in_orphan = False
    
    for i, line in enumerate(lines):
        stripped = line.strip()
        
        # Detect orphaned code (starts with tab but not inside a function)
        if stripped.startswith('if t != nil') or stripped.startswith('if s == nil') or stripped.startswith('if t == nil'):
            # Check if previous line is a function declaration
            if i > 0:
                prev = lines[i-1].strip()
                if not prev.startswith('func ') and not prev.startswith('//'):
                    # This is orphaned code - skip until closing brace
                    in_orphan = True
                    brace_count = 0
                    continue
        
        if in_orphan:
            if '{' in stripped:
                brace_count += stripped.count('{')
            if '}' in stripped:
                brace_count -= stripped.count('}')
            if brace_count <= 0 and stripped == '}':
                in_orphan = False
            continue
        
        result.append(line)
    
    if len(result) != len(lines):
        with open(filepath, 'w') as f:
            f.writelines(result)
        return True
    return False

def main():
    for root, dirs, files in os.walk(os.path.join(PROJECT_ROOT, 'service')):
        for fname in files:
            if fname == 'mapping.go':
                filepath = os.path.join(root, fname)
                if fix_file(filepath):
                    print(f"Fixed: {os.path.relpath(filepath, PROJECT_ROOT)}")

if __name__ == '__main__':
    main()
