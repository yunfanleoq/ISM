"""
ISM Automation — HTTP API 客户端
封装 ISM 后端所有 REST API，提供登录、Token 管理、重试、错误处理。

基于 httpx + tenacity，支持同步调用。
"""
from __future__ import annotations

import hashlib
import json
import time
from typing import Any, Dict, Optional, Union
from dataclasses import dataclass

import httpx
from tenacity import retry, stop_after_attempt, wait_exponential, retry_if_exception_type

from ism_automation.core.logger import ism_logger


# ── 异常定义 ──────────────────────────────────────

class ISMAPIError(Exception):
    """ISM API 调用失败"""
    def __init__(self, message: str, status_code: int = 0, response: Any = None):
        super().__init__(message)
        self.status_code = status_code
        self.response = response


class ISMAuthError(ISMAPIError):
    """认证失败（Token 过期、密码错误等）"""
    pass


class ISMValidationError(ISMAPIError):
    """请求参数校验失败"""
    pass


# ── 配置 ──────────────────────────────────────────

@dataclass
class ClientConfig:
    base_url: str = "http://127.0.0.1:8081"
    username: str = "admin"
    password: str = "123456"  # 明文密码，内部会 MD5 加密
    timeout: float = 30.0
    max_retries: int = 3


class ISMClient:
    """
    ISM 后端 API 客户端

    使用示例:
        client = ISMClient(config=ClientConfig(
            base_url="http://127.0.0.1:8081",
            username="admin",
            password="123456"
        ))
        client.login()
        models = client.model.list()
    """

    def __init__(self, config: Optional[ClientConfig] = None):
        self.config = config or ClientConfig()
        self._token: Optional[str] = None
        self._token_expires: float = 0
        self._client: Optional[httpx.Client] = None
        self._headers: Dict[str, str] = {}

    # ── 生命周期 ─────────────────────────────────────

    def __enter__(self):
        self._client = httpx.Client(
            base_url=self.config.base_url,
            timeout=httpx.Timeout(self.config.timeout),
            follow_redirects=True,
        )
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        if self._client:
            self._client.close()
            self._client = None

    def _ensure_client(self):
        if self._client is None:
            self._client = httpx.Client(
                base_url=self.config.base_url,
                timeout=httpx.Timeout(self.config.timeout),
                follow_redirects=True,
            )

    # ── 认证 ─────────────────────────────────────────

    def _md5_password(self, password: str) -> str:
        """ISM 前端将密码 MD5 后提交"""
        return hashlib.md5(password.encode("utf-8")).hexdigest()

    @retry(
        stop=stop_after_attempt(3),
        wait=wait_exponential(multiplier=1, min=1, max=10),
        retry=retry_if_exception_type((httpx.NetworkError, httpx.TimeoutException)),
    )
    def login(self, username: Optional[str] = None, password: Optional[str] = None) -> str:
        """
        登录 ISM 系统，获取 Token

        ISM 登录流程:
        1. POST /login 提交 username + md5(password)
        2. 后端返回 token
        3. 后续请求 Header 携带 token
        """
        self._ensure_client()
        u = username or self.config.username
        p = password or self.config.password
        md5_p = self._md5_password(p)

        ism_logger.info(f"🔐 登录 ISM: {u}@{self.config.base_url}")

        resp = self._client.post(
            "/login",
            data={"username": u, "password": md5_p},
            headers={"Content-Type": "application/x-www-form-urlencoded"},
        )
        resp.raise_for_status()

        data = resp.json()
        # ISM 返回格式: { "code": 200, "message": "...", "data": { "token": "..." } }
        if data.get("code") != 200:
            raise ISMAuthError(f"登录失败: {data.get('message', '未知错误')}", response=data)

        self._token = data.get("data", {}).get("token")
        if not self._token:
            raise ISMAuthError("登录响应中未找到 token", response=data)

        self._token_expires = time.time() + 3600 * 24  # 默认 24h
        self._headers = {"token": self._token}

        ism_logger.info("✅ 登录成功")
        return self._token

    def _ensure_auth(self):
        """确保已登录，Token 过期时自动重新登录"""
        if not self._token or time.time() > self._token_expires:
            self.login()

    # ── 基础请求方法 ─────────────────────────────────

    @retry(
        stop=stop_after_attempt(3),
        wait=wait_exponential(multiplier=1, min=1, max=10),
        retry=retry_if_exception_type((httpx.NetworkError, httpx.TimeoutException)),
    )
    def _request(
        self,
        method: str,
        path: str,
        json_data: Optional[Dict] = None,
        form_data: Optional[Dict] = None,
        params: Optional[Dict] = None,
        files: Optional[Dict] = None,
        extra_headers: Optional[Dict] = None,
    ) -> Dict[str, Any]:
        """发送 HTTP 请求，处理认证和错误"""
        self._ensure_auth()
        self._ensure_client()

        headers = dict(self._headers)
        if extra_headers:
            headers.update(extra_headers)

        if json_data is not None:
            headers.setdefault("Content-Type", "application/json")

        ism_logger.debug(f"📡 {method} {path}")

        try:
            resp = self._client.request(
                method=method,
                url=path,
                json=json_data,
                data=form_data,
                params=params,
                files=files,
                headers=headers,
            )
            resp.raise_for_status()
        except httpx.HTTPStatusError as e:
            if e.response.status_code == 401:
                # Token 过期，重新登录后重试
                ism_logger.warning("Token 过期，重新登录...")
                self.login()
                return self._request(method, path, json_data, form_data, params, files, extra_headers)
            raise ISMAPIError(
                f"HTTP {e.response.status_code}: {e.response.text}",
                status_code=e.response.status_code,
                response=e.response.text,
            )

        # 解析响应
        try:
            data = resp.json()
        except json.JSONDecodeError:
            # 有些接口返回纯文本
            return {"_raw": resp.text, "code": 200}

        # 检查 ISM 业务错误码
        code = data.get("code", 200)
        if code != 200:
            msg = data.get("message", data.get("msg", f"业务错误码: {code}"))
            if code == 401:
                raise ISMAuthError(msg, status_code=code, response=data)
            raise ISMAPIError(msg, status_code=code, response=data)

        return data

    def get(self, path: str, params: Optional[Dict] = None) -> Dict[str, Any]:
        return self._request("GET", path, params=params)

    def post(self, path: str, json_data: Optional[Dict] = None, form_data: Optional[Dict] = None) -> Dict[str, Any]:
        return self._request("POST", path, json_data=json_data, form_data=form_data)

    def put(self, path: str, json_data: Optional[Dict] = None) -> Dict[str, Any]:
        return self._request("PUT", path, json_data=json_data)

    def delete(self, path: str, params: Optional[Dict] = None) -> Dict[str, Any]:
        return self._request("DELETE", path, params=params)

    # ── 快捷 API 分组（懒加载）───────────────────────

    @property
    def model(self):
        """数据模型 API"""
        from ism_automation.api.model_api import ModelAPI
        return ModelAPI(self)

    @property
    def device(self):
        """设备管理 API"""
        from ism_automation.api.device_api import DeviceAPI
        return DeviceAPI(self)

    @property
    def display(self):
        """组态大屏 API"""
        from ism_automation.api.display_api import DisplayAPI
        return DisplayAPI(self)

    @property
    def project(self):
        """项目管理 API"""
        from ism_automation.api.project_api import ProjectAPI
        return ProjectAPI(self)


# ── 便捷函数 ──────────────────────────────────────

def get_client(**kwargs) -> ISMClient:
    """快速获取已登录的客户端实例"""
    config = ClientConfig(**kwargs)
    client = ISMClient(config)
    client.login()
    return client
