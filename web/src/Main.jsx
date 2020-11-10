import React from "react";
import { Layout, Menu } from "antd";
import { Link } from "@reach/router";
import "./main.scss";

const { Content, Sider } = Layout;

function Main({ children, logout }) {
  return (
    <Layout className="site-layout">
      <Sider width={200}>
        <div className="logo" />
        <Menu theme="dark" mode="inline">
          <Menu.Item key="customers">
            <Link to="/customers">Customers</Link>
          </Menu.Item>
          <Menu.Item key="logout" onClick={logout}>
            Logout
          </Menu.Item>
        </Menu>
      </Sider>
      <Layout>
        <Content style={{ padding: "32px 16px 16px 16px" }} id="main">
          {children}
        </Content>
      </Layout>
    </Layout>
  );
}

export default Main;
