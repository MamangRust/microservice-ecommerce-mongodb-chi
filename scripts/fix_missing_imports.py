#!/usr/bin/env python3
"""
Fix missing pb imports by detecting used aliases and adding imports.
"""
import os
import re

PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# Alias -> import path mapping (using actual underscore aliases from migration)
ALIAS_MAP = {
    'pbuser': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/user',
    'pbrole': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/role',
    'pbcategory': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/category',
    'pborder': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/order',
    'pbproduct': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/product',
    'pbmerchant': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant',
    'pbmerchant_document': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant_document',
    'pbmerchant_detail': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant_detail',
    'pbmerchant_business': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant_business',
    'pbmerchant_award': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant_award',
    'pbmerchant_policy': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant_policy',
    'pbmerchant_social_link': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant_social_link',
    'pbtransaction': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/transaction',
    'pborder_item': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/order_item',
    'pbcart': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/cart',
    'pbauth': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/auth',
    'pbbanner': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/banner',
    'pbslider': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/slider',
    'pbshipping_address': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/shipping_address',
    'pbreview': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/review',
    'pbreview_detail': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/review_detail',
    'pbcommon': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/common',
    'pbstats': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/stats',
    'pbcatestats': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/category',
    'pborderstats': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/order',
    'pbtxstats': 'github.com/MamangRust/microservice-ecommerce-grpc/pb/transaction',
}

def detect_used_aliases(content):
    used = set()
    for m in re.finditer(r'\b(pb[a-z_]+)\.[A-Z]', content):
        alias = m.group(1)
        if alias in ALIAS_MAP:
            used.add(alias)
    return used

def detect_existing_imports(content):
    imported = set()
    for alias in ALIAS_MAP:
        if f'{alias} "{ALIAS_MAP[alias]}"' in content:
            imported.add(alias)
    return imported

def fix_file(filepath):
    with open(filepath) as f:
        content = f.read()
    
    used = detect_used_aliases(content)
    existing = detect_existing_imports(content)
    missing = used - existing
    
    if not missing:
        return False
    
    import_lines = []
    for alias in sorted(missing):
        import_lines.append(f'\t{alias} "{ALIAS_MAP[alias]}"')
    
    if 'import (' in content:
        new_imports = '\n'.join(import_lines)
        content = content.replace('import (', 'import (\n' + new_imports, 1)
    
    with open(filepath, 'w') as f:
        f.write(content)
    return True

def main():
    modified = 0
    for root, dirs, files in os.walk(PROJECT_ROOT):
        dirs[:] = [d for d in dirs if d not in ('.git', 'vendor', 'bin', 'pb')]
        for fname in files:
            if not fname.endswith('.go'):
                continue
            filepath = os.path.join(root, fname)
            if fix_file(filepath):
                modified += 1
                print(f"Fixed: {os.path.relpath(filepath, PROJECT_ROOT)}")
    
    print(f"\nFixed {modified} files")

if __name__ == '__main__':
    main()
