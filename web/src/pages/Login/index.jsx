import React, { useState, useMemo, useEffect, useContext } from "react";
import { Form, Input, Button, Alert, message } from "antd";
import { UserOutlined, LockOutlined } from "@ant-design/icons";

import { AuthContext } from "shared/contexts";
import { navigate } from "@reach/router";
import "./styles.scss";

import { useAPI } from "shared/hooks";

const decodeJWT = (token) => {
  try {
    return JSON.parse(atob(token.split(".")[1]));
  } catch (e) {
    return undefined;
  }
};

const logoUrl = "/logo.png";

function useLogin() {
  const [body, setBody] = useState(undefined);
  const args = useMemo(() => (body ? ["/api/v1/auth/login", { method: "POST", body: JSON.stringify(body) }] : [undefined, undefined]), [
    body,
  ]);

  const [data, status] = useAPI(...args);
  let user;
  if (status.isSuccess) {
    const token = data?.token;
    user = decodeJWT(token);
    user.token = token;
  }

  return [user, status, setBody];
}

function Login() {
  const [u, setUser] = useContext(AuthContext);
  const [user, status, login] = useLogin();

  const handleFormSubmit = (b) => {
    login(b);
  };

  useEffect(() => {
    if (status.isSuccess) {
      setUser({
        ...user,
        avatar: user.name
          .toUpperCase()
          .split(" ")
          .slice(0, 2)
          .map((w) => w[0])
          .join(""),
        token: user.token,
      });
    }
  }, [status, setUser, user]);

  useEffect(() => {
    if (u?.token) navigate("/");
  }, [u]);

  useEffect(() => {
    if (status.isError && status.statusCode !== 401) {
      message.error("Oops! Something went wrong.", 3);
    }
  }, [status]);

  return (
    <div className="login-page">
      <Form className="login-form white-bg" initialValues={{ remember: true }} onFinish={handleFormSubmit} autoComplete="off">
        <div className="logo-wrapper" style={{ backgroundImage: `url(${logoUrl})` }}></div>
        <Form.Item name="username" rules={[{ required: true, message: "Input your username" }]}>
          <Input prefix={<UserOutlined className="site-form-item-icon" />} placeholder="Username" size="large" autoFocus />
        </Form.Item>
        <Form.Item name="password" rules={[{ required: true, message: "Input your password" }]}>
          <Input prefix={<LockOutlined className="site-form-item-icon" />} type="password" placeholder="Password" size="large" />
        </Form.Item>
        <Alert
          style={{
            visibility: status.isError && status.statusCode === 401 ? "visible" : "hidden",
            marginBottom: 16,
          }}
          message="Invalid username or password"
          type="error"
        />

        <Form.Item style={{ textAlign: "right" }}>
          <Button type="primary" htmlType="submit" loading={status.isLoading} size="large" block>
            Log in
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
}

export default Login;
