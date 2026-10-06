#!/usr/bin/env python3
"""
Migrate Go imports from shared/pb to per-service pb/<service> packages.

For each Go file that imports shared/pb:
1. Find all pb.XxxType references (types AND functions)
2. Map each to its new package
3. Replace pb.Xxx with the correct package alias
4. Add the required imports
"""

import os
import re
import sys
from collections import defaultdict

PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
PB_DIR = os.path.join(PROJECT_ROOT, "pb")

# Step 1: Build type-to-package mapping from pb/ directory
def build_type_map():
    type_map = {}
    for svc_dir in sorted(os.listdir(PB_DIR)):
        svc_path = os.path.join(PB_DIR, svc_dir)
        if not os.path.isdir(svc_path) or svc_dir in ('go.mod', 'go.sum', '.git'):
            continue
        for fname in os.listdir(svc_path):
            if not fname.endswith('.go'):
                continue
            fpath = os.path.join(svc_path, fname)
            with open(fpath) as f:
                for line in f:
                    # Match exported types: type XxxName
                    m = re.match(r'^type ([A-Z][a-zA-Z0-9]+)', line)
                    if m:
                        type_map[m.group(1)] = svc_dir
                    # Match exported functions: func NewXxxClient, func RegisterXxxServer, func UnimplementedXxx
                    m = re.match(r'^func (New[A-Z][a-zA-Z0-9]+|Register[A-Z][a-zA-Z0-9]+|Unimplemented[A-Z][a-zA-Z0-9]+)', line)
                    if m:
                        type_map[m.group(1)] = svc_dir
    return type_map

# Step 2: For each Go file, find pb.Xxx references
def find_pb_refs(filepath):
    refs = set()
    with open(filepath) as f:
        content = f.read()
    # Match pb.TypeName patterns
    for m in re.finditer(r'\bpb\.([A-Z][a-zA-Z0-9]+)', content):
        refs.add(m.group(1))
    return refs

# Step 3: Build import aliases
def pkg_alias(pkg):
    """Create a consistent alias for each package"""
    return f"pb{pkg}"

def import_path(pkg):
    return f"github.com/MamangRust/microservice-ecommerce-grpc/pb/{pkg}"

# Step 4: Process each file
def process_file(filepath, type_map, dry_run=False):
    with open(filepath) as f:
        content = f.read()
    
    # Skip if no shared/pb import
    if 'microservice-ecommerce-shared/pb"' not in content:
        return False
    
    # Find all pb.Xxx references
    refs = find_pb_refs(filepath)
    if not refs:
        return False
    
    # Map types to packages
    pkg_types = defaultdict(set)
    unknown_types = set()
    for type_name in refs:
        if type_name in type_map:
            pkg_types[type_map[type_name]].add(type_name)
        else:
            unknown_types.add(type_name)
    
    if unknown_types:
        print(f"  WARNING: {os.path.relpath(filepath, PROJECT_ROOT)}: unknown types: {unknown_types}")
    
    if not pkg_types:
        return False
    
    # Build the replacement: pb.Xxx -> alias.Xxx
    new_content = content
    
    # Replace all pb.Xxx with the correct alias
    for pkg, types in pkg_types.items():
        alias = pkg_alias(pkg)
        for type_name in types:
            new_content = re.sub(r'\bpb\.' + re.escape(type_name) + r'\b', 
                               f'{alias}.{type_name}', new_content)
    
    # Replace the old import
    old_import = 'pb "github.com/MamangRust/microservice-ecommerce-shared/pb"'
    new_imports = []
    for pkg in sorted(pkg_types.keys()):
        alias = pkg_alias(pkg)
        new_imports.append(f'\t{alias} "{import_path(pkg)}"')
    
    new_import_block = '\n'.join(new_imports)
    new_content = new_content.replace(old_import, new_import_block)
    
    # Write back if changed
    if new_content != content:
        if not dry_run:
            with open(filepath, 'w') as f:
                f.write(new_content)
        return True
    return False

def main():
    dry_run = '--dry-run' in sys.argv
    type_map = build_type_map()
    print(f"Loaded {len(type_map)} type mappings")
    
    # Find all Go files with shared/pb import
    modified = 0
    for root, dirs, files in os.walk(PROJECT_ROOT):
        # Skip pb/, .git/, vendor/
        dirs[:] = [d for d in dirs if d not in ('.git', 'vendor', 'bin', 'pb')]
        for fname in files:
            if not fname.endswith('.go'):
                continue
            filepath = os.path.join(root, fname)
            if process_file(filepath, type_map, dry_run):
                modified += 1
                print(f"  Modified: {os.path.relpath(filepath, PROJECT_ROOT)}")
    
    print(f"\nModified {modified} files")

if __name__ == '__main__':
    main()
