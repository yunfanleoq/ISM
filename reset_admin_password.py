#!/usr/bin/env python3
"""
ISM 系统 Admin 密码重置工具
用法: python3 reset_admin_password.py [新密码]
如果不指定密码，默认重置为 "admin123"
"""

import sys
import sqlite3
import bcrypt
import hashlib
import os

# 数据库路径
DB_PATH = os.path.join(os.path.dirname(os.path.abspath(__file__)), 
                        "ism_server_user", "data", "db", "ism.db")

def reset_admin_password(new_password: str = "admin123"):
    """
    重置 admin 用户的密码。
    注意：前端会对密码做 MD5 后再发送到后端，所以数据库中存储的 bcrypt 值
    是对 MD5 哈希值的 bcrypt 结果。
    """
    if not os.path.exists(DB_PATH):
        print(f"❌ 数据库文件不存在: {DB_PATH}")
        print("请确认 ism_server_user/data/db/ism.db 文件路径是否正确")
        sys.exit(1)
    
    # 前端: md5(用户输入密码) -> 后端: bcrypt(md5值)
    md5_hash = hashlib.md5(new_password.encode()).hexdigest()
    bcrypt_hash = bcrypt.hashpw(md5_hash.encode(), bcrypt.gensalt()).decode()
    
    print(f"新密码: {new_password}")
    print(f"前端发送的 MD5: {md5_hash}")
    print(f"存储的 bcrypt: {bcrypt_hash}")
    
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()
    
    # 检查 admin 用户是否存在
    cursor.execute("SELECT id, username, role FROM user WHERE username = 'admin'")
    user = cursor.fetchone()
    
    if user is None:
        print("❌ 未找到 admin 用户！正在创建...")
        import uuid
        new_uuid = str(uuid.uuid4())
        cursor.execute(
            "INSERT INTO user (created_at, updated_at, deleted_at, username, password, name, phone, email, avatar, job, profile, role, uuid) "
            "VALUES (datetime('now'), datetime('now'), NULL, 'admin', ?, '超级管理员', '', '', '', '', '', 'Admin', ?)",
            (bcrypt_hash, new_uuid)
        )
        conn.commit()
        print(f"✅ admin 用户已创建，uuid: {new_uuid}")
    else:
        cursor.execute(
            "UPDATE user SET password = ? WHERE username = 'admin'",
            (bcrypt_hash,)
        )
        conn.commit()
        print(f"✅ admin 用户密码已重置 (id={user[0]}, role={user[2]})")
    
    # 验证
    cursor.execute("SELECT password FROM user WHERE username = 'admin'")
    stored_hash = cursor.fetchone()[0]
    if bcrypt.checkpw(md5_hash.encode(), stored_hash.encode()):
        print(f"✅ 密码验证通过！请使用用户名 'admin' 和密码 '{new_password}' 登录")
    else:
        print("⚠️  密码验证失败，请检查")
    
    conn.close()

if __name__ == "__main__":
    password = sys.argv[1] if len(sys.argv) > 1 else "admin123"
    reset_admin_password(password)
